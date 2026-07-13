package render

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRenderExample(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ex      Example
		indent  config.IndentConfig
		output  config.OutputConfig
		writer  io.Writer
		want    string
		wantErr string
	}{
		{
			name:   "description and command non-compact",
			ex:     Example{Description: "create archive", Command: "tar cf archive.tar"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			want:   "  create archive\n\n    tar cf archive.tar\n",
		},
		{
			name:   "description only no command",
			ex:     Example{Description: "just a description"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			want:   "  just a description\n\n",
		},
		{
			name:   "hyphens enabled",
			ex:     Example{Description: "create archive", Command: "tar cf archive.tar"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			output: config.OutputConfig{ShowHyphens: true, ExamplePrefix: "- "},
			want:   "  - create archive\n\n    tar cf archive.tar\n",
		},
		{
			name:   "compact mode no blank line",
			ex:     Example{Description: "create archive", Command: "tar cf archive.tar"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			output: config.OutputConfig{Compact: true},
			want:   "  create archive\n    tar cf archive.tar\n",
		},
		{
			name:    "write error",
			ex:      Example{Description: "error"},
			indent:  config.IndentConfig{Bullet: 2, Example: 4},
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				useColor: false,
				output:   tt.output,
				indent:   tt.indent,
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderExample(w, tt.ex)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
