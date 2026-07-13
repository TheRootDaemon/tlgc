package render

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestWrapText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		lineLength int
		s          string
		indent     string
		want       string
	}{
		{
			name:       "line length zero returns indent plus text",
			lineLength: 0,
			s:          "hello world",
			indent:     "  ",
			want:       "  hello world",
		},
		{
			name:       "line length negative returns indent plus text",
			lineLength: -1,
			s:          "short",
			indent:     "  ",
			want:       "  short",
		},
		{
			name:       "empty text returns indent",
			lineLength: 80,
			s:          "",
			indent:     ">>",
			want:       ">>",
		},
		{
			name:       "text fits within line length",
			lineLength: 80,
			s:          "hi",
			indent:     "  ",
			want:       "  hi",
		},
		{
			name:       "text wraps with indent on continuation lines",
			lineLength: 12,
			s:          "hello world foo",
			indent:     "> ",
			want:       "> hello world\n> foo",
		},
		{
			name:       "long word exceeds line length without splitting",
			lineLength: 5,
			s:          "superlongword",
			indent:     "- ",
			want:       "- superlongword",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				output: config.OutputConfig{LineLength: tt.lineLength},
			}
			got := r.wrapText(tt.s, tt.indent)
			assert.Equal(t, tt.want, got)
		})
	}
}
