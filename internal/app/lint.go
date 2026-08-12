package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/lint"
	"github.com/TheRootDaemon/tlgc/logger"
)

// lintPages validates the tldr pages
// under the given paths
// and reports any lint errors to stderr.
// Returns 0 when no errors are found, 1 otherwise.
func (a *App) lintPages(cli *cmd.CLI) int {
	files, err := collectFiles(cli.Page)
	if err != nil {
		logger.Error("%v", err)
		return 1
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		logger.Error("%v", err)
		return 1
	}
	defer func() {
		_ = root.Close()
	}()

	failed := false
	for _, file := range files {
		f, err := root.Open(file)
		if err != nil {
			logger.Error("%v", err)
			failed = true
			continue
		}

		result, err := lint.Lint(f, cli.Ignore...)
		_ = f.Close()
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

// writeLintError reports a single lint error to stderr
// in the default or tabular reference format.
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

// collectFiles expands each input path
// into a flat list of page files.
// Individual .md files are included as-is,
// while directories are walked recursively
// to collect .md files.
func collectFiles(paths []string) ([]string, error) {
	var files []string

	for _, path := range paths {
		pathFiles, err := collectPathFiles(path)
		if err != nil {
			return nil, err
		}
		files = append(files, pathFiles...)
	}

	return files, nil
}

// collectPathFiles expands a single path into page files.
// A .md file is returned as-is,
// a directory is walked recursively for .md files,
// and other file types are ignored.
func collectPathFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if filepath.Ext(path) == ".md" {
			return []string{path}, nil
		}
		return nil, nil
	}

	var files []string
	if err = filepath.WalkDir(
		path,
		func(entry string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(entry) == ".md" {
				files = append(files, entry)
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return files, nil
}
