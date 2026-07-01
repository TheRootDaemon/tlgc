package cache

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/TheRootDaemon/tlgc/logger"
)

// SearchResult represents a matching page from a search.
type SearchResult struct {
	// Page is the page name (filename without the .md extension).
	Page string

	// Language is the language of the page (e.g. "en").
	Language string

	// Platform is the platform of the page (e.g. "linux", "common").
	Platform string
}

// Search performs a case-insensitive substring search
// across the requested platforms and languages.
//
// If platform is empty, all platforms are searched.
// If a specific platform is requested,
// that platform and common are searched.
// Results are returned sorted by page name.
func (c *Cache) Search(query, platform string, languages []string) ([]SearchResult, error) {
	platforms, err := c.resolvePlatforms(platform)
	if err != nil {
		return nil, err
	}

	languageDirectories := c.languagesToDirectories(languages, false)
	if len(languageDirectories) == 0 {
		return nil, fmt.Errorf(
			"no installed languages match the requested languages",
		)
	}

	results := c.searchPages(
		query,
		platforms,
		languageDirectories,
	)

	if len(results) == 0 {
		return nil, fmt.Errorf("no pages matched your search term")
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return results[i].Page < results[j].Page
		},
	)

	return results, nil
}

// resolvePlatforms validates platform
// and returns the platforms
// that should be searched.
//
// If platform is empty, all available platforms are returned.
// If platform is "common", only common is returned.
// Otherwise, the requested platform and common are returned.
func (c *Cache) resolvePlatforms(platform string) ([]string, error) {
	platforms, err := c.getPlatforms()
	if err != nil {
		return nil, err
	}

	if platform != "" && !platformExists(platforms, platform) {
		return nil, fmt.Errorf(
			"platform %q does not exist, possible values are: %s",
			platform,
			strings.Join(platforms, ", "),
		)
	}

	switch {
	case platform == "common":
		return []string{"common"}, nil
	case platform != "":
		return []string{platform, "common"}, nil
	default:
		return platforms, nil
	}
}

// searchPages searches all platform/language combinations
// and returns the matching pages.
func (c *Cache) searchPages(
	query string,
	platforms, languageDirectories []string,
) []SearchResult {
	query = strings.ToLower(query)
	var results []SearchResult

	for _, languageDirectory := range languageDirectories {
		for _, platform := range platforms {
			matches := c.searchDirectory(
				query,
				platform,
				languageDirectory,
			)

			results = append(results, matches...)
		}
	}

	return results
}

// searchDirectory searches a single language/platform directory
// for pages whose names contain query.
func (c *Cache) searchDirectory(query, platform, languageDirectory string) []SearchResult {
	pages, err := c.listDirectory(platform, languageDirectory)
	if err != nil {
		logger.Debug(
			"error listing %s/%s: %s",
			languageDirectory,
			platform,
			err,
		)

		return nil
	}

	var results []SearchResult
	language := strings.TrimPrefix(
		languageDirectory,
		"pages.",
	)

	for _, page := range pages {
		if strings.Contains(
			strings.ToLower(page),
			query,
		) {
			results = append(
				results,
				SearchResult{
					Page:     page,
					Language: language,
					Platform: platform,
				},
			)
		}
	}

	return results
}

// platformExists checks whether the given platform is in the list.
func platformExists(platforms []string, platform string) bool {
	return slices.Contains(platforms, platform)
}
