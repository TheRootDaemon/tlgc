package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/TheRootDaemon/tlgc/termcolor"
)

// fmtFlagError wraps flag parse errors with clap-style formatting.
func fmtFlagError(fs *flag.FlagSet, err error) error {
	s := err.Error()

	switch {
	case strings.HasPrefix(s, "flag provided but not defined: "):
		raw := strings.TrimPrefix(s, "flag provided but not defined: ")
		name := strings.TrimLeft(raw, "-")
		var tip string
		if sim := similarFlag(fs, name); sim != "" {
			tip = fmt.Sprintf(
				"\n\n  %s a similar argument exists: %s",
				termcolor.Sprint("bold green", "tip:"),
				termcolor.Sprint("bold blue", flagDisplay(sim)),
			)
		}
		return fmtUsage(
			"unexpected argument %s found%s",
			termcolor.Sprint("bold cyan", flagDisplay(name)),
			tip,
		)
	case strings.HasPrefix(s, "flag needs an argument: "):
		raw := strings.TrimPrefix(s, "flag needs an argument: ")
		name := strings.TrimLeft(raw, "-")
		return fmtUsage("flag %s requires an argument", termcolor.Sprint("bold blue", flagDisplay(name)))
	default:
		return fmtUsage("%s", s)
	}
}

// fmtConflictError builds a clap-style error for conflicting operations.
func fmtConflictError(cli *CLI) error {
	ops := activeOps(cli)
	if len(ops) < 2 {
		return fmt.Errorf("only one operation can be specified at a time")
	}
	return fmtUsage(
		"argument %s cannot be used with %s",
		termcolor.Sprint("blue", ops[0]),
		termcolor.Sprint("blue", ops[1]),
	)
}

// fmtUsage wraps a formatted message with the standard usage footer.
func fmtUsage(format string, args ...any) error {
	usage := fmt.Sprintf(
		"\n\n%s %s [OPTIONS] [PAGE]...\n\nFor more information, try %s.",
		termcolor.Sprint("bold underline", "Usage:"),
		termcolor.Sprint("bold", "tlgc"),
		termcolor.Sprint("bold blue", "--help"),
	)
	return fmt.Errorf(format+usage, args...)
}

// flagDisplay returns the display form of a flag name,
// "-x" for short, "--xxx" for long.
func flagDisplay(name string) string {
	if len(name) == 1 {
		return "-" + name
	}
	return "--" + name
}

// activeOps returns display names for all active operations in cli.
func activeOps(cli *CLI) []string {
	var ops []string
	if len(cli.Page) > 0 {
		ops = append(ops, "[PAGE]...")
	}
	if cli.Update {
		ops = append(ops, "--update")
	}
	if cli.List {
		ops = append(ops, "--list")
	}
	if cli.ListAll {
		ops = append(ops, "--list-all")
	}
	if cli.Search != "" {
		ops = append(ops, "--search <KEYWORD>")
	}
	if cli.Browse {
		ops = append(ops, "--browse")
	}
	if cli.ListPlatforms {
		ops = append(ops, "--list-platforms")
	}
	if cli.ListLanguages {
		ops = append(ops, "--list-languages")
	}
	if cli.Info {
		ops = append(ops, "--info")
	}
	if cli.Render != "" {
		ops = append(ops, "--render <FILE>")
	}
	if cli.CleanCache {
		ops = append(ops, "--clean-cache")
	}
	if cli.GenConfig {
		ops = append(ops, "--gen-config")
	}
	if cli.ConfigPath {
		ops = append(ops, "--config-path")
	}
	return ops
}
