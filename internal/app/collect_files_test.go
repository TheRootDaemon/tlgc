package app

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(*testing.T) ([]string, []string, bool)
	}{
		{
			name: "file_and_directory_combined",
			setup: func(t *testing.T) ([]string, []string, bool) {
				dir := filepath.Join(t.TempDir(), "pages")
				assert.NoError(t, mkdirall(dir))

				solo := filepath.Join(t.TempDir(), "solo.md")
				md := filepath.Join(dir, "a.md")
				other := filepath.Join(dir, "b.txt")

				touch(t, solo)
				touch(t, md)
				touch(t, other)

				return []string{solo, dir}, []string{solo, md}, false
			},
		},
		{
			name: "nonexistent_aborts",
			setup: func(t *testing.T) ([]string, []string, bool) {
				dir := t.TempDir()
				md := filepath.Join(dir, "a.md")
				touch(t, md)

				return []string{md, filepath.Join(dir, "nope")}, nil, true
			},
		},
		{
			name: "empty_input",
			setup: func(t *testing.T) ([]string, []string, bool) {
				return nil, nil, false
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				paths, want, wantErr := tt.setup(t)

				got, err := collectFiles(paths)

				if wantErr {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)
				sort.Strings(got)
				sort.Strings(want)
				assert.Equal(t, want, got)
			},
		)
	}
}

func TestCollectPathFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(*testing.T) (string, []string, bool)
	}{
		{
			name: "single_md_file",
			setup: func(t *testing.T) (string, []string, bool) {
				path := filepath.Join(t.TempDir(), "a.md")
				touch(t, path)
				return path, []string{path}, false
			},
		},
		{
			name: "non_md_direct_file_included",
			setup: func(t *testing.T) (string, []string, bool) {
				path := filepath.Join(t.TempDir(), "a.txt")
				touch(t, path)
				return path, []string{path}, false
			},
		},
		{
			name: "directory_walked_recursively",
			setup: func(t *testing.T) (string, []string, bool) {
				dir := t.TempDir()
				sub := filepath.Join(dir, "sub")
				assert.NoError(t, mkdirall(sub))

				files := []string{
					filepath.Join(dir, "a.md"),
					filepath.Join(dir, "b.txt"),
					filepath.Join(sub, "c.md"),
				}

				for _, f := range files {
					touch(t, f)
				}
				return dir, []string{files[0], files[2]}, false
			},
		},
		{
			name: "empty_directory",
			setup: func(t *testing.T) (string, []string, bool) {
				return t.TempDir(), nil, false
			},
		},
		{
			name: "nonexistent_path",
			setup: func(t *testing.T) (string, []string, bool) {
				return filepath.Join(t.TempDir(), "nope"), nil, true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, want, wantErr := tt.setup(t)

			got, err := collectPathFiles(path)

			if wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			sort.Strings(got)
			sort.Strings(want)
			assert.Equal(t, want, got)
		})
	}
}

// touch creates an empty file at path with permissions 0600.
// It is intended for use in tests.
func touch(t *testing.T, path string) {
	t.Helper()
	assert.NoError(t, os.WriteFile(path, nil, 0o600))
}

// mkdirall creates path and any missing parent directories
// with permissions 0750.
func mkdirall(path string) error {
	return os.MkdirAll(path, 0o750)
}
