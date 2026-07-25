package cmd

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			name: "multiple",
			cli:  CLI{Update: true, Search: "foo"},
			want: []string{"--update", "--search <KEYWORD>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := activeOps(&tt.cli)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFmtUsage(t *testing.T) {
	t.Parallel()

	err := fmtUsage("test message %d", 42)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test message 42")
	assert.Contains(t, err.Error(), "Usage:")
	assert.Contains(t, err.Error(), "tlgc")
	assert.Contains(t, err.Error(), "--help")
}

func TestFmtFlagErrorUndefined(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("update", false, "")
	fs.Bool("search", false, "")

	err := fs.Parse([]string{"--bogus"})
	assert.Error(t, err)

	err = fmtFlagError(fs, err)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected argument")
	assert.Contains(t, err.Error(), "--bogus")
}

func TestFmtFlagErrorUndefinedTip(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("update", false, "")
	fs.Bool("search", false, "")

	err := fs.Parse([]string{"--searc"})
	assert.Error(t, err)

	err = fmtFlagError(fs, err)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected argument")
	assert.Contains(t, err.Error(), "--searc")
	assert.Contains(t, err.Error(), "similar argument")
	assert.Contains(t, err.Error(), "--search")
}

func TestFmtFlagErrorNeedsArg(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("search", "", "")

	err := fs.Parse([]string{"--search"})
	assert.Error(t, err)

	err = fmtFlagError(fs, err)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires an argument")
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
