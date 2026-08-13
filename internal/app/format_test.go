package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(*testing.T) ([]string, *cmd.CLI)
		wantCode   int
		wantStdout *string
		wantStderr []string
		check      func(*testing.T, string, *cmd.CLI, string)
	}{
		{
			name: "valid_file_prints_to_stdout",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "tar.md", validPage)
				return []string{path}, &cmd.CLI{}
			},
			wantCode:   0,
			wantStdout: new(validPage + "\n"),
		},
		{
			name: "invalid_file_still_formats_with_errors",
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "bad.md", "x tar\n> no\n")
				return []string{path}, &cmd.CLI{}
			},
			wantCode:   1,
			wantStderr: []string{"bad.md"},
			check: func(t *testing.T, _ string, _ *cmd.CLI, stdout string) {
				assert.NotEmpty(t, stdout)
			},
		},
		{
			name:       "in_place_rewrites_file",
			wantStdout: new(""),
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "bad.md", "x tar\n> no\n")
				return []string{path}, &cmd.CLI{InPlace: true}
			},
			wantCode: 1,
			check: func(t *testing.T, path string, _ *cmd.CLI, _ string) {
				got, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.NotEqual(t, "x tar\n> no\n", string(got))
				assert.NotEmpty(t, string(got))
			},
		},
		{
			name:       "output_flag_writes_file",
			wantStdout: new(""),
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				dir := t.TempDir()
				path := writePage(t, dir, "tar.md", validPage)
				return []string{path}, &cmd.CLI{Output: filepath.Join(dir, "out.md")}
			},
			wantCode: 0,
			check: func(t *testing.T, _ string, cli *cmd.CLI, _ string) {
				got, err := os.ReadFile(cli.Output)
				require.NoError(t, err)
				assert.NotEmpty(t, string(got))
			},
		},
		{
			name:       "output_flag_to_different_directory",
			wantStdout: new(""),
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				src := t.TempDir()
				dst := filepath.Join(t.TempDir(), "nested")
				require.NoError(t, os.MkdirAll(dst, 0o750))
				path := writePage(t, src, "tar.md", validPage)
				return []string{path}, &cmd.CLI{Output: filepath.Join(dst, "out.md")}
			},
			wantCode: 0,
			check: func(t *testing.T, _ string, cli *cmd.CLI, _ string) {
				got, err := os.ReadFile(cli.Output)
				require.NoError(t, err)
				assert.NotEmpty(t, string(got))
			},
		},
		{
			name:       "output_flag_requires_single_file",
			wantStdout: new(""),
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				dir := t.TempDir()
				p1 := writePage(t, dir, "a.md", validPage)
				p2 := writePage(t, dir, "b.md", validPage)
				return []string{p1, p2}, &cmd.CLI{Output: filepath.Join(dir, "out.md")}
			},
			wantCode: 1,
		},
		{
			name:       "unparseable_content_refrains",
			wantStdout: new(""),
			setup: func(t *testing.T) ([]string, *cmd.CLI) {
				path := writePage(t, t.TempDir(), "tar.md", "   \n")
				return []string{path}, &cmd.CLI{}
			},
			wantCode:   1,
			wantStderr: []string{"refraining from formatting"},
			check: func(t *testing.T, path string, _ *cmd.CLI, _ string) {
				got, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, "   \n", string(got))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, cli := tt.setup(t)
			a, stdout, stderr := newTestApp()
			cli.Page = paths

			got := a.formatPages(cli)

			assert.Equal(t, tt.wantCode, got)
			if tt.wantStdout != nil {
				assert.Equal(t, *tt.wantStdout, stdout.String())
			}

			for _, want := range tt.wantStderr {
				assert.Contains(t, stderr.String(), want)
			}

			if tt.check != nil {
				tt.check(t, paths[0], cli, stdout.String())
			}
		})
	}
}
