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

func TestRenderTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		title    string
		indent   int
		useColor bool
		style    config.OutputStyle
		writer   io.Writer
		want     string
		wantErr  string
	}{
		{
			name:   "title with indent",
			title:  "tar",
			indent: 2,
			want:   "  tar",
		},
		{
			name:   "zero indent",
			title:  "tar",
			indent: 0,
			want:   "tar",
		},
		{
			name:     "colorized title",
			title:    "tar",
			indent:   2,
			useColor: true,
			style: config.OutputStyle{
				Bold:  true,
				Color: config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorRed},
			},
		},
		{
			name:    "write error",
			title:   "tar",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				useColor: tt.useColor,
				style:    config.StyleConfig{Title: tt.style},
				indent:   config.IndentConfig{Title: tt.indent},
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderTitle(w, tt.title)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			got := buf.String()

			if tt.useColor {
				assert.Contains(t, got, "  tar")
				assert.Contains(t, got, "\x1b[")
				assert.True(t, strings.HasSuffix(got, "\x1b[0m"))
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRenderPageTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		platform      string
		page          *Page
		showTitle     bool
		platformTitle bool
		compact       bool
		writer        io.Writer
		want          string
		wantErr       string
	}{
		{
			name:      "show title disabled",
			showTitle: false,
			page:      &Page{Title: "tar"},
			want:      "",
		},
		{
			name:      "empty title",
			showTitle: true,
			page:      &Page{Title: ""},
			want:      "",
		},
		{
			name:      "normal title non-compact",
			showTitle: true,
			page:      &Page{Title: "tar"},
			want:      "  tar\n\n",
		},
		{
			name:      "normal title compact",
			showTitle: true,
			compact:   true,
			page:      &Page{Title: "tar"},
			want:      "  tar\n",
		},
		{
			name:          "platform prefix",
			showTitle:     true,
			platformTitle: true,
			platform:      "linux",
			page:          &Page{Title: "tar"},
			want:          "  linux/tar\n\n",
		},
		{
			name:          "platform title enabled no platform",
			showTitle:     true,
			platformTitle: true,
			platform:      "",
			page:          &Page{Title: "tar"},
			want:          "  tar\n\n",
		},
		{
			name:      "write error",
			showTitle: true,
			page:      &Page{Title: "tar"},
			writer:    &errorWriter{err: errors.New("write error")},
			wantErr:   "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			r := &Renderer{
				w: &buf,
				output: config.OutputConfig{
					ShowTitle:     tt.showTitle,
					PlatformTitle: tt.platformTitle,
					Compact:       tt.compact,
				},
				indent: config.IndentConfig{Title: 2},
			}

			if tt.writer != nil {
				r.w = tt.writer
			}

			err := r.renderPageTitle(tt.platform, tt.page)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
