package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/TheRootDaemon/tlgc/version"
)

// Parse parses the process command-line arguments into a CLI value.
func Parse() (*CLI, error) {
	return parse(os.Args[1:])
}

// parse parses the provided command-line arguments into a CLI value.
//
// It validates the parsed flags
// and ensures that exactly one operation has been requested.
// If the arguments are empty it prints help message.
func parse(args []string) (*CLI, error) {
	cli := &CLI{}

	fs := flag.NewFlagSet("tlgc", flag.ContinueOnError)

	// operations
	fs.BoolVar(&cli.Update, "u", false, "update the cache")
	fs.BoolVar(&cli.Update, "update", false, "update the cache")

	fs.BoolVar(
		&cli.List,
		"l",
		false,
		"list all pages for the current platform",
	)
	fs.BoolVar(
		&cli.List,
		"list",
		false,
		"list all pages for the current platform",
	)

	fs.BoolVar(&cli.ListAll, "a", false, "list all pages")
	fs.BoolVar(&cli.ListAll, "list-all", false, "list all pages")

	fs.StringVar(
		&cli.Search,
		"s",
		"",
		"search for pages containing a keyword",
	)
	fs.StringVar(
		&cli.Search,
		"search",
		"",
		"search for pages containing a keyword",
	)

	fs.BoolVar(
		&cli.ListPlatforms,
		"list-platforms",
		false,
		"list available platforms",
	)

	fs.BoolVar(
		&cli.ListLanguages,
		"list-languages",
		false,
		"list installed languages",
	)

	fs.BoolVar(&cli.Info, "i", false, "show cache information")
	fs.BoolVar(&cli.Info, "info", false, "show cache information")

	fs.StringVar(&cli.Render, "r", "", "render the specified tldr page")
	fs.StringVar(&cli.Render, "render", "", "render the specified tldr page")

	fs.BoolVar(
		&cli.CleanCache,
		"clean-cache",
		false,
		"interactively delete the cache contents",
	)

	fs.BoolVar(
		&cli.GenConfig,
		"gen-config",
		false,
		"print the default configuration",
	)

	fs.BoolVar(
		&cli.ConfigPath,
		"config-path",
		false,
		"print the configuration path",
	)

	fs.BoolVar(&cli.ShowVersion, "v", false, "display version")
	fs.BoolVar(&cli.ShowVersion, "version", false, "display version")

	fs.BoolVar(&cli.ShowHelp, "h", false, "display help")
	fs.BoolVar(&cli.ShowHelp, "help", false, "display help")

	// options
	fs.StringVar(
		&cli.Platform,
		"p",
		"",
		"specify the platform to use (linux, osx, windows, etc.)",
	)
	fs.StringVar(
		&cli.Platform,
		"platform",
		"",
		"specify the platform to use (linux, osx, windows, etc.)",
	)

	fs.Var(
		&stringListValue{
			values: &cli.Languages,
		},
		"L",
		"specify the languages to use",
	)
	fs.Var(
		&stringListValue{
			values: &cli.Languages,
		},
		"language",
		"specify the languages to use",
	)

	fs.BoolVar(
		&cli.ShortOptions,
		"short-options",
		false,
		"display short options wherever possible (e.g. '-s')",
	)
	fs.BoolVar(
		&cli.LongOptions,
		"long-options",
		false,
		"display long options wherever possible (e.g. '--long')",
	)

	fs.BoolVar(
		&cli.Edit,
		"edit",
		false,
		"display a link to edit the page on GitHub",
	)

	fs.BoolVar(
		&cli.Offline,
		"o",
		false,
		"do not update the cache, even if it is stale",
	)
	fs.BoolVar(
		&cli.Offline,
		"offline",
		false,
		"do not update the cache, even if it is stale",
	)

	fs.BoolVar(&cli.Compact, "c", false, "strip empty lines from output")
	fs.BoolVar(&cli.Compact, "compact", false, "strip empty lines from output")

	fs.BoolVar(
		&cli.Raw,
		"R",
		false,
		"print pages in raw markdown instead of rendering them",
	)
	fs.BoolVar(
		&cli.Raw,
		"raw",
		false,
		"print pages in raw markdown instead of rendering them",
	)

	fs.BoolVar(&cli.Quiet, "q", false, "suppress status messages and warnings")
	fs.BoolVar(&cli.Quiet, "quiet", false, "suppress status messages and warnings")

	fs.Var(
		&countValue{
			count: &cli.Verbose,
		},
		"verbose",
		"increase verbosity (repeat upto twice)",
	)

	fs.StringVar(
		&cli.Color,
		"color",
		"auto",
		"specify when to enable color (auto, always, never)",
	)

	fs.StringVar(
		&cli.Config,
		"config",
		"",
		"specify an alternative configuration file",
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	switch cli.Color {
	case "auto", "always", "never":
	default:
		return nil, fmt.Errorf("invalid value %q for --color (expected auto, always, never)", cli.Color)
	}

	// show version
	if cli.ShowVersion {
		fmt.Printf(
			"tlgc %s (implementing client specification v2.3)\n",
			version.Version,
		)
		return cli, nil
	}

	// show help
	if cli.ShowHelp {
		fmt.Println("TODO")
		return cli, nil
	}

	// positional arguments
	cli.Page = fs.Args()

	// validate that exactly one operation is active
	ops := cli.operationCount()
	if ops == 0 {
		fmt.Println("TODO")
		return cli, nil
	} else if ops > 1 {
		return nil, fmt.Errorf("only one operation can be specified at a time")
	}

	return cli, nil
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
