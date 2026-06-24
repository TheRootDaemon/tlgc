package cache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupDir  func(t *testing.T) string
		query     string
		platform  string
		languages []string
		want      *FindResult
		wantErr   bool
	}{
		{
			name: "match_in_requested_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				return d
			},
			query:     "apt",
			platform:  "linux",
			languages: []string{"en"},
			want: &FindResult{
				Matches:   []string{filepath.Join("testroot", "pages.en", "linux", "apt.md")},
				Fallbacks: nil,
			},
		},
		{
			name: "fallback_to_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:     "ls",
			platform:  "linux",
			languages: []string{"en"},
			want: &FindResult{
				Matches:   []string{filepath.Join("testroot", "pages.en", "common", "ls.md")},
				Fallbacks: nil,
			},
		},
		{
			name: "both_primary_and_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:     "apt",
			platform:  "linux",
			languages: []string{"en"},
			want: &FindResult{
				Matches:   []string{filepath.Join("testroot", "pages.en", "linux", "apt.md")},
				Fallbacks: nil,
			},
		},
		{
			name: "fallback_to_other_platforms",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "osx"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "osx", "brew.md"), nil, 0o644))
				return d
			},
			query:     "brew",
			platform:  "linux",
			languages: []string{"en"},
			want: &FindResult{
				Matches:   nil,
				Fallbacks: []string{filepath.Join("testroot", "pages.en", "osx", "brew.md")},
			},
		},
		{
			name: "no_match_returns_empty_lists",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			query:     "nonexistent",
			platform:  "linux",
			languages: []string{"en"},
			want: &FindResult{
				Matches:   nil,
				Fallbacks: nil,
			},
		},
		{
			name: "error_on_unknown_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			query:     "apt",
			platform:  "nonexistent",
			languages: []string{"en"},
			wantErr:   true,
		},
		{
			name: "common_platform_only",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			query:     "ls",
			platform:  "common",
			languages: []string{"en"},
			want: &FindResult{
				Matches:   []string{filepath.Join("testroot", "pages.en", "common", "ls.md")},
				Fallbacks: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got, err := c.Find(tt.query, tt.platform, tt.languages)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Fix up expected paths to use the actual dir.
			if tt.want != nil {
				for i := range tt.want.Matches {
					tt.want.Matches[i] = strings.Replace(tt.want.Matches[i], "testroot", dir, 1)
				}
				for i := range tt.want.Fallbacks {
					tt.want.Fallbacks[i] = strings.Replace(tt.want.Fallbacks[i], "testroot", dir, 1)
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind_PrimaryMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		platform string
		file     string
		langDirs []string
		want     []string
	}{
		{
			name: "found_in_requested_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				return d
			},
			platform: "linux",
			file:     "apt.md",
			langDirs: []string{"pages.en"},
			want:     []string{filepath.Join("testroot", "pages.en", "linux", "apt.md")},
		},
		{
			name: "fallback_to_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			platform: "linux",
			file:     "ls.md",
			langDirs: []string{"pages.en"},
			want:     []string{filepath.Join("testroot", "pages.en", "common", "ls.md")},
		},
		{
			name: "both_requested_and_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "apt.md"), nil, 0o644))
				return d
			},
			platform: "linux",
			file:     "apt.md",
			langDirs: []string{"pages.en"},
			want: []string{
				filepath.Join("testroot", "pages.en", "linux", "apt.md"),
				filepath.Join("testroot", "pages.en", "common", "apt.md"),
			},
		},
		{
			name: "common_platform_skips_self",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			platform: "common",
			file:     "ls.md",
			langDirs: []string{"pages.en"},
			want:     []string{filepath.Join("testroot", "pages.en", "common", "ls.md")},
		},
		{
			name: "not_found",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			platform: "linux",
			file:     "nonexistent.md",
			langDirs: []string{"pages.en"},
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			for i := range tt.want {
				tt.want[i] = strings.Replace(tt.want[i], "testroot", dir, 1)
			}
			got := c.primaryMatches(tt.platform, tt.file, tt.langDirs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind_FallbackMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupDir  func(t *testing.T) string
		file      string
		platform  string
		platforms []string
		langDirs  []string
		want      []string
	}{
		{
			name: "found_in_other_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "osx"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "osx", "brew.md"), nil, 0o644))
				return d
			},
			file:      "brew.md",
			platform:  "linux",
			platforms: []string{"common", "linux", "osx"},
			langDirs:  []string{"pages.en"},
			want:      []string{filepath.Join("testroot", "pages.en", "osx", "brew.md")},
		},
		{
			name: "skips_requested_platform",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "linux"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "linux", "apt.md"), nil, 0o644))
				return d
			},
			file:      "apt.md",
			platform:  "linux",
			platforms: []string{"common", "linux"},
			langDirs:  []string{"pages.en"},
			want:      nil,
		},
		{
			name: "skips_common",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "ls.md"), nil, 0o644))
				return d
			},
			file:      "ls.md",
			platform:  "linux",
			platforms: []string{"common", "linux"},
			langDirs:  []string{"pages.en"},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			for i := range tt.want {
				tt.want[i] = strings.Replace(tt.want[i], "testroot", dir, 1)
			}
			got := c.fallbackMatches(tt.file, tt.platform, tt.platforms, tt.langDirs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind_PageFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		file     string
		platform string
		langDirs []string
		want     string
	}{
		{
			name: "found_in_first_lang_dir",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				return d
			},
			file:     "git.md",
			platform: "common",
			langDirs: []string{"pages.en"},
			want:     filepath.Join("testroot", "pages.en", "common", "git.md"),
		},
		{
			name: "searches_lang_dirs_in_order",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.de", "common"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.en", "common", "git.md"), nil, 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(d, "pages.de", "common", "git.md"), nil, 0o644))
				return d
			},
			file:     "git.md",
			platform: "common",
			langDirs: []string{"pages.en", "pages.de"},
			want:     filepath.Join("testroot", "pages.en", "common", "git.md"),
		},
		{
			name: "not_found_returns_empty",
			setupDir: func(t *testing.T) string {
				d := t.TempDir()
				require.NoError(t, os.MkdirAll(filepath.Join(d, "pages.en", "common"), 0o755))
				return d
			},
			file:     "nonexistent.md",
			platform: "common",
			langDirs: []string{"pages.en"},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			tt.want = strings.Replace(tt.want, "testroot", dir, 1)
			got := c.findPageFor(tt.file, tt.platform, tt.langDirs)
			assert.Equal(t, tt.want, got)
		})
	}
}
