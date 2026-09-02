package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPageFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "page_in_platform_dir", path: "common/git.md", want: true},
		{name: "nested_page", path: "linux/arch/git.md", want: true},
		{name: "root_level_md", path: "LICENSE.md", want: false},
		{name: "non_md_in_dir", path: "common/notes.txt", want: false},
		{name: "root_non_md", path: "README", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isPageFile(tt.path))
		})
	}
}

func TestHashPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T, dir string)
		want  map[string]string
	}{
		{
			name:  "nonexistent_directory",
			setup: func(t *testing.T, dir string) {},
			want:  map[string]string{},
		},
		{
			name: "empty_directory",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(dir, 0o750))
			},
			want: map[string]string{},
		},
		{
			name: "hashes_nested_pages",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "common"), 0o750))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "common", "git.md"), []byte("# git\n"), 0o600))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "common", "ls.md"), nil, 0o600))
			},
			want: map[string]string{
				filepath.Join("common", "git.md"): contentHash("# git\n"),
				filepath.Join("common", "ls.md"):  contentHash(""),
			},
		},
		{
			name: "ignores_non_pages",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "common"), 0o750))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "LICENSE.md"), []byte("license"), 0o600))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "common", "notes.txt"), []byte("notes"), 0o600))
			},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "subdir")
			tt.setup(t, dir)

			got, err := hashPages(dir)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
