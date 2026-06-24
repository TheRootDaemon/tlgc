package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TheRootDaemon/tlgc/internal/config"
)

func TestNew(t *testing.T) {
	t.Run("from_initialized_config", func(t *testing.T) {
		config.ResetForTesting()
		defer config.ResetForTesting()

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.toml")
		err := os.WriteFile(cfgPath, []byte("[cache]\ndir = \"/custom/cache\"\n"), 0o644)
		require.NoError(t, err)

		t.Setenv("TLGC_CONFIG", cfgPath)
		err = config.Initialize()
		require.NoError(t, err)

		c := New()
		assert.Equal(t, "/custom/cache", c.Dir())
	})

	t.Run("default_dir_when_not_in_config", func(t *testing.T) {
		config.ResetForTesting()
		defer config.ResetForTesting()

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.toml")
		err := os.WriteFile(cfgPath, []byte("[output]\nshow_title = false\n"), 0o644)
		require.NoError(t, err)

		t.Setenv("TLGC_CONFIG", cfgPath)
		err = config.Initialize()
		require.NoError(t, err)

		c := New()
		assert.Equal(t, config.Cache().Dir, c.Dir())
	})

	t.Run("empty_dir_in_config", func(t *testing.T) {
		config.ResetForTesting()
		defer config.ResetForTesting()

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.toml")
		err := os.WriteFile(cfgPath, []byte("[cache]\ndir = \"\"\n"), 0o644)
		require.NoError(t, err)

		t.Setenv("TLGC_CONFIG", cfgPath)
		err = config.Initialize()
		require.NoError(t, err)

		c := New()
		assert.Equal(t, "", c.Dir())
	})
}

// TestDir tests Cache.Dir.
func TestDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dir  string
		want string
	}{
		{name: "simple_path", dir: "/tmp/cache", want: "/tmp/cache"},
		{name: "empty_string", dir: "", want: ""},
		{name: "relative_path", dir: "./test/cache", want: "./test/cache"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{dir: tt.dir}
			assert.Equal(t, tt.want, c.Dir())
		})
	}
}

// TestSubDirExists tests Cache.subDirExists.
func TestSubDirExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		subName  string
		want     bool
	}{
		{
			name: "existing_subdirectory",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "sub"), 0o755))
				return d
			},
			subName: "sub",
			want:    true,
		},
		{
			name: "non_existent_name",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			subName: "nonexistent",
			want:    false,
		},
		{
			name: "file_instead_of_dir",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(d, "file.txt"), nil, 0o644))
				return d
			},
			subName: "file.txt",
			want:    false,
		},
		{
			name: "empty_name_returns_true",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			subName: "",
			want:    true,
		},
		{
			name: "nested_path",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "a", "b"), 0o755))
				return d
			},
			subName: "a/b",
			want:    true,
		},
		{
			name: "path_traversal",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			subName: "../etc",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			assert.Equal(t, tt.want, c.subDirExists(tt.subName))
		})
	}
}

// TestGetPlatforms tests Cache.getPlatforms.
func TestGetPlatforms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupDir      func(t *testing.T) string
		want          []string
		wantErr       bool
		skipCacheSeed bool
		extraChecks   func(t *testing.T, c *Cache, dir string)
	}{
		{
			name: "no_pages_en_dir",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
		{
			name: "pages_en_empty_no_subdirs",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			wantErr: true,
		},
		{
			name: "pages_en_only_files",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "file1.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "file2.txt"), nil, 0o644))
				return d
			},
			wantErr: true,
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
			name: "multiple_platforms_sorted",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "osx"), 0o755))
				return d
			},
			want: []string{"common", "linux", "osx"},
		},
		{
			name: "caches_result",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				return d
			},
			skipCacheSeed: true,
			want:          []string{"common", "linux"},
			extraChecks: func(t *testing.T, c *Cache, dir string) {
				got1, err := c.getPlatforms()
				require.NoError(t, err)
				assert.Equal(t, []string{"common", "linux"}, got1)

				require.NoError(t, os.RemoveAll(filepath.Join(dir, "pages.en", "linux")))

				got2, err := c.getPlatforms()
				require.NoError(t, err)
				assert.Equal(t, got1, got2)
			},
		},
		{
			name: "re_read_after_cache_nil",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			skipCacheSeed: true,
			want:          []string{"common", "linux"},
			extraChecks: func(t *testing.T, c *Cache, dir string) {
				got1, err := c.getPlatforms()
				require.NoError(t, err)
				assert.Equal(t, []string{"common"}, got1)

				c.platforms.Store([]string(nil))

				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "linux"), 0o755))

				got2, err := c.getPlatforms()
				require.NoError(t, err)
				assert.Equal(t, []string{"common", "linux"}, got2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}

			if tt.extraChecks != nil {
				tt.extraChecks(t, c, dir)
				return
			}

			got, err := c.getPlatforms()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestGetLanguageDirectories tests Cache.getLanguageDirectories.
func TestGetLanguageDirectories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		want     []string
		wantErr  bool
	}{
		{
			name: "empty_cache_dir",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			want: nil,
		},
		{
			name: "only_non_pages_dirs",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "other"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "data"), 0o755))
				return d
			},
			want: nil,
		},
		{
			name: "single_pages_en",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			want: []string{"pages.en"},
		},
		{
			name: "multiple_pages_dirs",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.zh"), 0o755))
				return d
			},
			want: []string{"pages.de", "pages.en", "pages.zh"},
		},
		{
			name: "mixed_pages_and_other",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "other"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "data"), 0o755))
				return d
			},
			want: []string{"pages.de", "pages.en"},
		},
		{
			name: "handles_readdir_error",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				filePath := filepath.Join(d, "not_a_dir")
				require.NoError(t, os.WriteFile(filePath, nil, 0o644))
				return filePath
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}

			got, err := c.getLanguageDirectories()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestLanguagesToDirectories tests Cache.languagesToDirectories.
func TestLanguagesToDirectories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupDir  func(t *testing.T) string
		languages []string
		sortFlag  bool
		want      []string
	}{
		{
			name: "nil_languages",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			languages: nil,
			sortFlag:  false,
			want:      nil,
		},
		{
			name: "empty_languages",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			languages: []string{},
			sortFlag:  false,
			want:      nil,
		},
		{
			name: "single_language_exists",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			languages: []string{"en"},
			sortFlag:  false,
			want:      []string{"pages.en"},
		},
		{
			name: "single_language_does_not_exist",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				return d
			},
			languages: []string{"de"},
			sortFlag:  false,
			want:      nil,
		},
		{
			name: "multiple_languages_all_exist",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				return d
			},
			languages: []string{"en", "de"},
			sortFlag:  false,
			want:      []string{"pages.en", "pages.de"},
		},
		{
			name: "multiple_languages_some_missing",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.fr"), 0o755))
				return d
			},
			languages: []string{"en", "de", "fr"},
			sortFlag:  false,
			want:      []string{"pages.en", "pages.fr"},
		},
		{
			name: "duplicates_in_input",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				return d
			},
			languages: []string{"en", "de", "en"},
			sortFlag:  false,
			want:      []string{"pages.en", "pages.de"},
		},
		{
			name: "sort_flag_true",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.fr"), 0o755))
				return d
			},
			languages: []string{"fr", "en", "de"},
			sortFlag:  true,
			want:      []string{"pages.de", "pages.en", "pages.fr"},
		},
		{
			name: "sort_flag_true_with_duplicates",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.fr"), 0o755))
				return d
			},
			languages: []string{"fr", "en", "de", "en"},
			sortFlag:  true,
			want:      []string{"pages.de", "pages.en", "pages.fr"},
		},
		{
			name: "sort_flag_false_preserves_order",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.fr"), 0o755))
				return d
			},
			languages: []string{"fr", "en", "de"},
			sortFlag:  false,
			want:      []string{"pages.fr", "pages.en", "pages.de"},
		},
		{
			name: "all_languages_missing",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			languages: []string{"en", "de"},
			sortFlag:  false,
			want:      nil,
		},
		{
			name: "sort_with_missing",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.zh"), 0o755))
				return d
			},
			languages: []string{"zh", "de", "en"},
			sortFlag:  true,
			want:      []string{"pages.en", "pages.zh"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got := c.languagesToDirectories(tt.languages, tt.sortFlag)
			assert.Equal(t, tt.want, got)
		})
	}
}
