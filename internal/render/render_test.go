package render

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []RenderOption
		want bool
	}{
		{name: "overrides to false", opts: []RenderOption{WithColor(false)}, want: false},
		{name: "overrides to true", opts: []RenderOption{WithColor(true)}, want: true},
		{name: "last option wins", opts: []RenderOption{WithColor(false), WithColor(true)}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			r := New(&buf, tt.opts...)
			assert.Equal(t, tt.want, r.useColor)
		})
	}
}

func TestWithWriter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		writer io.Writer
	}{
		{name: "os.Stdout", writer: os.Stdout},
		{name: "nil writer", writer: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(&strings.Builder{}, WithWriter(tt.writer))
			assert.Equal(t, tt.writer, r.w)
		})
	}
}

func TestWithStyle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		style config.StyleConfig
	}{
		{
			name: "custom style",
			style: config.StyleConfig{
				Title: config.OutputStyle{
					Bold:  true,
					Color: config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorRed},
				},
			},
		},
		{
			name:  "zero value style",
			style: config.StyleConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			r := New(&buf, WithStyle(tt.style))
			assert.Equal(t, tt.style, r.style)
		})
	}
}

func TestWithOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		output config.OutputConfig
	}{
		{
			name: "custom output",
			output: config.OutputConfig{
				ShowTitle:  false,
				LineLength: 50,
				Compact:    true,
			},
		},
		{
			name:   "zero value output",
			output: config.OutputConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			r := New(&buf, WithOutput(tt.output))
			assert.Equal(t, tt.output, r.output)
		})
	}
}

func TestWithIndent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		indent config.IndentConfig
	}{
		{
			name: "custom indent",
			indent: config.IndentConfig{
				Title:       0,
				Description: 1,
				Bullet:      2,
				Example:     3,
			},
		},
		{
			name:   "zero value indent",
			indent: config.IndentConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			r := New(&buf, WithIndent(tt.indent))
			assert.Equal(t, tt.indent, r.indent)
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("defaults with no config loaded", func(t *testing.T) {
		config.ResetForTesting()
		t.Cleanup(config.ResetForTesting)
		t.Setenv("NO_COLOR", "1")

		var buf strings.Builder
		r := New(&buf)

		assert.False(t, r.useColor)
		assert.Equal(t, &buf, r.w)
		assert.Equal(t, config.DefaultStyleConfig(), r.style)
		assert.Equal(t, config.DefaultOutputConfig(), r.output)
		assert.Equal(t, config.DefaultIndentConfig(), r.indent)
	})

	t.Run("picks up custom config", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.toml")
		content := `
[output]
show_title = false
line_length = 60

[indent]
title = 4
example = 6
`
		require.NoError(t, os.WriteFile(cfgPath, []byte(content), 0o644))
		t.Setenv("TLGC_CONFIG", cfgPath)
		t.Setenv("NO_COLOR", "1")
		config.ResetForTesting()
		require.NoError(t, config.Initialize())
		t.Cleanup(config.ResetForTesting)

		var buf strings.Builder
		r := New(&buf)

		assert.Equal(t, 4, r.indent.Title)
		assert.Equal(t, 6, r.indent.Example)
		assert.Equal(t, 2, r.indent.Description)
		assert.Equal(t, 2, r.indent.Bullet)
		assert.False(t, r.output.ShowTitle)
		assert.Equal(t, 60, r.output.LineLength)
		assert.False(t, r.output.ShowHyphens)
	})

	t.Run("color disabled when TERM is dumb", func(t *testing.T) {
		t.Setenv("TERM", "dumb")

		var buf strings.Builder
		r := New(&buf)

		assert.False(t, r.useColor)
	})
}

func TestRender(t *testing.T) {
	t.Parallel()

	fullPageWant := "\n" +
		"  tar\n\n" +
		"  archive utility.\n" +
		"  More information: https://example.org.\n" +
		"\n" +
		"  create archive\n\n" +
		"    tar cf archive.tar\n" +
		"\n" +
		"  extract\n\n" +
		"    tar xf archive.tar\n" +
		"\n"

	compactFullPageWant := "" +
		"  tar\n" +
		"  archive utility.\n" +
		"  More information: https://example.org.\n" +
		"  create archive\n" +
		"    tar cf archive.tar\n" +
		"  extract\n" +
		"    tar xf archive.tar\n"

	fullPage := &Page{
		Title:       "tar",
		Description: []string{"archive utility."},
		URL:         "https://example.org",
		Examples: []Example{
			{Description: "create archive", Command: "tar cf archive.tar"},
			{Description: "extract", Command: "tar xf archive.tar"},
		},
	}

	tests := []struct {
		name     string
		renderer *Renderer
		page     *Page
		want     string
		wantErr  string
	}{
		{
			name:     "nil page returns nil",
			renderer: &Renderer{},
			want:     "",
		},
		{
			name: "full page renders all sections in order",
			renderer: &Renderer{
				output: config.OutputConfig{ShowTitle: true},
				indent: config.IndentConfig{Title: 2, Description: 2, Bullet: 2, Example: 4},
			},
			page: fullPage,
			want: fullPageWant,
		},
		{
			name: "full page compact mode",
			renderer: &Renderer{
				output: config.OutputConfig{ShowTitle: true, Compact: true},
				indent: config.IndentConfig{Title: 2, Description: 2, Bullet: 2, Example: 4},
			},
			page: fullPage,
			want: compactFullPageWant,
		},
		{
			name:     "raw markdown mode writes content from RawContent",
			renderer: &Renderer{output: config.OutputConfig{RawMarkdown: true}},
			page:     &Page{RawContent: "# test page\n\n> description.\n"},
			want:     "\n# test page\n\n> description.\n",
		},
		{
			name:     "leading and trailing newlines in non-compact mode",
			renderer: &Renderer{output: config.OutputConfig{ShowTitle: false}},
			page:     &Page{},
			want:     "\n\n",
		},
		{
			name:     "compact suppresses leading and trailing newlines",
			renderer: &Renderer{output: config.OutputConfig{Compact: true}},
			page:     &Page{},
			want:     "",
		},
		{
			name: "write error on leading newline propagates",
			renderer: &Renderer{
				w:      &errorWriter{err: errors.New("write error")},
				output: config.OutputConfig{},
			},
			page:    &Page{},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			if tt.renderer.w == nil {
				tt.renderer.w = &buf
			}

			err := tt.renderer.Render("", tt.page)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			if tt.want != "" {
				assert.Equal(t, tt.want, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}
