package app

import (
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
)

// cleanCache interactively removes all cached entries.
// Returns 0 on success, 1 on error.
func (a *App) cleanCache() int {
	c := cache.New()
	if err := c.Clean(a.Stdin); err != nil {
		logger.Error("failed to clean cache: %v", err)
		return 1
	}
	return 0
}
