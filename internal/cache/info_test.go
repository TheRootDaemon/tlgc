package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAge(t *testing.T) {
	t.Parallel()

	t.Run("uses_checksum_file_mtime", func(t *testing.T) {
		dir := t.TempDir()
		sumPath := filepath.Join(dir, checksumFile)
		err := os.WriteFile(sumPath, []byte("sums"), 0o644)
		require.NoError(t, err)
		err = os.Chtimes(
			sumPath,
			time.Now().Add(-1*time.Hour),
			time.Now().Add(-1*time.Hour),
		)
		require.NoError(t, err)

		c := &Cache{dir: dir}
		age, err := c.Age()
		require.NoError(t, err)
		assert.Greater(t, age, 55*time.Minute)
		assert.Less(t, age, 65*time.Minute)
	})

	t.Run("falls_back_to_cache_dir_mtime", func(t *testing.T) {
		dir := t.TempDir()
		err := os.Chtimes(
			dir,
			time.Now().Add(-2*time.Hour),
			time.Now().Add(-2*time.Hour),
		)
		require.NoError(t, err)

		c := &Cache{dir: dir}
		age, err := c.Age()
		require.NoError(t, err)
		assert.Greater(t, age, 115*time.Minute)
		assert.Less(t, age, 125*time.Minute)
	})

	t.Run("error_on_non_existent_dir", func(t *testing.T) {
		c := &Cache{dir: "/nonexistent/path"}
		_, err := c.Age()
		assert.Error(t, err)
	})

	t.Run("error_on_future_mtime", func(t *testing.T) {
		dir := t.TempDir()
		future := time.Now().Add(1 * time.Hour)
		err := os.Chtimes(dir, future, future)
		require.NoError(t, err)

		c := &Cache{dir: dir}
		_, err = c.Age()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "future")
	})
}

func TestInfo(t *testing.T) {
	t.Run("returns_info_for_valid_cache", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "linux"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "linux", "apt.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "linux", "pacman.md"), nil, 0o644))

		past := time.Now().Add(-1 * time.Hour)
		require.NoError(t, os.Chtimes(dir, past, past))

		c := &Cache{dir: dir}
		info, err := c.Info()
		require.NoError(t, err)
		assert.Equal(t, dir, info.CacheDir)
		assert.Equal(t, 3, info.TotalPages)
		assert.Len(t, info.LanguageStats, 1)
		assert.Equal(t, "en", info.LanguageStats[0].Language)
		assert.Equal(t, 3, info.LanguageStats[0].Pages)
		assert.NotEmpty(t, info.Age)
		assert.True(t, info.AutoUpdate)
		assert.Equal(t, uint64(336), info.MaxAge)
		assert.NotEmpty(t, info.Mirror)
		assert.NotEmpty(t, info.Platforms)
		assert.Contains(t, info.Platforms, "common")
		assert.Contains(t, info.Platforms, "linux")
		assert.Greater(t, info.AgeDuration, 55*time.Minute)
		assert.Less(t, info.AgeDuration, 65*time.Minute)
		assert.Len(t, info.LanguageStats[0].Platforms, 2)
		assert.Equal(t, "common", info.LanguageStats[0].Platforms[0].Name)
		assert.Equal(t, 1, info.LanguageStats[0].Platforms[0].Pages)
		assert.Equal(t, "linux", info.LanguageStats[0].Platforms[1].Name)
		assert.Equal(t, 2, info.LanguageStats[0].Platforms[1].Pages)
	})

	t.Run("error_on_non_existent_dir", func(t *testing.T) {
		c := &Cache{dir: "/nonexistent/path"}
		_, err := c.Info()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache directory")
	})

	t.Run("error_on_file_instead_of_dir", func(t *testing.T) {
		dir := t.TempDir()
		filePath := filepath.Join(dir, "not_a_dir")
		require.NoError(t, os.WriteFile(filePath, nil, 0o644))

		c := &Cache{dir: filePath}
		_, err := c.Info()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})

	t.Run("empty_cache_returns_zero_pages", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}
		_, err := c.Info()
		assert.Error(t, err)
	})

	t.Run("cache_with_multiple_languages", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.zh", "common"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.zh", "common", "git.md"), nil, 0o644))

		c := &Cache{dir: dir}
		info, err := c.Info()
		require.NoError(t, err)
		assert.Equal(t, 2, info.TotalPages)
		assert.Len(t, info.LanguageStats, 2)
		assert.Len(t, info.LanguageStats[0].Platforms, 1)
		assert.Equal(t, "common", info.LanguageStats[0].Platforms[0].Name)
		assert.Equal(t, 1, info.LanguageStats[0].Platforms[0].Pages)
	})
}

func TestLanguageStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		platforms []string
		langDirs  []string
		setup     func(t *testing.T, dir string)
		wantStats []LanguageInfo
		wantTotal int
	}{
		{
			name:      "single_language",
			platforms: []string{"common", "linux"},
			langDirs:  []string{"pages.en"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "linux"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "ls.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "linux", "apt.md"), nil, 0o644))
			},
			wantStats: []LanguageInfo{
				{
					Language: "en",
					Pages:    3,
					Platforms: []PlatformInfo{
						{Name: "common", Pages: 2},
						{Name: "linux", Pages: 1},
					},
				},
			},
			wantTotal: 3,
		},
		{
			name:      "multiple_languages",
			platforms: []string{"common"},
			langDirs:  []string{"pages.en", "pages.zh"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.zh", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.zh", "common", "git.md"), nil, 0o644))
			},
			wantStats: []LanguageInfo{
				{
					Language:  "en",
					Pages:     1,
					Platforms: []PlatformInfo{{Name: "common", Pages: 1}},
				},
				{
					Language:  "zh",
					Pages:     1,
					Platforms: []PlatformInfo{{Name: "common", Pages: 1}},
				},
			},
			wantTotal: 2,
		},
		{
			name:      "strips_pages_prefix",
			platforms: []string{"common"},
			langDirs:  []string{"pages.en", "pages.de"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.de", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.de", "common", "git.md"), nil, 0o644))
			},
			wantStats: []LanguageInfo{
				{
					Language:  "en",
					Pages:     1,
					Platforms: []PlatformInfo{{Name: "common", Pages: 1}},
				},
				{
					Language:  "de",
					Pages:     1,
					Platforms: []PlatformInfo{{Name: "common", Pages: 1}},
				},
			},
			wantTotal: 2,
		},
		{
			name:      "empty_directories_list",
			platforms: []string{"common"},
			langDirs:  nil,
			setup:     func(t *testing.T, dir string) {},
			wantStats: []LanguageInfo{},
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(t, dir)

			c := &Cache{dir: dir}
			stats, total, err := c.languageStats(tt.platforms, tt.langDirs)
			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.Equal(t, tt.wantStats, stats)
		})
	}
}

func TestPlatformStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		platforms []string
		setup     func(t *testing.T, dir string)
		want      []PlatformInfo
		wantTotal int
	}{
		{
			name:      "single_platform",
			platforms: []string{"common"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "ls.md"), nil, 0o644))
			},
			want:      []PlatformInfo{{Name: "common", Pages: 2}},
			wantTotal: 2,
		},
		{
			name:      "multiple_platforms",
			platforms: []string{"common", "linux"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "linux"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "linux", "apt.md"), nil, 0o644))
			},
			want: []PlatformInfo{
				{Name: "common", Pages: 1},
				{Name: "linux", Pages: 1},
			},
			wantTotal: 2,
		},
		{
			name:      "skips_missing_platform",
			platforms: []string{"common", "linux"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
			},
			want:      []PlatformInfo{{Name: "common", Pages: 1}},
			wantTotal: 1,
		},
		{
			name:      "empty_platforms_list",
			platforms: []string{},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
			},
			want:      nil,
			wantTotal: 0,
		},
		{
			name:      "ignores_non_md_files",
			platforms: []string{"common"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "notes.txt"), nil, 0o644))
			},
			want:      []PlatformInfo{{Name: "common", Pages: 1}},
			wantTotal: 1,
		},
		{
			name:      "empty_directory",
			platforms: []string{"common"},
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
			},
			want:      []PlatformInfo{{Name: "common", Pages: 0}},
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(t, dir)

			c := &Cache{dir: dir}
			got, total, err := c.platformStats("pages.en", tt.platforms)
			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.Equal(t, tt.want, got)
		})
	}
}
