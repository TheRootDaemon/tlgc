package app

import (
	"context"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/internal/upstream"
	"github.com/TheRootDaemon/tlgc/logger"
)

// updateCache downloads the latest tldr-pages
// for the configured languages.
// Returns 0 on success, 1 on error.
func (a *App) updateCache(cli *cmd.CLI) int {
	c := cache.New()
	languages := a.resolveLanguages(cli.Languages)
	client := upstream.New()

	if err := c.Update(context.Background(), languages, client); err != nil {
		logger.Error("failed to update cache: %v", err)
		return 1
	}
	return 0
}

// shouldAutoUpdate reports whether a stale cache maybe updated
// automatically for the current invocation.
// Automatic updates run only for page-serving operations
// and are disabled in offline mode,
// and when explicit cache maintenance was requested.
func shouldAutoUpdate(cli *cmd.CLI) bool {
	if cli.Offline || cli.Update || cli.CleanCache {
		return false
	}
	return (len(cli.Page) > 0 && !cli.Browse && !cli.Lint && !cli.Format) ||
		cli.Browse ||
		cli.Search != "" ||
		cli.List
}

// autoUpdate downloads the latest tldr-pages
// for the configured languages.
// It returns the error when the update fails.
func (a *App) autoUpdate(cli *cmd.CLI) error {
	c := cache.New()
	languages := a.resolveLanguages(cli.Languages)
	client := upstream.New()

	return c.Update(context.Background(), languages, client)
}
