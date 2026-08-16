package cmd

import (
	"strings"

	"github.com/TheRootDaemon/tlgc/termcolor"
)

// flagDependency describes a modifier flag
// and the operation flags that it depends on.
type flagDependency struct {
	flag    string          // flag is the modifier flag being validated.
	parents []string        // parents contains the names of the operations that permit the modifier.
	present func(*CLI) bool // present reports whether the modifier flag is enabled in the CLI configuration.
	valid   func(*CLI) bool // valid reports whether the modifier can be used with the current CLI configuration.
}

// modifierDependencies contains the dependency rules for modifier flags.
//
// Each entry describes which operation
// or operations must be active for a particular modifier flag to be valid.
// These rules are evaluated by validateFlagDependencies.
var modifierDependencies = []flagDependency{
	{
		flag:    "--output",
		parents: []string{"--format"},
		present: func(c *CLI) bool { return c.Output != "" },
		valid:   func(c *CLI) bool { return c.Format },
	},
	{
		flag:    "--in-place",
		parents: []string{"--format"},
		present: func(c *CLI) bool { return c.InPlace },
		valid:   func(c *CLI) bool { return c.Format },
	},
	{
		flag:    "--tabular",
		parents: []string{"--lint", "--format"},
		present: func(c *CLI) bool { return c.Tabular },
		valid:   func(c *CLI) bool { return c.Lint || c.Format },
	},
	{
		flag:    "--ignore",
		parents: []string{"--lint", "--format"},
		present: func(c *CLI) bool { return len(c.Ignore) > 0 },
		valid:   func(c *CLI) bool { return c.Lint || c.Format },
	},
	{
		flag:    "--platform",
		parents: []string{"a page", "--browse", "--list", "--search"},
		present: func(c *CLI) bool { return c.Platform != "" },
		valid: func(c *CLI) bool {
			return c.pageLookup() || c.Browse || c.List || c.Search != ""
		},
	},
	{
		flag:    "--language",
		parents: []string{"a page", "--browse", "--search", "--update"},
		present: func(c *CLI) bool { return len(c.Languages) > 0 },
		valid: func(c *CLI) bool {
			return c.pageLookup() || c.Browse || c.Search != "" || c.Update
		},
	},
	{
		flag:    "--offline",
		parents: []string{"a page", "--browse"},
		present: func(c *CLI) bool { return c.Offline },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Browse },
	},
	{
		flag:    "--compact",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.Compact },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
	{
		flag:    "--no-compact",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.NoCompact },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
	{
		flag:    "--raw",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.Raw },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
	{
		flag:    "--no-raw",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.NoRaw },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
	{
		flag:    "--short-options",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.ShortOptions },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
	{
		flag:    "--long-options",
		parents: []string{"a page", "--render"},
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
		present: func(c *CLI) bool { return c.LongOptions },
	},
	{
		flag:    "--edit",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.Edit },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
	{
		flag:    "--color",
		parents: []string{"a page", "--render"},
		present: func(c *CLI) bool { return c.Color != "auto" },
		valid:   func(c *CLI) bool { return c.pageLookup() || c.Render != "" },
	},
}

// validateFlagDependencies validates the dependencies
// between modifier flags and their parent operations.
//
// A modifier flag is valid only when at least one of the operations
// listed in its corresponding flagDependency is active.
// If a modifier is used without a valid parent operation,
// validateFlagDependencies returns a usage error
// describing the required operation.
func validateFlagDependencies(cli *CLI) error {
	for _, dependency := range modifierDependencies {
		if !dependency.present(cli) || dependency.valid(cli) {
			continue
		}
		return fmtUsage(
			"flag %s requires %s",
			termcolor.Sprint("bold blue", dependency.flag),
			formatParents(dependency.parents),
		)
	}

	return nil
}

// pageLookup reports whether the CLI represents a bare page lookup.
//
// A bare page lookup occurs when at least one page argument is provided
// and none of the explicit page-processing operations are active
// that is, --browse, --lint, or --format are active.
func (c *CLI) pageLookup() bool {
	return len(c.Page) > 0 && !c.Browse && !c.Lint && !c.Format
}

// formatParents formats the parent flags as a human-readable list
// separated by "or", such as "X or Y" or "X, Y, or Z".
func formatParents(items []string) string {
	if len(items) == 1 {
		return termcolor.Sprint("bold blue", items[0])
	}

	styled := make([]string, len(items))
	for i, item := range items {
		styled[i] = termcolor.Sprint("bold blue", item)
	}

	head := strings.Join(styled[:len(styled)-1], ", ")
	return head + " or " + styled[len(styled)-1]
}
