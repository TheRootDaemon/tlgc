package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		platform string
		want     []string
		wantErr  bool
	}{
		{
			name: "pages_in_platform_and_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			platform: "linux",
			want:     []string{"apt", "ls"},
		},
		{
			name: "common_platform_returns_only_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			platform: "common",
			want:     []string{"ls"},
		},
		{
			name: "dedupes_duplicate_pages",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				return d
			},
			platform: "linux",
			want:     []string{"git"},
		},
		{
			name: "empty_platform_dir",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform: "linux",
			want:     nil,
		},
		{
			name: "error_when_no_pages_en",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			platform: "linux",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.ListFor(tt.platform)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		want     []string
		wantErr  bool
	}{
		{
			name: "all_pages_across_platforms",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			want: []string{"apt", "ls"},
		},
		{
			name: "dedupes_duplicates",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				return d
			},
			want: []string{"git"},
		},
		{
			name: "error_when_no_pages_en",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.ListAll()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListPlatforms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		want     []string
		wantErr  bool
	}{
		{
			name: "returns_platforms",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			want: []string{"common", "linux"},
		},
		{
			name: "single_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			want: []string{"common"},
		},
		{
			name: "error_when_no_pages_en",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.ListPlatforms()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListLanguages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		want     []string
		wantErr  bool
	}{
		{
			name: "single_language",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			want: []string{"en"},
		},
		{
			name: "multiple_languages_sorted",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				return d
			},
			want: []string{"de", "en"},
		},
		{
			name: "no_pages_dirs",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.ListLanguages()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListDirectory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupDir    func(t *testing.T) string
		platform    string
		languageDir string
		want        []string
		wantErr     bool
	}{
		{
			name: "returns_md_files",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			platform:    "common",
			languageDir: "pages.en",
			want:        []string{"git", "ls"},
		},
		{
			name: "skips_non_md_files",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "notes.txt"), nil, 0o644))
				return d
			},
			platform:    "common",
			languageDir: "pages.en",
			want:        []string{"git"},
		},
		{
			name: "skips_directories",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common", "subdir"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				return d
			},
			platform:    "common",
			languageDir: "pages.en",
			want:        []string{"git"},
		},
		{
			name: "non_existent_dir_returns_nil",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform:    "nonexistent",
			languageDir: "pages.en",
			want:        nil,
		},
		{
			name: "empty_dir_returns_nil",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform:    "common",
			languageDir: "pages.en",
			want:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.listDirectory(tt.platform, tt.languageDir)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
