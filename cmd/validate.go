package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/TheRootDaemon/tlgc/termcolor"
	"github.com/TheRootDaemon/tlgc/version"
)

// Validate parses the flag set, validates the CLI, and returns the result.
func Validate(cli *CLI, fs *flag.FlagSet, args []string) (*CLI, error) {
	fs.Usage = func() {}
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return nil, fmtFlagError(fs, err)
	}

	cli.Page = fs.Args()
	cli.HasArgs = len(args) > 0

	if err := validate(cli); err != nil {
		return nil, err
	}

	return cli, nil
}

// validate checks that the parsed CLI has valid flags.
func validate(cli *CLI) error {
	switch cli.Color {
	case "auto", "always", "never":
	default:
		return fmtUsage(
			"invalid value for %s (expected %s, %s, %s)",
			termcolor.Sprint("bold blue", "--color"),
			termcolor.Sprint("blue", "auto"),
			termcolor.Sprint("blue", "always"),
			termcolor.Sprint("blue", "never"),
		)
	}

	// show version
	if cli.ShowVersion {
		fmt.Printf(
			"tlgc %s (implementing client specification v2.3)\n",
			version.String(),
		)
		return nil
	}

	// show help
	if cli.ShowHelp {
		help()
		return nil
	}

	// validate that exactly one operation is active
	ops := cli.operationCount()
	if ops == 0 {
		if cli.HasArgs {
			return fmtUsage("no operation specified")
		}
		help()
		return nil
	}
	if ops > 1 {
		return fmtConflictError(cli)
	}

	// browse requires a page
	if cli.Browse && len(cli.Page) == 0 {
		return fmtUsage(
			"flag %s requires a page",
			termcolor.Sprint("bold blue", "--browse"),
		)
	}

	// lint and format require a file or directory
	if cli.Lint && len(cli.Page) == 0 {
		return fmtUsage(
			"flag %s requires a file or directory",
			termcolor.Sprint("bold blue", "--lint"),
		)
	}
	if cli.Format && len(cli.Page) == 0 {
		return fmtUsage(
			"flag %s requires a file or directory",
			termcolor.Sprint("bold blue", "--format"),
		)
	}

	if cli.Output != "" && !cli.Format {
		return fmtUsage(
			"flag %s requires %s",
			termcolor.Sprint("bold blue", "--output"),
			termcolor.Sprint("bold blue", "--format"),
		)
	}

	return nil
}

// operationCount returns how many operation-group flags are active.
func (c *CLI) operationCount() int {
	count := 0

	if len(c.Page) > 0 &&
		!c.Browse &&
		!c.Lint &&
		!c.Format {
		count++
	}
	if c.Update {
		count++
	}
	if c.List {
		count++
	}
	if c.ListAll {
		count++
	}
	if c.Lint {
		count++
	}
	if c.Format {
		count++
	}
	if c.Search != "" {
		count++
	}
	if c.Browse {
		count++
	}
	if c.ListPlatforms {
		count++
	}
	if c.ListLanguages {
		count++
	}
	if c.Info {
		count++
	}
	if c.Render != "" {
		count++
	}
	if c.CleanCache {
		count++
	}
	if c.GenConfig {
		count++
	}
	if c.ConfigPath {
		count++
	}

	return count
}
