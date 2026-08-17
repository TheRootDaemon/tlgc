package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/TheRootDaemon/tlgc/format"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/termcolor"
)

// cacheInfo prints detailed cache metadata including directory,
// age, auto-update status, mirror, platforms, and per-language breakdown.
// Returns 0 on success, 1 on error.
func (a *App) cacheInfo() int {
	c := cache.New()
	info, err := c.Info()
	if err != nil {
		logger.Error("failed to get cache info: %v", err)
		return 1
	}

	if err := a.printCacheHeader(info); err != nil {
		logger.Error("%v", err)
		return 1
	}

	if err := a.printAutoUpdate(info); err != nil {
		logger.Error("%v", err)
		return 1
	}

	if err := a.printMirror(info); err != nil {
		logger.Error("%v", err)
		return 1
	}

	if err := a.printLanguages(info); err != nil {
		logger.Error("%v", err)
		return 1
	}

	return 0
}

// printCacheHeader prints the cache directory and last update time.
func (a *App) printCacheHeader(info *cache.InfoResult) error {
	cacheDir := termcolor.Sprint("red", info.CacheDir)
	lastUpdated := termcolor.Sprint("bold blue", info.Age)

	_, err := fmt.Fprintf(
		a.Stdout,
		"Cache: %s (last update: %s ago)\n",
		cacheDir,
		lastUpdated,
	)

	return err
}

// printAutoUpdate prints the automatic cache update configuration.
func (a *App) printAutoUpdate(info *cache.InfoResult) error {
	if !info.AutoUpdate {
		_, err := fmt.Fprintf(
			a.Stdout,
			"Auto update: %s\n",
			termcolor.Sprint("bold red", "disabled"),
		)
		return err
	}

	maxAge, err := format.ValidateDurationOverflow(info.MaxAge)
	if err != nil {
		return err
	}

	frequency := termcolor.Sprint(
		"bold blue",
		format.DurationFmt(maxAge),
	)
	remaining := termcolor.Sprint(
		"bold blue",
		format.DurationFmt(
			max(
				0,
				maxAge*time.Hour-info.AgeDuration,
			),
		),
	)

	_, err = fmt.Fprintf(
		a.Stdout,
		"Auto update: every %s (next in %s)\n",
		frequency,
		remaining,
	)
	return err
}

// printMirror prints the configured tldr-pages mirror URL.
func (a *App) printMirror(info *cache.InfoResult) error {
	mirror := termcolor.Sprint("green", info.Mirror)
	_, err := fmt.Fprintf(
		a.Stdout,
		"Mirror: %s\n",
		mirror,
	)
	return err
}

// printLanguages prints per-language page counts and the total number
// of cached pages.
func (a *App) printLanguages(info *cache.InfoResult) error {
	if _, err := fmt.Fprintln(a.Stdout, "Installed languages:"); err != nil {
		return err
	}

	maxWidth := len("total")
	for _, ls := range info.LanguageStats {
		if len(ls.Language) > maxWidth {
			maxWidth = len(ls.Language)
		}
	}

	totalPages := termcolor.Fprintf("bold blue", "%d", info.TotalPages)
	if _, err := fmt.Fprintf(
		a.Stdout,
		"  %-*s : %s pages\n",
		maxWidth,
		"total",
		totalPages,
	); err != nil {
		return err
	}

	for _, ls := range info.LanguageStats {
		pages := termcolor.Fprintf("bold blue", "%d", ls.Pages)
		if _, err := fmt.Fprintf(
			a.Stdout,
			"  %-*s : %s pages%s\n",
			maxWidth,
			ls.Language,
			pages,
			formatPlatformBreakdown(ls.Platforms),
		); err != nil {
			return err
		}
	}

	return nil
}

// formatPlatformBreakdown formats per-platform page counts
// as a parenthesized suffix, e.g. " (common: 3000, linux: 2500)".
func formatPlatformBreakdown(platforms []cache.PlatformInfo) string {
	if len(platforms) == 0 {
		return ""
	}

	parts := make([]string, len(platforms))
	for i, p := range platforms {
		pages := termcolor.Fprintf("bold blue", "%d", p.Pages)
		parts[i] = fmt.Sprintf("%s: %s", p.Name, pages)
	}

	return " (" + strings.Join(parts, ", ") + ")"
}
