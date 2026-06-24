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
	})
}

func TestLanguageStats(t *testing.T) {
	t.Parallel()

	t.Run("counts_pages_across_platforms", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "linux"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "ls.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "linux", "apt.md"), nil, 0o644))

		c := &Cache{dir: dir}
		stats, total, err := c.languageStats(
			[]string{"common", "linux"},
			[]string{"pages.en"},
		)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, stats, 1)
		assert.Equal(t, "en", stats[0].Language)
		assert.Equal(t, 3, stats[0].Pages)
	})

	t.Run("multiple_languages", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.zh", "common"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.zh", "common", "git.md"), nil, 0o644))

		c := &Cache{dir: dir}
		stats, total, err := c.languageStats(
			[]string{"common"},
			[]string{"pages.en", "pages.zh"},
		)
		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, stats, 2)
	})

	t.Run("skips_non_existent_platform_dirs", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))

		c := &Cache{dir: dir}
		stats, total, err := c.languageStats(
			[]string{"common", "linux"},
			[]string{"pages.en"},
		)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, stats, 1)
	})

	t.Run("empty_directories_list", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}
		stats, total, err := c.languageStats(
			[]string{"common"},
			nil,
		)
		require.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Empty(t, stats)
	})

	t.Run("ignores_non_md_files", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "pages.en", "common"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "git.md"), nil, 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "pages.en", "common", "notes.txt"), nil, 0o644))

		c := &Cache{dir: dir}
		_, total, err := c.languageStats(
			[]string{"common"},
			[]string{"pages.en"},
		)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
	})
}
