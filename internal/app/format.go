package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/lint"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/termcolor"
)

// formatPages formats every page reachable from the CLI paths,
// writing each result to stdout,
// back to the source file, or to a single output file
// depending on the flags.
//
// It returns 0 when every page formatted without lint errors and 1
// otherwise.
func (a *App) formatPages(cli *cmd.CLI) int {
	files, err := collectFiles(cli.Page)
	if err != nil {
		logger.Error("%v", err)
		return 1
	}

	if cli.Output != "" && len(files) != 1 {
		logger.Error(
			"flag %s requires a single file",
			termcolor.Sprint("bold blue", "--output"),
		)
		return 1
	}

	failed := false
	for _, file := range files {
		hasErrors, err := a.formatFile(file, cli)
		if err != nil {
			logger.Error("%v", err)
		}
		if hasErrors {
			failed = true
		}
	}

	if failed {
		return 1
	}

	return 0
}

// formatFile lints, reformats,
// and writes a single page,
// reporting whether the run failed.
func (a *App) formatFile(
	file string,
	cli *cmd.CLI,
) (bool, error) {
	result, err := a.lintFile(file, cli.Ignore...)
	if err != nil {
		return true, err
	}

	for _, e := range result.Errors {
		a.writeLintError(file, e, cli.Tabular)
	}

	root, err := os.OpenRoot(filepath.Dir(file))
	if err != nil {
		return true, err
	}
	defer func() {
		_ = root.Close()
	}()

	content, err := root.ReadFile(filepath.Base(file))
	if err != nil {
		return true, err
	}

	formatted := lint.Format(string(content))
	if formatted == "" {
		_, _ = fmt.Fprintln(
			a.Stderr,
			"refraining from formatting because of a fatal error",
		)
		return true, nil
	}

	err = a.writeFormatted(root, file, formatted, cli)
	return len(result.Errors) > 0 || err != nil, err
}

// writeFormatted emits the formatted page content
// according to the CLI flags.
//
// With --in-place it writes the content back over the source file
// with 0600 permissions, through the caller-supplied root
// for that file's directory.
//
// With --output it instead opens a root anchored at the output file's directory
// and writes filepath.Base of the output path into it,
// so the destination may live in a different directory
// from the source.
//
// With neither flag set it prints the content to a.Stdout
// followed by a newline, mirroring the reference linter's console output.
//
// It returns the underlying write error, if any.
func (a *App) writeFormatted(
	root *os.Root,
	path string,
	content string,
	cli *cmd.CLI,
) error {
	switch {
	case cli.InPlace:
		return root.WriteFile(
			filepath.Base(path),
			[]byte(content),
			0o600,
		)
	case cli.Output != "":
		outputRoot, err := os.OpenRoot(filepath.Dir(cli.Output))
		if err != nil {
			return err
		}
		defer func() {
			_ = outputRoot.Close()
		}()

		return outputRoot.WriteFile(
			filepath.Base(cli.Output),
			[]byte(content),
			0o600,
		)
	default:
		_, err := fmt.Fprint(a.Stdout, content, "\n")
		return err
	}
}
