package render

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rawContent string
		path       string
		content    string
		writer     io.Writer
		want       string
		wantErr    string
		wantAnyErr bool
	}{
		{
			name:       "uses RawContent when set",
			rawContent: "# hello from content",
			want:       "# hello from content",
		},
		{
			name:       "RawContent takes precedence over Path",
			rawContent: "from content",
			path:       "/some/nonexistent/path.md",
			want:       "from content",
		},
		{
			name:       "falls back to file when RawContent is empty",
			rawContent: "",
			content:    "# file content",
			want:       "# file content",
		},
		{
			name:       "file not found when RawContent empty and path invalid",
			rawContent: "",
			path:       "/nonexistent/file.md",
			wantAnyErr: true,
		},
		{
			name:       "write error",
			rawContent: "data",
			writer:     &errorWriter{err: errors.New("write error")},
			wantErr:    "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if tt.content != "" {
				path = filepath.Join(t.TempDir(), "page.md")
				require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o644))
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			r := &Renderer{w: w}
			err := r.renderRaw(&Page{RawContent: tt.rawContent, Path: path})

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			if tt.wantAnyErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
