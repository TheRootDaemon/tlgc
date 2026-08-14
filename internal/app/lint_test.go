package app

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/lint"
	"github.com/stretchr/testify/assert"
)

func TestLintPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(*testing.T) ([]string, *cmd.CLI)
		wantCode   int
		wantStderr []string
	}{
		{
			name: "valid_file_returns_zero",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "tar.md", validPage)
				return []string{path}, &cmd.CLI{}
			},
			wantCode: 0,
		},
		{
			name: "invalid_file_reports_to_stderr",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "bad.md", "x tar\n> no\n")
				return []string{path}, &cmd.CLI{}
			},
			wantCode:   1,
			wantStderr: []string{"bad.md"},
		},
		{
			name: "directory_walk_skips_non_md",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				dir := t.TempDir()
				writePage(t, dir, "bad.txt", "x tar\n> no\n")
				writePage(t, dir, "good.md", validPage)
				return []string{dir}, &cmd.CLI{}
			},
			wantCode: 0,
		},
		{
			name: "non_md_direct_file_still_linted",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "page.txt", validPage)
				return []string{path}, &cmd.CLI{}
			},
			wantCode:   1,
			wantStderr: []string{"page.txt"},
		},
		{
			name: "ignore_codes_suppress_errors",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "bad.md", "x tar\n> no\n")
				return []string{path}, &cmd.CLI{Ignore: []string{"TLDR106"}}
			},
			wantCode: 0,
		},
		{
			name: "nonexistent_path_returns_one",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				return []string{filepath.Join(t.TempDir(), "nope")}, &cmd.CLI{}
			},
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, cli := tt.setup(t)
			a, _, stderr := newTestApp()
			cli.Page = paths

			got := a.lintPages(cli)

			assert.Equal(t, tt.wantCode, got)
			out := stderr.String()
			for _, want := range tt.wantStderr {
				assert.Contains(t, out, want)
			}

			if len(tt.wantStderr) == 0 {
				assert.Empty(t, out, strings.TrimSpace(out))
			}
		})
	}
}

func TestWriteLintError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	a := &App{
		Stderr: &buf,
	}

	a.writeLintError(
		"pages/a.md",
		lint.Error{Code: "TLDR001", Line: 1, Description: "leading whitespace"},
	)

	assert.Equal(t, "pages/a.md:1: TLDR001 leading whitespace\n", buf.String())
}

func TestWriteTabular(t *testing.T) {
	t.Parallel()

	a, _, stderr := newTestApp()

	a.writeTabular(
		[]lintViolation{
			{
				path: "a.md",
				err:  lint.Error{Line: 1, Code: "TLDR001", Description: "leading whitespace"},
			},
			{
				path: "b.md",
				err:  lint.Error{Line: 2, Code: "TLDR002", Description: "space"},
			},
		},
	)

	assert.Equal(
		t,
		"File Line Code    Description\n"+
			"a.md 1    TLDR001 leading whitespace\n"+
			"b.md 2    TLDR002 space\n",
		stderr.String(),
	)
}

func TestWriteTabularNoRows(t *testing.T) {
	t.Parallel()

	a, _, stderr := newTestApp()

	a.writeTabular(nil)

	assert.Empty(t, stderr.String())
}

func TestLintColumnWidths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rows     []lintViolation
		wantFile int
		wantLine int
		wantCode int
	}{
		{
			name:     "empty_rows_uses_header_widths",
			rows:     nil,
			wantFile: len("File"),
			wantLine: len("Line"),
			wantCode: len("Code"),
		},
		{
			name: "header_wins_over_short_values",
			rows: []lintViolation{
				{
					path: "a.md",
					err:  lint.Error{Line: 1, Code: "TLDR001", Description: "x"},
				},
			},
			wantFile: len("File"),
			wantLine: len("Line"),
			wantCode: len("TLDR001"),
		},
		{
			name: "long_values_widen_columns",
			rows: []lintViolation{
				{
					path: "pages/somepage.md",
					err:  lint.Error{Line: 123, Code: "TLDR102", Description: "some very long description"},
				},
			},
			wantFile: len("pages/somepage.md"),
			wantLine: len("Line"),
			wantCode: len("TLDR102"),
		},
		{
			name: "picks_max_across_all_rows",
			rows: []lintViolation{
				{
					path: "a.md",
					err:  lint.Error{Line: 1, Code: "TLDR001", Description: "short"},
				},
				{
					path: "pages/very/long/path.md",
					err:  lint.Error{Line: 4567, Code: "TLDR104", Description: "a much longer description than before"},
				},
			},
			wantFile: len("pages/very/long/path.md"),
			wantLine: len("4567"),
			wantCode: len("TLDR104"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile, gotLine, gotCode := lintColumnWidths(tt.rows)
			assert.Equal(t, tt.wantFile, gotFile)
			assert.Equal(t, tt.wantLine, gotLine)
			assert.Equal(t, tt.wantCode, gotCode)
		})
	}
}
