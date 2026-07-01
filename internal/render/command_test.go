package render

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRenderCommand(t *testing.T) {
	t.Parallel()

	noColorRenderer := &Renderer{
		useColor: false,
		output: config.OutputConfig{
			OptionStyle: config.OptionStyleLong,
			LineLength:  0,
		},
		indent: config.IndentConfig{
			Example: 4,
		},
	}

	tests := []struct {
		name     string
		renderer *Renderer
		raw      string
		want     string
	}{
		{
			name:     "empty segments writes nothing",
			renderer: noColorRenderer,
			raw:      "",
			want:     "",
		},
		{
			name:     "simple text command",
			renderer: noColorRenderer,
			raw:      "tar cf archive.tar",
			want:     "    tar cf archive.tar\n",
		},
		{
			name:     "command with placeholders",
			renderer: noColorRenderer,
			raw:      "tar cf {{archive.tar}} {{dest}}",
			want:     "    tar cf archive.tar dest\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := ParseCommand(tt.raw)
			var buf strings.Builder
			err := tt.renderer.renderCommand(&buf, segments)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}

	t.Run("wrapping produces multiple lines", func(t *testing.T) {
		r := &Renderer{
			useColor: false,
			output: config.OutputConfig{
				OptionStyle: config.OptionStyleLong,
				LineLength:  15,
			},
			indent: config.IndentConfig{
				Example: 4,
			},
		}
		var buf strings.Builder
		err := r.renderCommand(&buf, ParseCommand("some very long command"))
		assert.NoError(t, err)
		assert.Equal(t, "    some very long\n    command\n", buf.String())
	})

	t.Run("option rendered with short style", func(t *testing.T) {
		r := &Renderer{
			useColor: false,
			output: config.OutputConfig{
				OptionStyle: config.OptionStyleShort,
				LineLength:  0,
			},
			indent: config.IndentConfig{
				Example: 4,
			},
		}
		var buf strings.Builder
		err := r.renderCommand(&buf, ParseCommand("cmd {{[-s|--long]}}"))
		assert.NoError(t, err)
		assert.Equal(t, "    cmd -s\n", buf.String())
	})

	t.Run("option rendered with combined style", func(t *testing.T) {
		r := &Renderer{
			useColor: false,
			output: config.OutputConfig{
				OptionStyle: config.OptionStyleCombined,
				LineLength:  0,
			},
			indent: config.IndentConfig{
				Example: 4,
			},
		}
		var buf strings.Builder
		err := r.renderCommand(&buf, ParseCommand("cmd {{[-s|--long]}}"))
		assert.NoError(t, err)
		assert.Equal(t, "    cmd [-s|--long]\n", buf.String())
	})

	t.Run("colorized output contains ANSI sequences", func(t *testing.T) {
		r := &Renderer{
			useColor: true,
			style:    config.DefaultStyleConfig(),
			output: config.OutputConfig{
				OptionStyle: config.OptionStyleLong,
				LineLength:  0,
			},
			indent: config.IndentConfig{
				Example: 4,
			},
		}
		var buf strings.Builder
		err := r.renderCommand(&buf, ParseCommand("echo hello"))
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "\x1b[36m")
		assert.Contains(t, output, "\x1b[0m")
		assert.Contains(t, output, "echo")
		assert.Contains(t, output, "hello")
	})

	t.Run("error from renderCommandLine propagates", func(t *testing.T) {
		r := &Renderer{
			useColor: false,
			output: config.OutputConfig{
				OptionStyle: config.OptionStyleLong,
				LineLength:  0,
			},
			indent: config.IndentConfig{
				Example: 4,
			},
		}
		err := r.renderCommand(&errorWriter{err: errors.New("write error")}, ParseCommand("echo hi"))
		assert.ErrorContains(t, err, "write error")
	})
}

func TestMapWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		segments    []Segment
		optionStyle config.OptionStyle
		want        []mappedWord
	}{
		{
			name:        "empty segments",
			segments:    nil,
			optionStyle: config.OptionStyleLong,
			want:        nil,
		},
		{
			name:        "single text segment",
			segments:    []Segment{{Kind: Text, Text: "hello"}},
			optionStyle: config.OptionStyleLong,
			want:        []mappedWord{{text: "hello", segmentIndex: 0}},
		},
		{
			name:        "text segment with multiple words",
			segments:    []Segment{{Kind: Text, Text: "tar cf archive.tar"}},
			optionStyle: config.OptionStyleLong,
			want: []mappedWord{
				{text: "tar", segmentIndex: 0},
				{text: "cf", segmentIndex: 0},
				{text: "archive.tar", segmentIndex: 0},
			},
		},
		{
			name: "multiple text segments",
			segments: []Segment{
				{Kind: Text, Text: "mv "},
				{Kind: Text, Text: "src "},
				{Kind: Text, Text: "dst"},
			},
			optionStyle: config.OptionStyleLong,
			want: []mappedWord{
				{text: "mv", segmentIndex: 0},
				{text: "src", segmentIndex: 1},
				{text: "dst", segmentIndex: 2},
			},
		},
		{
			name: "option with short style",
			segments: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-s", Long: "--long"},
			},
			optionStyle: config.OptionStyleShort,
			want: []mappedWord{
				{text: "cmd", segmentIndex: 0},
				{text: "-s", segmentIndex: 1},
			},
		},
		{
			name: "option with long style",
			segments: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-s", Long: "--long"},
			},
			optionStyle: config.OptionStyleLong,
			want: []mappedWord{
				{text: "cmd", segmentIndex: 0},
				{text: "--long", segmentIndex: 1},
			},
		},
		{
			name: "option with combined style",
			segments: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-s", Long: "--long"},
			},
			optionStyle: config.OptionStyleCombined,
			want: []mappedWord{
				{text: "cmd", segmentIndex: 0},
				{text: "[-s|--long]", segmentIndex: 1},
			},
		},
		{
			name: "placeholder segment",
			segments: []Segment{
				{Kind: Text, Text: "echo "},
				{Kind: Placeholder, Text: "hello"},
			},
			optionStyle: config.OptionStyleLong,
			want: []mappedWord{
				{text: "echo", segmentIndex: 0},
				{text: "hello", segmentIndex: 1},
			},
		},
		{
			name: "mixed text option and placeholder",
			segments: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-o", Long: "--output"},
				{Kind: Text, Text: " "},
				{Kind: Placeholder, Text: "file"},
			},
			optionStyle: config.OptionStyleShort,
			want: []mappedWord{
				{text: "cmd", segmentIndex: 0},
				{text: "-o", segmentIndex: 1},
				{text: "file", segmentIndex: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapWords(tt.segments, tt.optionStyle)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommandText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		words []mappedWord
		want  string
	}{
		{
			name:  "empty",
			words: nil,
			want:  "",
		},
		{
			name: "single word",
			words: []mappedWord{
				{text: "hello"},
			},
			want: "hello",
		},
		{
			name: "multiple words",
			words: []mappedWord{
				{text: "tar"},
				{text: "cf"},
				{text: "archive.tar"},
			},
			want: "tar cf archive.tar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandText(tt.words)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWrapLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		width  int
		indent string
		text   string
		want   []string
	}{
		{
			name:   "width zero returns single line",
			width:  0,
			indent: "    ",
			text:   "some very long command that exceeds any reasonable width",
			want:   []string{"some very long command that exceeds any reasonable width"},
		},
		{
			name:   "width negative returns single line",
			width:  -1,
			indent: "    ",
			text:   "short",
			want:   []string{"short"},
		},
		{
			name:   "width exceeds text returns single line",
			width:  100,
			indent: "    ",
			text:   "short text",
			want:   []string{"short text"},
		},
		{
			name:   "width less than text wraps with indent",
			width:  5,
			indent: "  ",
			text:   "a b c d e",
			want: []string{
				"a b c",
				"  d e",
			},
		},
		{
			name:   "empty text returns single empty string",
			width:  80,
			indent: "    ",
			text:   "",
			want:   []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapLines(tt.width, tt.indent, tt.text)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderCommandLine(t *testing.T) {
	t.Parallel()

	r := &Renderer{
		useColor: false,
	}

	tests := []struct {
		name        string
		words       []string
		mappedWords []mappedWord
		segments    []Segment
		indent      string
		want        string
		wantOffset  int
		writer      io.Writer
		wantErr     string
	}{
		{
			name:  "single word",
			words: []string{"hello"},
			mappedWords: []mappedWord{
				{text: "hello", segmentIndex: 0},
			},
			segments: []Segment{
				{Kind: Text, Text: "hello"},
			},
			indent:     "    ",
			want:       "    hello\n",
			wantOffset: 1,
		},
		{
			name:  "multiple words",
			words: []string{"tar", "cf", "archive.tar"},
			mappedWords: []mappedWord{
				{text: "tar", segmentIndex: 0},
				{text: "cf", segmentIndex: 0},
				{text: "archive.tar", segmentIndex: 0},
			},
			segments: []Segment{
				{Kind: Text, Text: "tar cf archive.tar"},
			},
			indent:     "    ",
			want:       "    tar cf archive.tar\n",
			wantOffset: 3,
		},
		{
			name:  "break when mapped words exhausted",
			words: []string{"a", "b", "c"},
			mappedWords: []mappedWord{
				{text: "a", segmentIndex: 0},
			},
			segments: []Segment{
				{Kind: Text, Text: "a"},
			},
			indent:     "    ",
			want:       "    a \n",
			wantOffset: 1,
		},
		{
			name:  "write error",
			words: []string{"hello"},
			mappedWords: []mappedWord{
				{text: "hello", segmentIndex: 0},
			},
			segments: []Segment{
				{Kind: Text, Text: "hello"},
			},
			indent:  "    ",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := 0

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderCommandLine(
				w,
				tt.words,
				tt.mappedWords,
				tt.segments,
				tt.indent,
				&offset,
			)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantOffset, offset)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
