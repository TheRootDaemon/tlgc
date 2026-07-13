package render

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
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
			r := &Renderer{}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderDescriptionLine(w, tt.text, tt.indent)

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
			name:    "write error",
			descs:   []string{"hello"},
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				indent: config.IndentConfig{Description: 2},
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderDescriptions(w, tt.descs, tt.url)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
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
			r := &Renderer{}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderBulletLine(w, tt.text, tt.indent)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
