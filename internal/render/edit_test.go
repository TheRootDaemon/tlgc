package render

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderEditLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		compact bool
		url     string
		writer  io.Writer
		want    string
		wantErr string
	}{
		{
			name:    "non-compact adds newline",
			compact: false,
			url:     "https://example.com",
			want:    "https://example.com\n",
		},
		{
			name:    "compact omits newline",
			compact: true,
			url:     "https://example.com",
			want:    "https://example.com",
		},
		{
			name:    "write error",
			compact: false,
			url:     "url",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				output: config.OutputConfig{Compact: tt.compact},
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderEditLink(w, tt.url)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestRenderPageEditLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		page     *Page
		editLink bool
		compact  bool
		writer   io.Writer
		want     string
		wantErr  string
	}{
		{
			name:     "edit link disabled",
			editLink: false,
			page:     &Page{},
			want:     "",
		},
		{
			name:     "empty path and url",
			editLink: true,
			page:     &Page{},
			want:     "",
		},
		{
			name:     "renders edit link from path",
			editLink: true,
			page:     &Page{Path: "/pages/common/tar.md"},
			want:     "https://github.com/tldr-pages/tldr/edit/main/pages/common/tar.md\n",
		},
		{
			name:     "custom url returned as-is",
			editLink: true,
			page:     &Page{URL: "https://custom.com/edit"},
			want:     "https://custom.com/edit\n",
		},
		{
			name:     "compact mode omits newline",
			editLink: true,
			compact:  true,
			page:     &Page{Path: "/pages/common/tar.md"},
			want:     "https://github.com/tldr-pages/tldr/edit/main/pages/common/tar.md",
		},
		{
			name:     "write error",
			editLink: true,
			page:     &Page{URL: "https://example.com"},
			writer:   &errorWriter{err: errors.New("write error")},
			wantErr:  "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			r := &Renderer{
				w:      &buf,
				output: config.OutputConfig{EditLink: tt.editLink, Compact: tt.compact},
			}

			if tt.writer != nil {
				r.w = tt.writer
			}

			err := r.renderPageEditLink(tt.page)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestBuildEditURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		url  string
		want string
	}{
		{
			name: "non-empty url returned as-is",
			path: "/pages/common/tar.md",
			url:  "https://example.com",
			want: "https://example.com",
		},
		{
			name: "empty url constructs from path",
			path: "/pages/common/tar.md",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/common/tar.md",
		},
		{
			name: "linux platform extracted correctly",
			path: "/pages/linux/apt.md",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/linux/apt.md",
		},
		{
			name: "windows platform extracted correctly",
			path: "/pages/windows/dir.md",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/windows/dir.md",
		},
		{
			name: "empty path and empty url returns empty",
			path: "",
			url:  "",
			want: "",
		},
		{
			name: "path without .md extension adds .md",
			path: "/pages/common/some-page",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/common/some-page.md",
		},
		{
			name: "url takes precedence over path",
			path: "/pages/common/tar.md",
			url:  "https://custom.com/edit",
			want: "https://custom.com/edit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildEditURL(tt.path, tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}
