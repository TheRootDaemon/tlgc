package app

import (
	"fmt"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
)

// listPages lists pages for a specific platform from the cache.
// Returns 0 on success, 1 on error.
func (a *App) listPages(cli *cmd.CLI) int {
	c := cache.New()
	p := a.resolvePlatform(cli.Platform)
	languages := a.resolveLanguages(cli.Languages)
	pages, err := c.ListFor(p, languages)
	if err != nil {
		logger.Error("failed to list pages: %v", err)
		return 1
	}

	for _, page := range pages {
		if _, err := fmt.Fprintln(a.Stdout, page); err != nil {
			logger.Error("%v", err)
			return 1
		}
	}
	return 0
}

// listAllPages lists all cached pages across all platforms.
// Returns 0 on success, 1 on error.
func (a *App) listAllPages(cli *cmd.CLI) int {
	c := cache.New()
	languages := a.resolveLanguages(cli.Languages)
	pages, err := c.ListAll(languages)
	if err != nil {
		logger.Error("failed to list pages: %v", err)
		return 1
	}

	for _, page := range pages {
		if _, err := fmt.Fprintln(a.Stdout, page); err != nil {
			logger.Error("%v", err)
			return 1
		}
	}
	return 0
}

// listPlatforms lists all available platforms from the cache.
// Returns 0 on success, 1 on error.
func (a *App) listPlatforms() int {
	c := cache.New()
	platforms, err := c.ListPlatforms()
	if err != nil {
		logger.Error("failed to list platforms: %v", err)
		return 1
	}

	for _, p := range platforms {
		if _, err := fmt.Fprintln(a.Stdout, p); err != nil {
			logger.Error("%v", err)
			return 1
		}
	}
	return 0
}

// listLanguages lists all available languages from the cache.
// Returns 0 on success, 1 on error.
func (a *App) listLanguages() int {
	c := cache.New()
	languages, err := c.ListLanguages()
	if err != nil {
		logger.Error("failed to list languages: %v", err)
		return 1
	}

	for _, l := range languages {
		if _, err := fmt.Fprintln(a.Stdout, l); err != nil {
			logger.Error("%v", err)
			return 1
		}
	}
	return 0
}
