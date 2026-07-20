package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/internal/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		results    *cache.FindResult
		query      string
		platform   string
		wantPath   string
		wantPlat   string
		wantErr    bool
		wantStderr string
	}{
		{
			name:     "exact_match",
			results:  &cache.FindResult{Matches: []string{"/pages.en/common/tar.md"}},
			query:    "tar",
			platform: "linux",
			wantPath: "/pages.en/common/tar.md",
			wantPlat: "linux",
		},
		{
			name:     "multiple_matches_uses_first",
			results:  &cache.FindResult{Matches: []string{"/pages.en/common/tar.md", "/pages.en/linux/tar.md"}},
			query:    "tar",
			platform: "linux",
			wantPath: "/pages.en/common/tar.md",
			wantPlat: "linux",
		},
		{
			name:       "fallback_only",
			results:    &cache.FindResult{Fallbacks: []string{"/pages.en/osx/tar.md"}},
			query:      "tar",
			platform:   "linux",
			wantPath:   "/pages.en/osx/tar.md",
			wantPlat:   "osx",
			wantStderr: "1. osx (tldr --platform osx tar)\n",
		},
		{
			name: "matches_preferred_over_fallbacks",
			results: &cache.FindResult{
				Matches:   []string{"/pages.en/common/tar.md"},
				Fallbacks: []string{"/pages.en/osx/tar.md"},
			},
			query:      "tar",
			platform:   "linux",
			wantPath:   "/pages.en/common/tar.md",
			wantPlat:   "linux",
			wantStderr: "1. osx (tldr --platform osx tar)\n",
		},
		{
			name:     "no_results",
			results:  &cache.FindResult{},
			query:    "tar",
			platform: "linux",
			wantErr:  true,
		},
		{
			name:     "nil_matches_and_fallbacks",
			results:  &cache.FindResult{},
			query:    "tar",
			platform: "linux",
			wantErr:  true,
		},
		{
			name:       "multiple_fallbacks",
			results:    &cache.FindResult{Fallbacks: []string{"/pages.en/osx/tar.md", "/pages.en/windows/tar.md"}},
			query:      "tar",
			platform:   "linux",
			wantPath:   "/pages.en/osx/tar.md",
			wantPlat:   "osx",
			wantStderr: "1. osx (tldr --platform osx tar)\n2. windows (tldr --platform windows tar)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			a := &App{Stderr: &stderr}

			gotPath, gotPlat, err := a.selectPage(tt.results, tt.query, tt.platform)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantPath, gotPath)
			assert.Equal(t, tt.wantPlat, gotPlat)

			if tt.wantStderr != "" {
				assert.Equal(t, tt.wantStderr, stderr.String())
			} else {
				assert.Empty(t, stderr.String())
			}
		})
	}
}

func TestRenderOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cli     cmd.CLI
		wantLen int
	}{
		{
			name:    "default",
			cli:     cmd.CLI{},
			wantLen: 2,
		},
		{
			name:    "color_always",
			cli:     cmd.CLI{Color: "always"},
			wantLen: 3,
		},
		{
			name:    "color_never",
			cli:     cmd.CLI{Color: "never"},
			wantLen: 3,
		},
		{
			name:    "edit",
			cli:     cmd.CLI{Edit: true},
			wantLen: 2,
		},
		{
			name:    "compact",
			cli:     cmd.CLI{Compact: true},
			wantLen: 2,
		},
		{
			name:    "no_compact",
			cli:     cmd.CLI{NoCompact: true},
			wantLen: 2,
		},
		{
			name:    "raw",
			cli:     cmd.CLI{Raw: true},
			wantLen: 2,
		},
		{
			name:    "no_raw",
			cli:     cmd.CLI{NoRaw: true},
			wantLen: 2,
		},
		{
			name:    "short_options",
			cli:     cmd.CLI{ShortOptions: true},
			wantLen: 2,
		},
		{
			name:    "long_options",
			cli:     cmd.CLI{LongOptions: true},
			wantLen: 2,
		},
		{
			name:    "both_options",
			cli:     cmd.CLI{ShortOptions: true, LongOptions: true},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{Stdout: &bytes.Buffer{}}
			got := a.renderOptions(&tt.cli)
			assert.Equal(t, tt.wantLen, len(got))
		})
	}
}

func TestLoadPage(t *testing.T) {
	t.Parallel()

	validPage := "# tar\n\n> archive utility.\n\n- create an archive:\n\n`tar cf archive.tar`\n"

	tests := []struct {
		name    string
		content string
		wantErr bool
		check   func(t *testing.T, page *render.Page, path string)
	}{
		{
			name:    "valid_page",
			content: validPage,
			wantErr: false,
			check: func(t *testing.T, page *render.Page, path string) {
				assert.Equal(t, "tar", page.Title)
				assert.Equal(t, path, page.Path)
				assert.Equal(t, validPage, page.RawContent)
				assert.NotEmpty(t, page.Examples)
			},
		},
		{
			name:    "valid_page_with_url",
			content: "# tar\n\n> archive utility.\n> More information: <https://example.org/tar>.\n\n- create:\n\n`tar cf archive.tar`\n",
			wantErr: false,
			check: func(t *testing.T, page *render.Page, _ string) {
				assert.Equal(t, "tar", page.Title)
				assert.Equal(t, "https://example.org/tar", page.URL)
			},
		},
		{
			name:    "nonexistent_file",
			content: "",
			wantErr: true,
		},
		{
			name:    "invalid_page_content",
			content: "some random text\n",
			wantErr: true,
		},
		{
			name:    "empty_file",
			content: "",
			wantErr: false,
			check: func(t *testing.T, page *render.Page, _ string) {
				assert.Empty(t, page.Title)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "page.md")

			if tt.name != "nonexistent_file" {
				require.NoError(t, writeTestFile(path, tt.content))
			}

			got, err := loadPage(path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			if tt.check != nil {
				tt.check(t, got, path)
			}
		})
	}
}

// writeTestFile creates or overwrites a test file
// at path with the given content.
func writeTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
