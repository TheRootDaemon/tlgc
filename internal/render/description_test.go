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

func TestRenderDescriptionLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		indent  string
		writer  io.Writer
		want    string
		wantErr string
	}{
		{
			name:   "writes text with trailing newline",
			text:   "hello",
			indent: "  ",
			want:   "  hello\n",
		},
		{
			name:    "write error",
			text:    "hello",
			indent:  "  ",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			w := io.Writer(&buf)
			r := &Renderer{
				w: w,
			}
			if tt.writer != nil {
				r.w = tt.writer
			}

			err := r.renderDescriptionLine(tt.text, tt.indent)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestRenderDescriptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		descs   []string
		url     string
		compact bool
		writer  io.Writer
		want    string
		wantErr string
	}{
		{
			name:  "no descriptions and no url returns nil",
			descs: nil,
			url:   "",
			want:  "",
		},
		{
			name:  "single description",
			descs: []string{"hello"},
			want:  "  hello\n\n",
		},
		{
			name:    "single description compact omits trailing newline",
			descs:   []string{"hello"},
			compact: true,
			want:    "  hello\n",
		},
		{
			name:  "multiple descriptions",
			descs: []string{"first", "second"},
			want:  "  first\n  second\n\n",
		},
		{
			name:  "description with URL",
			descs: []string{"hello"},
			url:   "https://example.org",
			want:  "  hello\n  More information: https://example.org.\n\n",
		},
		{
			name: "URL only no descriptions",
			url:  "https://example.org",
			want: "  More information: https://example.org.\n\n",
		},
		{
			name:    "URL only compact",
			url:     "https://example.org",
			compact: true,
			want:    "  More information: https://example.org.\n",
		},
		{
			name:    "write error",
			descs:   []string{"hello"},
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			w := io.Writer(&buf)
			r := &Renderer{
				w:      w,
				indent: config.IndentConfig{Description: 2},
				output: config.OutputConfig{Compact: tt.compact},
			}

			if tt.writer != nil {
				r.w = tt.writer
			}

			err := r.renderDescriptions(tt.descs, tt.url)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestRenderDescriptionURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		url        string
		indent     string
		lineLength int
		useColor   bool
		descStyle  config.OutputStyle
		urlStyle   config.OutputStyle
		writer     io.Writer
		want       string
		wantANSI   bool
		wantErr    string
	}{
		{
			name:   "writes more information line",
			url:    "https://example.org",
			indent: "  ",
			want:   "  More information: https://example.org.\n",
		},
		{
			name:      "url includes ANSI styles when color enabled",
			url:       "https://example.org",
			indent:    "  ",
			useColor:  true,
			descStyle: config.OutputStyle{Bold: true},
			urlStyle:  config.OutputStyle{Italic: true},
			wantANSI:  true,
		},
		{
			name:    "write error",
			url:     "https://example.org",
			indent:  "  ",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			w := io.Writer(&buf)
			r := &Renderer{
				w:        w,
				useColor: tt.useColor,
				style: config.StyleConfig{
					Description: tt.descStyle,
					URL:         tt.urlStyle,
				},
				output: config.OutputConfig{LineLength: tt.lineLength},
			}

			if tt.writer != nil {
				r.w = tt.writer
			}

			err := r.renderDescriptionURL(tt.url, tt.indent)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			got := buf.String()

			if tt.wantANSI {
				assert.Contains(t, got, "\x1b[")
				assert.Contains(t, got, "More information:")
				assert.Contains(t, got, "https://example.org")
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRenderBulletLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		text    string
		indent  string
		writer  io.Writer
		want    string
		wantErr string
	}{
		{
			name:   "writes text without trailing newline",
			text:   "hello",
			indent: "  ",
			want:   "  hello",
		},
		{
			name:    "write error",
			text:    "hello",
			indent:  "  ",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			w := io.Writer(&buf)
			r := &Renderer{
				w: w,
			}

			if tt.writer != nil {
				r.w = tt.writer
			}

			err := r.renderBulletLine(tt.text, tt.indent)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
