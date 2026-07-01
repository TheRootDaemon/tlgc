package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TheRootDaemon/tlgc/format"
	"github.com/TheRootDaemon/tlgc/internal/config"
)

// LanguageInfo contains cache statistics for a single language.
type LanguageInfo struct {
	// Pages is the number of cached pages for this language.
	Pages int

	// Language is the language name (e.g. "en", "pt", "es").
	Language string
}

// InfoResult contains information about the current cache state.
type InfoResult struct {
	// AutoUpdate indicates whether automatic cache updates are enabled.
	AutoUpdate bool

	// TotalPages is the total number of cached pages across all languages.
	TotalPages int

	// MaxAge is the maximum cache age in seconds before a refresh is due.
	MaxAge uint64

	// CacheDir is the absolute path to the cache directory.
	CacheDir string

	// Age is a human-readable string representing the cache age.
	Age string

	// LanguageStats contains per-language page statistics.
	LanguageStats []LanguageInfo
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

// Info returns a snapshot of the current cache state,
// including its location, age, configuration,
// per-language page counts, and total page count.
func (c *Cache) Info() (*InfoResult, error) {
	fi, err := os.Stat(c.dir)
	if err != nil {
		return nil, fmt.Errorf("cache directory %q: %s", c.dir, err)
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("cache path %q is not a directory", c.dir)
	}

	age, err := c.Age()
	if err != nil {
		return nil, err
	}

	cfg := config.Cache()

	languageDirectories, err := c.getLanguageDirectories()
	if err != nil {
		return nil, err
	}

	platforms, err := c.getPlatforms()
	if err != nil {
		return nil, err
	}

	languageStats, total, err := c.languageStats(
		platforms,
		languageDirectories,
	)
	if err != nil {
		return nil, err
	}

	return &InfoResult{
		CacheDir:      c.dir,
		Age:           format.DurationFmt(age),
		MaxAge:        cfg.MaxAge,
		AutoUpdate:    cfg.AutoUpdate,
		LanguageStats: languageStats,
		TotalPages:    total,
	}, nil
}

// languageStats counts cached pages for each language and
// returns the per-language statistics
// along with the total page count.
func (c *Cache) languageStats(
	platforms,
	languageDirectories []string,
) ([]LanguageInfo, int, error) {
	var languageStats []LanguageInfo
	total := 0

	for _, languageDirectory := range languageDirectories {
		lang := strings.TrimPrefix(
			languageDirectory,
			"pages.",
		)
		count := 0

		for _, platform := range platforms {
			if !c.subDirExists(
				filepath.Join(
					languageDirectory,
					platform,
				),
			) {
				continue
			}

			pages, err := c.listDirectory(
				platform,
				languageDirectory,
			)
			if err != nil {
				return nil, 0, err
			}

			count += len(pages)
		}

		languageStats = append(languageStats, LanguageInfo{
			Language: lang,
			Pages:    count,
		})
		total += count
	}

	return languageStats, total, nil
}
