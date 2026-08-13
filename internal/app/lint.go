package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/lint"
	"github.com/TheRootDaemon/tlgc/logger"
)

// lintPages runs the linter over every page reachable
// from the CLI paths and reports each violation to stderr.
//
// It returns 0 when every file passed without violations
// and 1 otherwise;
// a failure to collect or open any file is logged and also yields 1.
func (a *App) lintPages(cli *cmd.CLI) int {
	files, err := collectFiles(cli.Page)
	if err != nil {
		logger.Error("%v", err)
		return 1
	}

	failed := false
	for _, file := range files {
		result, err := a.lintFile(file, cli.Ignore...)
		if err != nil {
			logger.Error("%v", err)
			failed = true
			continue
		}

		for _, e := range result.Errors {
			a.writeLintError(file, e, cli.Tabular)
			failed = true
		}
	}

	if failed {
		return 1
	}

	return 0
}

// lintFile opens the page at path and runs every applicable lint rule
// over it, returning the violations found.
//
// Error codes passed in ignore are suppressed by the linter.
//
// On any open or lint failure it returns an empty Result together with the error.
func (a *App) lintFile(
	path string,
	ignore ...string,
) (*lint.Result, error) {
	root, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return &lint.Result{}, err
	}
	defer func() {
		_ = root.Close()
	}()

	f, err := root.Open(filepath.Base(path))
	if err != nil {
		return &lint.Result{}, err
	}
	defer func() {
		_ = f.Close()
	}()

	return lint.Lint(f, ignore...)
}

// writeLintError reports a single lint violation for path to a.Stderr.
func (a *App) writeLintError(
	path string,
	e lint.Error,
	tabular bool,
) {
	if tabular {
		_, _ = fmt.Fprintf(
			a.Stderr,
			"%s\t%d\t%s\t%s\t\n",
			path,
			e.Line,
			e.Code,
			e.Description,
		)
		return
	}

	_, _ = fmt.Fprintf(
		a.Stderr,
		"%s:%d: %s %s\n",
		path,
		e.Line,
		e.Code,
		e.Description,
	)
}
