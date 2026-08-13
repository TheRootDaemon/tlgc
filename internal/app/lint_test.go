package app

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/cmd"
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
