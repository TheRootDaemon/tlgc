package render

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
)

// failOnNewline is a writer that succeeds on all writes except standalone "\n".
// Used to test error propagation on the newline write between wrapped lines.
type failOnNewline struct {
	buf strings.Builder
}

// Write writes p to the underlying buffer
// unless p is exactly "\n",
// in which case it returns a sentinel write error.
// This allows tests to simulate a failure during newline writes.
func (w *failOnNewline) Write(p []byte) (int, error) {
	if string(p) == "\n" {
		return 0, errors.New("write error")
	}
	return w.buf.Write(p)
}

func TestRenderStyledInline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		text       string
		indent     string
		lineLength int
		useColor   bool
		baseStyle  config.OutputStyle
		codeStyle  config.OutputStyle
		writer     io.Writer
		want       string
		wantErr    string
		wantANSI   bool
	}{
		{
			name:   "fast path single segment no backticks",
			text:   "hello world",
			indent: "  ",
			want:   "  hello world",
		},
		{
			name:   "multiple segments no wrap",
			text:   "hello `code` world",
			indent: "  ",
			want:   "  hello code world",
		},
		{
			name:       "wrapping across lines",
			text:       "some `code` here and there",
			indent:     "  ",
			lineLength: 15,
			want:       "  some code here\n  and there",
		},
		{
			name:      "color output contains ANSI escapes",
			text:      "hello `code` world",
			indent:    "  ",
			useColor:  true,
			baseStyle: config.OutputStyle{Bold: true},
			codeStyle: config.OutputStyle{Italic: true},
			wantANSI:  true,
		},
		{
			name:   "no color output is plain text",
			text:   "hello `code` world",
			indent: "  ",
			want:   "  hello code world",
		},
		{
			name:       "write error on newline between lines",
			text:       "aaa bbbb `cccc` dddd eeee",
			indent:     "  ",
			lineLength: 15,
			writer:     &failOnNewline{},
			wantErr:    "write error",
		},
		{
			name:    "write error on first content write",
			text:    "hello world",
			indent:  "  ",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				useColor: tt.useColor,
				output:   config.OutputConfig{LineLength: tt.lineLength},
				style: config.StyleConfig{
					Description: tt.baseStyle,
					InlineCode:  tt.codeStyle,
				},
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderStyledInline(
				w,
				tt.text,
				tt.indent,
				r.style.Description,
				r.style.InlineCode,
			)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			got := buf.String()

			if tt.wantANSI {
				assert.Contains(t, got, "\x1b[")
				assert.Contains(t, got, "hello")
				assert.Contains(t, got, "code")
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWrapInlineSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		segments   []inlineSegment
		indent     string
		lineLength int
		want       []string
	}{
		{
			name:     "single segment one line",
			segments: []inlineSegment{{text: "hello"}},
			indent:   "  ",
			want:     []string{"  hello"},
		},
		{
			name:     "multiple segments no wrap",
			segments: []inlineSegment{{text: "hello "}, {text: "code", code: true}},
			indent:   "  ",
			want:     []string{"  hello code"},
		},
		{
			name:       "wrapping at boundary produces two lines",
			segments:   []inlineSegment{{text: "aaa "}, {text: "bbb", code: true}, {text: " ccc"}},
			indent:     "  ",
			lineLength: 10,
			want:       []string{"  aaa bbb", "  ccc"},
		},
		{
			name:       "wrapping produces multiple lines",
			segments:   []inlineSegment{{text: "aaa bbb ccc ddd eee"}},
			indent:     "  ",
			lineLength: 10,
			want:       []string{"  aaa bbb", "  ccc ddd", "  eee"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				output: config.OutputConfig{LineLength: tt.lineLength},
			}
			got := r.wrapInlineSegments(tt.segments, tt.indent)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderStyledLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		line     string
		segments []inlineSegment
		writer   io.Writer
		want     string
		wantErr  string
	}{
		{
			name:     "single text segment",
			line:     "hello",
			segments: []inlineSegment{{text: "hello"}},
			want:     "hello",
		},
		{
			name:     "single code segment",
			line:     "hello",
			segments: []inlineSegment{{text: "hello", code: true}},
			want:     "hello",
		},
		{
			name: "multiple segments mixed",
			line: "hello world",
			segments: []inlineSegment{
				{text: "hello "},
				{text: "world", code: true},
			},
			want: "hello world",
		},
		{
			name:     "segment text not found writes remaining",
			line:     "hello",
			segments: []inlineSegment{{text: "xyz"}},
			want:     "hello",
		},
		{
			name:     "empty line returns nil",
			line:     "",
			segments: []inlineSegment{{text: "hello"}},
			want:     "",
		},
		{
			name: "remaining text after segments consumed",
			line: "hello world extra",
			segments: []inlineSegment{
				{text: "hello "},
				{text: "world", code: true},
			},
			want: "hello world extra",
		},
		{
			name: "break writes remaining when segment not found mid-line",
			line: "hello world",
			segments: []inlineSegment{
				{text: "hello"},
				{text: "xyz"},
			},
			want: "hello world",
		},
		{
			name:     "write error on before text",
			line:     "  hello",
			segments: []inlineSegment{{text: "hello"}},
			writer:   &errorWriter{err: errors.New("write error")},
			wantErr:  "write error",
		},
		{
			name:     "write error on code segment",
			line:     "hello",
			segments: []inlineSegment{{text: "hello", code: true}},
			writer:   &errorWriter{err: errors.New("write error")},
			wantErr:  "write error",
		},
		{
			name:     "write error on remaining text",
			line:     "hello",
			segments: []inlineSegment{{text: "xyz"}},
			writer:   &errorWriter{err: errors.New("write error")},
			wantErr:  "write error",
		},
		{
			name: "empty segment text skipped",
			line: "hello",
			segments: []inlineSegment{
				{text: ""},
				{text: "hello"},
			},
			want: "hello",
		},
		{
			name: "repeated segment text positional matching",
			line: "a a",
			segments: []inlineSegment{
				{text: "a"},
				{text: "a"},
			},
			want: "a a",
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

			err := r.renderStyledLine(
				w,
				tt.line,
				tt.segments,
				config.OutputStyle{},
				config.OutputStyle{},
			)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestParseInline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		want []inlineSegment
	}{
		{
			name: "no backticks",
			s:    "hello world",
			want: []inlineSegment{{text: "hello world"}},
		},
		{
			name: "empty string",
			s:    "",
			want: []inlineSegment{{text: ""}},
		},
		{
			name: "single code span",
			s:    "`code`",
			want: []inlineSegment{{text: "code", code: true}},
		},
		{
			name: "plain text with code span",
			s:    "hello `code` world",
			want: []inlineSegment{
				{text: "hello "},
				{text: "code", code: true},
				{text: " world"},
			},
		},
		{
			name: "multiple code spans",
			s:    "`a` and `b`",
			want: []inlineSegment{
				{text: "a", code: true},
				{text: " and "},
				{text: "b", code: true},
			},
		},
		{
			name: "leading plain text before code",
			s:    "text `code`",
			want: []inlineSegment{
				{text: "text "},
				{text: "code", code: true},
			},
		},
		{
			name: "trailing plain text after code",
			s:    "`code` text",
			want: []inlineSegment{
				{text: "code", code: true},
				{text: " text"},
			},
		},
		{
			name: "adjacent code spans",
			s:    "`a``b`",
			want: []inlineSegment{
				{text: "a", code: true},
				{text: "b", code: true},
			},
		},
		{
			name: "only double backticks",
			s:    "``",
			want: []inlineSegment{{text: "", code: true}},
		},
		{
			name: "only triple backticks",
			s:    "```",
			want: []inlineSegment{
				{text: "", code: true},
				{text: "", code: true},
			},
		},
		{
			name: "three consecutive code spans",
			s:    "`x``y``z`",
			want: []inlineSegment{
				{text: "x", code: true},
				{text: "y", code: true},
				{text: "z", code: true},
			},
		},
		{
			name: "code plain code",
			s:    "`a` x `b`",
			want: []inlineSegment{
				{text: "a", code: true},
				{text: " x "},
				{text: "b", code: true},
			},
		},
		{
			name: "plain text then backtick",
			s:    "text`",
			want: []inlineSegment{
				{text: "text"},
				{text: "", code: true},
			},
		},
		{
			name: "opening backtick only",
			s:    "`hello",
			want: []inlineSegment{{text: "hello", code: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseInline(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
