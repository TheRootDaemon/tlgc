package app

import (
	"fmt"

	"github.com/TheRootDaemon/tlgc/exitcode"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
)

// genConfig prints a default configuration file to stdout.
// It returns [exitcode.Success] on success.
// Other failures return [exitcode.Default].
func (a *App) genConfig() int {
	cfg, err := config.DefaultConfig()
	if err != nil {
		logger.Error("failed to generate config: %v", err)
		return exitcode.Default
	}

	if _, err := fmt.Fprint(a.Stdout, cfg); err != nil {
		logger.Error("%v", err)
		return exitcode.Default
	}
	return exitcode.Success
}

// configPath prints the configuration file path to stdout.
// It returns [exitcode.Success] on success.
// Other failures return [exitcode.Default].
func (a *App) configPath() int {
	if _, err := fmt.Fprintln(a.Stdout, config.ConfigPath()); err != nil {
		logger.Error("%v", err)
		return exitcode.Default
	}
	return exitcode.Success
}
