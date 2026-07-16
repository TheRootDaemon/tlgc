package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TheRootDaemon/tlgc/format"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
)

// LanguageInfo contains cache statistics for a single language.
type LanguageInfo struct {
	// Pages is the number of cached pages for this language.
	Pages int

	// Language is the language name (e.g. "en", "pt", "es").
	Language string

	// Platforms contains per-platform page counts for this language.
	Platforms []PlatformInfo
}

// PlatformInfo contains page statistics for a single platform.
type PlatformInfo struct {
	// Name is the platform name (e.g. "common", "linux", "osx").
	Name string

	// Pages is the number of cached pages for this platform.
	Pages int
}

// InfoResult contains information about the current cache state.
type InfoResult struct {
	// CacheDir is the absolute path to the cache directory.
	CacheDir string

	// Age is a human-readable string representing the cache age.
	Age string

	// AgeDuration is the raw cache age duration.
	AgeDuration time.Duration

	// MaxAge is the maximum cache age in hours before a refresh is due.
	MaxAge uint64

	// AutoUpdate indicates whether automatic cache updates are enabled.
	AutoUpdate bool

	// Mirror is the URL used to download tldr-pages archives.
	Mirror string

	// Platforms lists the available platforms discovered in the cache.
	Platforms []string

	// LanguageStats contains per-language page statistics.
	LanguageStats []LanguageInfo

	// TotalPages is the total number of cached pages across all languages.
	TotalPages int
}

// Age returns the cache age based on the checksum file's mtime.
// Falls back to the cache directory mtime
// if the checksum file does not exist.
func (c *Cache) Age() (time.Duration, error) {
	sumfile := filepath.Join(c.dir, checksumFile)
	fi, err := os.Stat(sumfile)
	if err != nil {
		logger.Debug(
			"cache age: stat failed for %s, falling back to %s",
			sumfile,
			c.dir,
		)
		fi, err = os.Stat(c.dir)
		if err != nil {
			return 0, err
		}
	}

	mod := fi.ModTime()
	age := time.Since(mod)

	if age < 0 {
		logger.Warn("cache mtime is in the future: clock may be wrong")
		return 0, fmt.Errorf("cache mtime is in the future: clock issue")
	}

	return age, nil
}

// Info returns a snapshot of the current cache state,
// including its location, age, configuration,
// per-language page counts, and total page count.
func (c *Cache) Info() (*InfoResult, error) {
	logger.Debug("cache dir=%q", c.dir)

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
		AgeDuration:   age,
		MaxAge:        cfg.MaxAge,
		AutoUpdate:    cfg.AutoUpdate,
		Mirror:        cfg.Mirror,
		Platforms:     platforms,
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
	total := 0
	languageStats := make(
		[]LanguageInfo,
		0,
		len(languageDirectories),
	)

	for _, languageDirectory := range languageDirectories {
		lang := strings.TrimPrefix(
			languageDirectory,
			"pages.",
		)

		platformInfos, count, err := c.platformStats(
			languageDirectory,
			platforms,
		)
		if err != nil {
			return nil, 0, err
		}

		languageStats = append(
			languageStats,
			LanguageInfo{
				Language:  lang,
				Pages:     count,
				Platforms: platformInfos,
			},
		)

		total += count
	}

	return languageStats, total, nil
}

func (c *Cache) platformStats(
	languageDirectory string,
	platforms []string,
) ([]PlatformInfo, int, error) {
	total := 0
	var infos []PlatformInfo

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

		count := len(pages)
		infos = append(
			infos,
			PlatformInfo{
				Name:  platform,
				Pages: count,
			},
		)

		total += count
	}

	return infos, total, nil
}
