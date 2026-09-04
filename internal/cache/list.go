package cache

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/slice"
)

// ListFor returns all page names in the give platform (plus common).
// across the specified languages.
//
// If nothing languages is empty it defaults to pages.en.
func (c *Cache) ListFor(platform string, languages []string) ([]string, error) {
	logger.Debug("platform=%q, language=%q", platform, languages)
	if _, err := c.getPlatforms(); err != nil {
		return nil, err
	}

	languageDirectories := c.languagesToDirectories(languages, false)
	if len(languageDirectories) == 0 {
		languageDirectories = append(languageDirectories, englishDirectory)
	}

	var pages []string
	for _, languageDirectory := range languageDirectories {
		platformPages, err := c.listDirectory(platform, languageDirectory)
		if err != nil {
			return nil, err
		}

		pages = append(pages, platformPages...)
		if platform != "common" {
			commonPages, err := c.listDirectory("common", languageDirectory)
			if err != nil {
				return nil, err
			}
			pages = append(pages, commonPages...)
		}
	}

	sort.Strings(pages)
	return slice.Dedup(pages), nil
}

// ListAll returns all page names across all platforms
// for the specified languages.
//
// If nothing languages is empty it defaults to pages.en.
func (c *Cache) ListAll(languages []string) ([]string, error) {
	platforms, err := c.getPlatforms()
	logger.Debug("%d platforms", len(platforms))
	if err != nil {
		return nil, err
	}

	languageDirectories := c.languagesToDirectories(languages, false)
	if len(languageDirectories) == 0 {
		languageDirectories = append(languageDirectories, englishDirectory)
	}

	var pages []string
	for _, languageDirectory := range languageDirectories {
		for _, platform := range platforms {
			platformPages, err := c.listDirectory(platform, languageDirectory)
			if err != nil {
				return nil, err
			}

			pages = append(pages, platformPages...)
		}
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
	logger.Debug("found %d languages", len(directories))
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
