package app

import (
	"fmt"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
)

// searchPages searches cached pages for the given query.
// Returns 0 on success, 1 on error.
func (a *App) searchPages(cli *cmd.CLI) int {
	c := cache.New()
	p := a.resolvePlatform(cli.Platform)
	languages := a.resolveLanguages(cli.Languages)

	results, err := c.Search(cli.Search, p, languages)
	if err != nil {
		logger.Error("search failed: %v", err)
		return 1
	}

	for _, r := range results {
		if _, err := fmt.Fprintf(
			a.Stdout,
			"%s/%s\n",
			r.Platform,
			r.Page,
		); err != nil {
			logger.Error("%w", err)
			return 1
		}
	}
	return 0
}
