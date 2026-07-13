package render

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
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
