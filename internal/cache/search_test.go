package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSearch tests Cache.Search.
func TestSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupDir  func(t *testing.T) string
		query     string
		platform  string
		languages []string
		want      []SearchResult
		wantErr   bool
	}{
		{
			name: "case_insensitive_substring_match",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git-add.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git-commit.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:     "GIT",
			platform:  "common",
			languages: []string{"en"},
			want: []SearchResult{
				{Page: "git-add", Language: "en", Platform: "common"},
				{Page: "git-commit", Language: "en", Platform: "common"},
			},
		},
		{
			name: "results_sorted_by_page_name",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "java.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "apt.md"), nil, 0o644))
				return d
			},
			query:     "a",
			platform:  "common",
			languages: []string{"en"},
			want: []SearchResult{
				{Page: "apt", Language: "en", Platform: "common"},
				{Page: "java", Language: "en", Platform: "common"},
			},
		},
		{
			name: "empty_platform_searches_all",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "osx"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "osx", "bat.md"), nil, 0o644))
				return d
			},
			query:     "a",
			platform:  "",
			languages: []string{"en"},
			want: []SearchResult{
				{Page: "apt", Language: "en", Platform: "linux"},
				{Page: "bat", Language: "en", Platform: "osx"},
			},
		},
		{
			name: "no_match_returns_error",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:     "zzzznotfound",
			platform:  "common",
			languages: []string{"en"},
			wantErr:   true,
		},
		{
			name: "error_on_unknown_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			query:     "ls",
			platform:  "nonexistent",
			languages: []string{"en"},
			wantErr:   true,
		},
		{
			name: "error_on_no_matching_languages",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			query:     "ls",
			platform:  "common",
			languages: []string{"de"},
			wantErr:   true,
		},
		{
			name: "specific_platform_only",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "osx"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "osx", "brew.md"), nil, 0o644))
				return d
			},
			query:     "a",
			platform:  "linux",
			languages: []string{"en"},
			want: []SearchResult{
				{Page: "apt", Language: "en", Platform: "linux"},
			},
		},
		{
			name: "searches_across_languages",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.de", "common", "git.md"), nil, 0o644))
				return d
			},
			query:     "git",
			platform:  "common",
			languages: []string{"en", "de"},
			want: []SearchResult{
				{Page: "git", Language: "en", Platform: "common"},
				{Page: "git", Language: "de", Platform: "common"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.Search(tt.query, tt.platform, tt.languages)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestResolvePlatforms tests Cache.resolvePlatforms.
func TestResolvePlatforms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		platform string
		want     []string
		wantErr  bool
	}{
		{
			name: "empty_returns_all",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform: "",
			want:     []string{"common", "linux"},
		},
		{
			name: "common_returns_only_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform: "common",
			want:     []string{"common"},
		},
		{
			name: "specific_platform_and_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform: "linux",
			want:     []string{"linux", "common"},
		},
		{
			name: "error_on_unknown_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.resolvePlatforms(tt.platform)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestPlatformExists tests platformExists.
func TestPlatformExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		platforms []string
		platform  string
		want      bool
	}{
		{
			name:      "platform_exists",
			platforms: []string{"linux", "osx", "common"},
			platform:  "linux",
			want:      true,
		},
		{
			name:      "platform_does_not_exist",
			platforms: []string{"linux", "osx"},
			platform:  "windows",
			want:      false,
		},
		{
			name:      "empty_platforms_list",
			platforms: []string{},
			platform:  "linux",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := platformExists(tt.platforms, tt.platform)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestSearchDirectory tests Cache.searchDirectory.
func TestSearchDirectory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		setupDir          func(t *testing.T) string
		query             string
		platform          string
		languageDirectory string
		want              []SearchResult
	}{
		{
			name: "matches_substring_case_insensitive",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git-add.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git-commit.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:             "commit",
			platform:          "common",
			languageDirectory: "pages.en",
			want: []SearchResult{
				{Page: "git-commit", Language: "en", Platform: "common"},
			},
		},
		{
			name: "no_match_returns_empty",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:             "zzzz",
			platform:          "common",
			languageDirectory: "pages.en",
			want:              nil,
		},
		{
			name: "nonexistent_directory_returns_empty",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			query:             "ls",
			platform:          "nonexistent",
			languageDirectory: "pages.en",
			want:              nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got := c.searchDirectory(tt.query, tt.platform, tt.languageDirectory)
			assert.Equal(t, tt.want, got)
		})
	}
}
