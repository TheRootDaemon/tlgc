package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/slice"
)

// Cache manages the local tldr-pages cache on disk.
type Cache struct {
	dir       string
	platforms atomic.Value
}

// New creates a Cache using the cache directory from the config singleton.
func New() *Cache {
	return &Cache{
		dir: config.Cache().Dir,
	}
}

// Dir returns the cache directory path.
func (c *Cache) Dir() string {
	return c.dir
}

// subdirExists reports whether name is a subdirectory of the cache.
func (c *Cache) subDirExists(name string) bool {
	fi, err := os.Stat(filepath.Join(c.dir, name))
	return err == nil && fi.IsDir()
}

// Age returns the cache age based on the checksum file's mtime.
// Falls back to the cache directory mtime
// if the checksum file does not exist.
func (c *Cache) Age() (time.Duration, error) {
	sumfile := filepath.Join(c.dir, checksumFile)
	fi, err := os.Stat(sumfile)
	if err != nil {
		fi, err = os.Stat(c.dir)
		if err != nil {
			return 0, err
		}
	}

	mod := fi.ModTime()
	age := time.Since(mod)

	if age < 0 {
		return 0, fmt.Errorf("cache mtime is in the future: clock issue")
	}

	return age, nil
}

// getPlatforms discovers available platforms
// by reading directories under pages.en/.
// Results are cached after first load.
func (c *Cache) getPlatforms() ([]string, error) {
	if p, ok := c.platforms.Load().([]string); ok {
		return p, nil
	}

	entries, err := os.ReadDir(filepath.Join(c.dir, englishDirectory))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %s", englishDirectory, err)
	}

	var platforms []string
	for _, e := range entries {
		if e.IsDir() {
			platforms = append(platforms, e.Name())
		}
	}

	if len(platforms) == 0 {
		return nil, fmt.Errorf("'%s' contains no platform directories", englishDirectory)
	}

	sort.Strings(platforms)
	c.platforms.Store(platforms)
	return platforms, nil
}

// getLanguageDirectories returns all pages.* directories in the cache.
func (c *Cache) getLanguageDirectories() ([]string, error) {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "pages.") {
			dirs = append(dirs, e.Name())
		}
	}

	sort.Strings(dirs)
	return dirs, nil
}

// languagesToDirs converts language codes to pages.xx dirs that exist.
// If sort is true, results are sorted; otherwise only adjacent duplicates are removed.
func (c *Cache) languagesToDirectories(languages []string, sortFlag bool) []string {
	var dirs []string
	for _, lang := range languages {
		dir := "pages." + lang
		if c.subDirExists(dir) {
			dirs = append(dirs, dir)
		}
	}

	if sortFlag {
		sort.Strings(dirs)
	}

	dirs = slice.Dedup(dirs)
	return dirs
}
