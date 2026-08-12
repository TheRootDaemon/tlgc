package cmd

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtFlagError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*flag.FlagSet)
		args     []string
		contains []string
	}{
		{
			name: "undefined",
			setup: func(fs *flag.FlagSet) {
				fs.Bool("update", false, "")
				fs.Bool("search", false, "")
			},
			args: []string{"--bogus"},
			contains: []string{
				"unexpected argument",
			},
		},
		{
			name: "undefined_with_tip",
			setup: func(fs *flag.FlagSet) {
				fs.Bool("update", false, "")
				fs.Bool("search", false, "")
			},
			args: []string{"--searc"},
			contains: []string{
				"unexpected argument",
			},
		},
		{
			name: "missing_argument",
			setup: func(fs *flag.FlagSet) {
				fs.String("search", "", "")
			},
			args: []string{"--search"},
			contains: []string{
				"requires an argument",
			},
		},
		{
			name: "missing_argument_output",
			setup: func(fs *flag.FlagSet) {
				fs.String("output", "", "")
			},
			args: []string{"--output"},
			contains: []string{
				"--format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			tt.setup(fs)

			err := fs.Parse(tt.args)
			assert.Error(t, err)

			err = fmtFlagError(fs, err)
			assert.Error(t, err)

			for _, s := range tt.contains {
				assert.True(t, strings.Contains(err.Error(), s))
			}
		})
	}
}

func TestFmtConflictError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		cli      CLI
		contains []string
	}{
		{
			name:     "two_operations",
			cli:      CLI{Update: true, List: true},
			contains: []string{"cannot be used with", "--update", "--list"},
		},
		{
			name:     "page_and_search",
			cli:      CLI{Page: []string{"tar"}, Search: "foo"},
			contains: []string{"cannot be used with", "[PAGE]...", "--search"},
		},
		{
			name:     "lint_and_format",
			cli:      CLI{Lint: true, Format: true, Page: []string{"file.md"}},
			contains: []string{"cannot be used with", "--lint <FILE|DIR>", "--format <FILE|DIR>"},
		},
		{
			name:     "one_operation_fallback",
			cli:      CLI{Update: true},
			contains: []string{"only one operation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fmtConflictError(&tt.cli)
			assert.Error(t, err)
			for _, s := range tt.contains {
				assert.True(t, strings.Contains(err.Error(), s),
					"error %q should contain %q", err.Error(), s)
			}
		})
	}
}

func TestFmtUsage(t *testing.T) {
	t.Parallel()

	err := fmtUsage("test message %d", 42)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test message 42")
	assert.Contains(t, err.Error(), "Usage:")
	assert.Contains(t, err.Error(), "tldr")
	assert.Contains(t, err.Error(), "--help")
}

func TestFlagDisplay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "short", arg: "x", want: "-x"},
		{name: "long", arg: "update", want: "--update"},
		{name: "single_char_short", arg: "u", want: "-u"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flagDisplay(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestActiveOps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cli  CLI
		want []string
	}{
		{
			name: "empty",
			cli:  CLI{},
			want: nil,
		},
		{
			name: "page",
			cli:  CLI{Page: []string{"tar"}},
			want: []string{"[PAGE]..."},
		},
		{
			name: "update",
			cli:  CLI{Update: true},
			want: []string{"--update"},
		},
		{
			name: "browse",
			cli:  CLI{Browse: true, Page: []string{"tar"}},
			want: []string{"--browse [PAGE]..."},
		},
		{
			name: "search",
			cli:  CLI{Search: "ngi"},
			want: []string{"--search <KEYWORD>"},
		},
		{
			name: "render",
			cli:  CLI{Render: "file.md"},
			want: []string{"--render <FILE>"},
		},
		{
			name: "lint",
			cli:  CLI{Lint: true, Page: []string{"file.md"}},
			want: []string{"--lint <FILE|DIR>"},
		},
		{
			name: "format",
			cli:  CLI{Format: true, Page: []string{"file.md"}},
			want: []string{"--format <FILE|DIR>"},
		},
		{
			name: "multiple",
			cli:  CLI{Update: true, Search: "foo"},
			want: []string{"--update", "--search <KEYWORD>"},
		},
		{
			name: "lint_with_update",
			cli:  CLI{Lint: true, Page: []string{"file.md"}, Update: true},
			want: []string{"--update", "--lint <FILE|DIR>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := activeOps(&tt.cli)
			assert.Equal(t, tt.want, got)
		})
	}
}
