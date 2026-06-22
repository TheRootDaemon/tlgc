package cache

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TheRootDaemon/tlgc/slice"
)

// ListFor returns all page names in the give platform (plus common).
func (c *Cache) ListFor(platform string) ([]string, error) {
	if _, err := c.getPlatforms(); err != nil {
		return nil, err
	}

	pages, err := c.listDirectory(platform, englishDirectory)
	if err != nil {
		return nil, err
	}

	if platform != "common" {
		common, err := c.listDirectory("common", englishDirectory)
		if err != nil {
			return nil, err
		}

		pages = append(pages, common...)
	}

	sort.Strings(pages)
	return slice.Dedup(pages), nil
}

// ListAll returns all page names across all platforms in English.
func (c *Cache) ListAll() ([]string, error) {
	platforms, err := c.getPlatforms()
	if err != nil {
		return nil, err
	}

	var pages []string
	for _, platform := range platforms {
		platformPages, err := c.listDirectory(platform, englishDirectory)
		if err != nil {
			return nil, err
		}

		pages = append(pages, platformPages...)
	}

	sort.Strings(pages)
	return slice.Dedup(pages), nil
}

// ListPlatforms returns the available platform directories.
func (c *Cache) ListPlatforms() ([]string, error) {
	return c.getPlatforms()
}

// ListLanguages returns the installed language codes (without the "pages." prefix).
func (c *Cache) ListLanguages() ([]string, error) {
	directories, err := c.getLanguageDirectories()
	if err != nil {
		return nil, err
	}

	languages := make([]string, len(directories))
	for i, directory := range directories {
		languages[i] = strings.TrimPrefix(directory, "pages.")
	}

	return languages, nil
}

// listDir lists page names (without .md extension) in lang/platform.
// Returns empty slice if the directory does not exist.
func (c *Cache) listDirectory(platform, languageDirectory string) ([]string, error) {
	dir := filepath.Join(c.dir, languageDirectory, platform)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var pages []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if before, ok := strings.CutSuffix(name, ".md"); ok {
			pages = append(
				pages,
				before,
			)
		}
	}

	return pages, nil
}
