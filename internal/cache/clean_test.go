package cache

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		input    string
		wantErr  bool
		check    func(t *testing.T, dir string)
	}{
		{
			name: "cache_does_not_exist",
			setupDir: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			input: "",
			check: func(t *testing.T, dir string) {},
		},
		{
			name: "cache_is_empty",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			input: "",
			check: func(t *testing.T, dir string) {},
		},
		{
			name: "user_confirms_with_y",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			input: "y\n",
			check: func(t *testing.T, dir string) {
				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Empty(t, entries)
			},
		},
		{
			name: "user_confirms_with_Y",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			input: "Y\n",
			check: func(t *testing.T, dir string) {
				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Empty(t, entries)
			},
		},
		{
			name: "user_declines_with_n",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			input: "n\n",
			check: func(t *testing.T, dir string) {
				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Len(t, entries, 1)
			},
		},
		{
			name: "user_declines_with_N",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			input: "N\n",
			check: func(t *testing.T, dir string) {
				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Len(t, entries, 1)
			},
		},
		{
			name: "multiple_entries_removed",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				return d
			},
			input: "y\n",
			check: func(t *testing.T, dir string) {
				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Empty(t, entries)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}

			err := c.Clean(strings.NewReader(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.check(t, dir)
		})
	}
}

func TestGetEntries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantLen int
		wantNil bool
		wantErr bool
	}{
		{
			name: "directory_with_entries",
			setup: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(d, "a.txt"), nil, 0o644))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "sub"), 0o755))
				return d
			},
			wantLen: 2,
		},
		{
			name: "directory_does_not_exist",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			wantNil: true,
		},
		{
			name: "empty_directory",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantNil: true,
		},
		{
			name: "path_is_a_file",
			setup: func(t *testing.T) string {
				d := t.TempDir()
				f := filepath.Join(d, "file.txt")
				require.NoError(t, os.WriteFile(f, nil, 0o644))
				return f
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			entries, err := getEntries(path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, entries)
				return
			}
			assert.Len(t, entries, tt.wantLen)
		})
	}
}

func TestParseInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty_input_means_yes",
			input: "\n",
			want:  true,
		},
		{
			name:  "lowercase_y",
			input: "y\n",
			want:  true,
		},
		{
			name:  "uppercase_Y",
			input: "Y\n",
			want:  true,
		},
		{
			name:  "full_yes",
			input: "yes\n",
			want:  true,
		},
		{
			name:  "lowercase_n",
			input: "n\n",
			want:  false,
		},
		{
			name:  "uppercase_N",
			input: "N\n",
			want:  false,
		},
		{
			name:  "full_no",
			input: "no\n",
			want:  false,
		},
		{
			name:  "arbitrary_text",
			input: "garbage\n",
			want:  false,
		},
		{
			name:  "whitespace_only",
			input: "  \n",
			want:  true,
		},
		{
			name:  "trailing_spaces",
			input: "  yes  \n",
			want:  true,
		},
		{
			name:  "starts_with_y",
			input: "yellow\n",
			want:  true,
		},
		{
			name:  "starts_with_Y",
			input: "Yellow\n",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseInput(bufio.NewReader(strings.NewReader(tt.input)))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInput_ReaderError(t *testing.T) {
	t.Parallel()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	_ = w.Close()

	got := parseInput(bufio.NewReader(r))
	assert.False(t, got)
}
