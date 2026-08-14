package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/lint"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/termcolor"
)

// lintViolation couples a lint violation with the page it was found in.
type lintViolation struct {
	path string
	err  lint.Error
}

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
	var rows []lintViolation
	for _, file := range files {
		result, err := a.lintFile(file, cli.Ignore...)
		if err != nil {
			logger.Error("%v", err)
			failed = true
			continue
		}

		for _, e := range result.Errors {
			if cli.Tabular {
				rows = append(rows, lintViolation{path: file, err: e})
			} else {
				a.writeLintError(file, e)
			}
			failed = true
		}
	}

	if cli.Tabular {
		a.writeTabular(rows)
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
func (a *App) writeLintError(path string, e lint.Error) {
	_, _ = fmt.Fprintf(
		a.Stderr,
		"%s:%d: %s %s\n",
		path,
		e.Line,
		e.Code,
		e.Description,
	)
}

// writeTabular writes a decorated,
// aligned table of lint violations to a.Stderr,
// mirroring the header style of the search table.
// It writes nothing when there are no rows.
func (a *App) writeTabular(rows []lintViolation) {
	if len(rows) == 0 {
		return
	}

	fileW, lineW, codeW := lintColumnWidths(rows)

	_, _ = fmt.Fprintln(
		a.Stderr,
		termcolor.Fprintf(
			"bold",
			"%-*s %-*s %-*s %s",
			fileW,
			"File",
			lineW,
			"Line",
			codeW,
			"Code",
			"Description",
		),
	)

	for _, r := range rows {
		_, _ = fmt.Fprintf(
			a.Stderr,
			"%-*s %-*d %-*s %s\n",
			fileW,
			r.path,
			lineW,
			r.err.Line,
			codeW,
			r.err.Code,
			r.err.Description,
		)
	}
}

// lintColumnWidths returns the widths required to display the
// File, Line, and Code columns without truncation.
func lintColumnWidths(rows []lintViolation) (int, int, int) {
	fileW := len("File")
	lineW := len("Line")
	codeW := len("Code")

	for _, r := range rows {
		fileW = max(fileW, len(r.path))
		lineW = max(lineW, len(strconv.Itoa(r.err.Line)))
		codeW = max(codeW, len(r.err.Code))
	}

	return fileW, lineW, codeW
}
