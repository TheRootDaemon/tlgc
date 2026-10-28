package app

import (
	"github.com/TheRootDaemon/tlgc/exitcode"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
)

// cleanCache interactively removes all cached entries.
// It returns [exitcode.Success] on success.
// Other failures return [exitcode.Default].
func (a *App) cleanCache() int {
	c := cache.New()
	if err := c.Clean(a.Stdin); err != nil {
		logger.Error("failed to clean cache: %v", err)
		return exitcode.Default
	}
	return exitcode.Success
}
