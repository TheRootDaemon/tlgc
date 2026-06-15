package termcolor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSprint(t *testing.T) {
	tests := []struct {
		name       string
		style      string
		text       string
		expectSame bool
	}{
		{
			name:       "empty string returns empty string",
			style:      "",
			text:       "",
			expectSame: true,
		},
		{
			name:       "empty style returns same text",
			style:      "",
			text:       "hello",
			expectSame: true,
		},
		{
			name:       "default style returns same text",
			style:      "default",
			text:       "hello",
			expectSame: true,
		},
		{
			name:       "valid style wraps text",
			style:      "red",
			text:       "hello",
			expectSame: false,
		},
		{
			name:       "multiple effects",
			style:      "bold underline green",
			text:       "hi",
			expectSame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sprint(tt.style, tt.text)

			if tt.expectSame {
				require.Equal(t, tt.text, got)
				return
			}

			require.Contains(t, got, tt.text)
			require.Contains(t, got, Reset())
			require.NotEqual(t, tt.text, got)
		})
	}
}

func TestFprintf(t *testing.T) {
	tests := []struct {
		name   string
		style  string
		format string
		args   []any
	}{
		{
			name:   "basic formatting",
			style:  "red",
			format: "hello %s",
			args:   []any{"world"},
		},
		{
			name:   "no style",
			style:  "",
			format: "%d + %d = %d",
			args:   []any{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fprintf(tt.style, tt.format, tt.args...)

			expectedContent := fmt.Sprintf(tt.format, tt.args...)

			require.Contains(t, got, expectedContent)

			if tt.style == "" {
				require.Equal(t, expectedContent, got)
			} else {
				require.NotEqual(t, expectedContent, got)
				require.Contains(t, got, Reset())
			}
		})
	}
}
