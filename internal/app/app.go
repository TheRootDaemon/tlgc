package app

import (
	"io"
	"os"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/locale"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/platform"
)

// App is the main application struct that holds I/O streams and configuration.
type App struct {
	// Stdout is the writer for standard output.
	Stdout io.Writer

	// Stderr is the writer for diagnostic output.
	Stderr io.Writer

	// ConfigPath is the path to the configuration file, if set.
	ConfigPath string
}

// Option configures the App.
type Option func(*App)

// WithStdout sets the standard output writer for the App.
func WithStdout(w io.Writer) Option {
	return func(a *App) {
		a.Stdout = w
	}
}

// WithStderr sets the standard error writer for the App.
func WithStderr(w io.Writer) Option {
	return func(a *App) {
		a.Stderr = w
	}
}

// WithConfigPath sets the configuration file path for the App.
func WithConfigPath(path string) Option {
	return func(a *App) {
		a.ConfigPath = path
	}
}

// New creates a new App with the given options.
// It defaults Stdout to os.Stdout and Stderr to os.Stderr.
func New(opts ...Option) *App {
	a := &App{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Run dispatches the CLI command to the appropriate handler.
// It initializes the config if needed, then delegates to the matching
// sub-handler based on the CLI flags. Returns 0 on success, 1 on error.
func (a *App) Run(cli *cmd.CLI) int {
	needsConfig := !cli.GenConfig && !cli.ConfigPath && !cli.ShowVersion && !cli.ShowHelp

	if needsConfig {
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}
	}

	switch {
	case cli.Update:
		return a.updateCache(cli)
	case cli.List:
		return a.listPages(cli)
	case cli.ListAll:
		return a.listAllPages()
	case cli.Search != "":
		return a.searchPages(cli)
	case cli.ListPlatforms:
		return a.listPlatforms()
	case cli.ListLanguages:
		return a.listLanguages()
	case cli.Info:
		return a.cacheInfo()
	case cli.Render != "":
		return a.renderLocalFile(cli)
	case cli.GenConfig:
		return a.genConfig()
	case cli.ConfigPath:
		return a.configPath()
	case len(cli.Page) > 0:
		return a.lookupAndRenderPage(cli)
	default:
		return 0
	}
}

// resolveLanguages returns the resolved list of languages from the CLI flag,
// config file, or system locale, in that order of precedence.
func (a *App) resolveLanguages(flagLangs []string) []string {
	if len(flagLangs) > 0 {
		return flagLangs
	}
	if cfgLangs := config.Cache().Languages; len(cfgLangs) > 0 {
		return cfgLangs
	}
	var langs []string
	locale.GetLanguages(&langs)
	return langs
}

// resolvePlatform returns the resolved platform string from the CLI flag
// or the default platform.
func (a *App) resolvePlatform(flagPlatform string) string {
	if flagPlatform != "" {
		return platform.Resolve(flagPlatform)
	}
	return platform.Default()
}
