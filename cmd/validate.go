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
		help()
		return nil
	}
	if ops > 1 {
		return fmtConflictError(cli)
	}

	return nil
}

// operationCount returns how many operation-group flags are active.
func (c *CLI) operationCount() int {
	count := 0

	if len(c.Page) > 0 {
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
	if c.Search != "" {
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
