package app

import (
	"fmt"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
)

// genConfig prints a default configuration file to stdout.
// Returns 0 on success, 1 on error.
func (a *App) genConfig() int {
	cfg, err := config.DefaultConfig()
	if err != nil {
		logger.Error("failed to generate config: %w", err)
		return 1
	}

	if _, err := fmt.Fprint(a.Stdout, cfg); err != nil {
		logger.Error("%w", err)
		return 1
	}
	return 0
}

// configPath prints the configuration file path to stdout.
// Returns 0 on success.
func (a *App) configPath() int {
	if _, err := fmt.Fprintln(a.Stdout, config.ConfigPath()); err != nil {
		logger.Error("%w", err)
		return 1
	}
	return 0
}
