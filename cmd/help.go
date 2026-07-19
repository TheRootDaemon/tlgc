package cmd

import (
	"fmt"
	"strings"

	"github.com/TheRootDaemon/tlgc/termcolor"
	"github.com/TheRootDaemon/tlgc/version"
)

func help() {
	printUsage()
	printFlags()
	printFooter()
}

func printUsage() {
	fmt.Printf(
		"tlgc %s (implementing client specification v2.3)\n\n",
		version.String(),
	)
	fmt.Printf(
		"%s tldr [OPTIONS] [PAGE]...\n\n",
		termcolor.Sprint("bold underline", "Usage:"),
	)
	fmt.Printf(
		"%s\n",
		termcolor.Sprint("bold underline", "Arguments:"),
	)
	fmt.Printf("  [PAGE]...  The tldr page to show\n\n")
}

func printFlags() {
	type flagEntry struct {
		short       string
		long        string
		arg         string
		description string
	}

	flags := []flagEntry{
		{
			short:       "-u",
			long:        "--update",
			description: "Update the cache",
		},
		{
			short:       "-l",
			long:        "--list",
			description: "List all pages in the current platform",
		},
		{
			short:       "-a",
			long:        "--list-all",
			description: "List all pages",
		},
		{
			short: "-s",
			long:  "--search",
			arg:   "<KEYWORD>", description: "Search for pages containing a keyword",
		},
		{
			long:        "--list-platforms",
			description: "List available platforms",
		},
		{
			long:        "--list-languages",
			description: "List installed languages",
		},
		{
			short:       "-i",
			long:        "--info",
			description: "Show cache information",
		},
		{
			short: "-r",
			long:  "--render",
			arg:   "<FILE>", description: "Render the specified tldr page",
		},
		{
			long:        "--clean-cache",
			description: "Interactively delete contents of the cache directory",
		},

		{
			long:        "--gen-config",
			description: "Print the default config",
		},

		{
			long:        "--config-path",
			description: "Print the default config path",
		},
		{
			short: "-p",
			long:  "--platform",
			arg:   "<PLATFORM>", description: "Specify the platform to use (linux, osx, windows, etc.)",
		},
		{
			short: "-L",
			long:  "--language",
			arg:   "<LANGUAGE_CODE>", description: "Specify the languages to use",
		},
		{
			long:        "--short-options",
			description: "Display short options wherever possible (e.g. '-s')",
		},

		{
			long:        "--long-options",
			description: "Display long options wherever possible (e.g. '--long')",
		},

		{
			long:        "--edit",
			description: "Display a link to edit the shown page on GitHub",
		},

		{
			short:       "-o",
			long:        "--offline",
			description: "Do not update the cache, even if it is stale",
		},
		{
			short:       "-c",
			long:        "--compact",
			description: "Strip empty lines from output",
		},
		{
			long:        "--no-compact",
			description: "Do not strip empty lines from output (overrides --compact)",
		},
		{
			short:       "-R",
			long:        "--raw",
			description: "Print pages in raw markdown instead of rendering them",
		},
		{long: "--no-raw", description: "Render pages instead of printing raw file contents (overrides --raw)"},
		{
			short:       "-q",
			long:        "--quiet",
			description: "Suppress status messages and warnings",
		},
		{long: "--verbose...", description: "Be more verbose (can be specified twice)"},
		{
			long:        "--color",
			arg:         "<WHEN>",
			description: "Specify when to enable color [default: auto] [possible values: auto, always, never]",
		},
		{
			long:        "--config",
			arg:         "<FILE>",
			description: "Specify an alternative path to the config file",
		},
		{
			short:       "-v",
			long:        "--version",
			description: "Print version",
		},
		{
			short:       "-h",
			long:        "--help",
			description: "Print help",
		},
	}

	maxShort, maxLong := 0, 0
	for _, f := range flags {
		updateColumnWidths(
			&maxShort,
			&maxLong,
			f.short,
			f.long,
			f.arg,
		)
	}

	fmt.Printf(
		"%s\n",
		termcolor.Sprint("bold underline", "Options:"),
	)

	for _, f := range flags {
		printFlag(
			maxShort,
			maxLong,
			f.short,
			f.long,
			f.arg,
			f.description,
		)
	}
}

func printFooter() {
	fmt.Printf("\nSee https://github.com/TheRootDaemon/tlgc for more information.\n")
}

func updateColumnWidths(
	maxShort,
	maxLong *int,
	short, long, arg string,
) {
	shortWidth := len(short)
	if short != "" && long != "" {
		shortWidth++
	}

	*maxShort = max(*maxShort, shortWidth)

	if long == "" {
		return
	}

	longWidth := len(long) + 1
	if arg != "" {
		longWidth += len(arg)
	}

	*maxLong = max(*maxLong, longWidth)
}

func printFlag(
	maxShort,
	maxLong int,
	short, long, arg, description string,
) {
	shortText := short
	if shortText != "" && long != "" {
		shortText += ","
	}

	longText := long
	if arg != "" {
		longText += " " + arg
	}

	longDisplay := termcolor.Sprint("bold", long)
	if arg != "" {
		longDisplay += " " + arg
	}

	fmt.Printf(
		"  %s%s %s%s %s\n",
		termcolor.Sprint("bold", shortText),
		strings.Repeat(" ", maxShort-len(shortText)),
		longDisplay,
		strings.Repeat(" ", maxLong-len(longText)),
		" "+description,
	)
}
