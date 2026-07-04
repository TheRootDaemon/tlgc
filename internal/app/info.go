package app

import (
	"fmt"

	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
)

// cacheInfo prints cache metadata (directory, age, page count, etc.).
// Returns 0 on success, 1 on error.
func (a *App) cacheInfo() int {
	c := cache.New()
	info, err := c.Info()
	if err != nil {
		logger.Error("failed to get cache info: %v", err)
		return 1
	}

	if _, err := fmt.Fprintf(
		a.Stdout,
		`Cache: %s
Cache age: %s
Total pages: %d
Auto update: %v
Max age (hours): %d
`,
		info.CacheDir,
		info.Age,
		info.TotalPages,
		info.AutoUpdate,
		info.MaxAge,
	); err != nil {
		logger.Error("%w", err)
		return 1
	}

	for _, ls := range info.LanguageStats {
		if _, err := fmt.Fprintf(
			a.Stdout,
			"%s: %d pages\n",
			ls.Language,
			ls.Pages,
		); err != nil {
			logger.Error("%w", err)
			return 1
		}
	}
	return 0
}
