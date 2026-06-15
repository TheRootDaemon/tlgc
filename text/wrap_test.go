package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		indent string
		want   string
	}{
		{name: "short line", input: "hello world", maxLen: 80, indent: "  ", want: "hello world"},
		{name: "long line", input: "this is a long line that should be wrapped at word boundaries", maxLen: 20, indent: "", want: "this is a long line\nthat should be\nwrapped at word\nboundaries"},
		{name: "with indent", input: "this is a long line that should be indented", maxLen: 15, indent: "> ", want: "this is a long\n> line that\n> should be\n> indented"},
		{name: "empty", input: "", maxLen: 80, indent: "", want: ""},
		{name: "zero maxLen", input: "hello world", maxLen: 0, indent: "", want: "hello\nworld"},
		{name: "negative maxLen", input: "hello world", maxLen: -1, indent: "", want: "hello\nworld"},
		{name: "word longer than maxLen", input: "superlongword that is long", maxLen: 5, indent: "", want: "superlongword\nthat\nis\nlong"},
		{name: "exact fit", input: "1234567890", maxLen: 10, indent: "", want: "1234567890"},
		{name: "exact fit with space", input: "hello world", maxLen: 11, indent: "", want: "hello world"},
		{name: "single word", input: "hello", maxLen: 10, indent: "", want: "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.input, tt.maxLen, tt.indent)
			assert.Equal(t, tt.want, got)
		})
	}
}
