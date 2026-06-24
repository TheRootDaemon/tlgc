package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindResult contains the pages found for a command lookup.
//
// Matches contains pages found in the requested platform
// and/or the common platform.
//
// Fallbacks contains pages found only in other platforms.
type FindResult struct {
	Matches   []string
	Fallbacks []string
}

// Find locates a command page in the cache.
//
// It searches the requested platform first,
// then the common platform.
// If no match is found there, it searches the remaining platforms
// and returns those matches as fallbacks.
//
// Language directories are searched in the order provided by languages.
func (c *Cache) Find(query, platform string, languages []string) (*FindResult, error) {
	languageDirectories := c.languagesToDirectories(languages, false)
	if len(languageDirectories) == 0 {
		return nil, fmt.Errorf("no matching language directories found in cache")
	}

	platforms, err := c.getPlatforms()
	if err != nil {
		return nil, err
	}

	if platform != "common" {
		if !platformExists(platforms, platform) {
			return nil, fmt.Errorf("platform %q does not exist", platform)
		}
	}

	file := query + ".md"

	matches := c.primaryMatches(
		platform,
		file,
		languageDirectories,
	)
	fallbacks := c.fallbackMatches(
		file,
		platform,
		platforms,
		languageDirectories,
	)

	return &FindResult{
		Matches:   matches,
		Fallbacks: fallbacks,
	}, nil
}

// primaryMatches searches for file in the requested platform
// and the common platform, returning matches in priority order.
func (c *Cache) primaryMatches(
	platform,
	file string,
	languageDirectories []string,
) []string {
	var results []string

	if platform != "common" {
		if path := c.findPageFor(
			file,
			platform,
			languageDirectories,
		); path != "" {
			results = append(results, path)
		}
	}

	if path := c.findPageFor(
		file,
		"common",
		languageDirectories,
	); path != "" {
		results = append(results, path)
	}

	return results
}

// fallbackMatches searches for file in all platforms
// other than the requested platform and common,
// returning any matches found.
func (c *Cache) fallbackMatches(
	file,
	platform string,
	platforms,
	languageDirectories []string,
) []string {
	var results []string

	for _, p := range platforms {
		if p == platform || p == "common" {
			continue
		}

		if path := c.findPageFor(
			file,
			p,
			languageDirectories,
		); path != "" {
			results = append(results, path)
		}
	}

	return results
}

// findPageFor searches for fname within platform across language dirs.
// Returns the first match found (respects language priority).
func (c *Cache) findPageFor(fname, platform string, languageDirectories []string) string {
	for _, languageDirectory := range languageDirectories {
		path := filepath.Join(
			c.dir,
			languageDirectory,
			platform,
			fname,
		)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
