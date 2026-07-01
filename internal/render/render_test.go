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
	fullPageWant := "  tar\n\n" +
		"  archive utility.\n" +
		"  More information: https://example.org.\n" +
		"\n" +
		"  create archive\n" +
		"    tar cf archive.tar\n" +
		"\n" +
		"  extract\n" +
		"    tar xf archive.tar\n"

	tests := []struct {
		name        string
		platform    string
		renderer    *Renderer
		page        *Page
		want        string
		contains    []string
		notContains []string
		wantErr     string
	}{
		{
			name:     "nil page returns nil",
			renderer: &Renderer{},
			want:     "",
		},
		{
			name:     "empty page with only title",
			renderer: &Renderer{output: config.OutputConfig{ShowTitle: true}, indent: config.IndentConfig{Title: 2}},
			page:     &Page{Title: "tar"},
			want:     "  tar\n\n",
		},
		{
			name: "full page renders all sections in order",
			renderer: &Renderer{
				output: config.OutputConfig{ShowTitle: true},
				indent: config.IndentConfig{Title: 2, Description: 2, Bullet: 2, Example: 4},
			},
			page: &Page{
				Title:       "tar",
				Description: []string{"archive utility."},
				URL:         "https://example.org",
				Examples: []Example{
					{Description: "create archive", Command: "tar cf archive.tar"},
					{Description: "extract", Command: "tar xf archive.tar"},
				},
			},
			want: fullPageWant,
		},
		{
			name:     "platform title includes platform prefix",
			platform: "linux",
			renderer: &Renderer{
				output: config.OutputConfig{ShowTitle: true, PlatformTitle: true},
				indent: config.IndentConfig{Title: 2},
			},
			page: &Page{Title: "tar"},
			want: "  linux (tar)\n\n",
		},
		{
			name: "edit link rendered before title",
			renderer: &Renderer{
				output: config.OutputConfig{EditLink: true, ShowTitle: true},
				indent: config.IndentConfig{Title: 2},
			},
			page: &Page{
				Path:  "/pages/common/tar.md",
				Title: "tar",
			},
			contains: []string{
				"https://github.com/tldr-pages/tldr/edit/main/pages/common/tar.md\n",
				"  tar\n\n",
			},
		},
		{
			name:     "raw markdown mode writes file content",
			renderer: &Renderer{output: config.OutputConfig{RawMarkdown: true}},
			page:     &Page{},
			want:     "# test page\n\n> description.\n",
		},
		{
			name: "compact mode omits blank lines between examples",
			renderer: &Renderer{
				output: config.OutputConfig{Compact: true},
				indent: config.IndentConfig{Bullet: 2, Example: 4},
			},
			page: &Page{
				Examples: []Example{
					{Description: "first", Command: "echo a"},
					{Description: "second", Command: "echo b"},
				},
			},
			notContains: []string{"\n\n"},
		},
		{
			name: "title hidden when ShowTitle is false",
			renderer: &Renderer{
				output: config.OutputConfig{ShowTitle: false},
				indent: config.IndentConfig{Title: 2, Description: 2, Bullet: 2, Example: 4},
			},
			page: &Page{
				Title:       "tar",
				Description: []string{"desc."},
				Examples:    []Example{{Description: "ex", Command: "cmd"}},
			},
			notContains: []string{"tar"},
		},
		{
			name: "write error propagates",
			renderer: &Renderer{
				w:      &errorWriter{err: errors.New("write error")},
				output: config.OutputConfig{ShowTitle: true},
				indent: config.IndentConfig{Title: 2},
			},
			page:    &Page{Title: "tar"},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			if tt.renderer.w == nil {
				tt.renderer.w = &buf
			}

			if tt.renderer.output.RawMarkdown && tt.want != "" {
				path := filepath.Join(t.TempDir(), "page.md")
				require.NoError(t, os.WriteFile(path, []byte(tt.want), 0o644))
				tt.page.Path = path
			}

			err := tt.renderer.Render(tt.platform, tt.page)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			got := buf.String()

			if tt.want != "" {
				assert.Equal(t, tt.want, got)
			}
			for _, s := range tt.contains {
				assert.Contains(t, got, s)
			}
			for _, s := range tt.notContains {
				assert.NotContains(t, got, s)
			}
		})
	}
}

func TestRenderEditLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		compact bool
		url     string
		writer  io.Writer
		want    string
		wantErr string
	}{
		{
			name:    "non-compact adds newline",
			compact: false,
			url:     "https://example.com",
			want:    "https://example.com\n",
		},
		{
			name:    "compact omits newline",
			compact: true,
			url:     "https://example.com",
			want:    "https://example.com",
		},
		{
			name:    "write error",
			compact: false,
			url:     "url",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				output: config.OutputConfig{Compact: tt.compact},
			}

			var buf strings.Builder
			w := io.Writer(&buf)
			if tt.writer != nil {
				w = tt.writer
			}

			err := r.renderEditLink(w, tt.url)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestRenderRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		content    string
		writer     io.Writer
		want       string
		wantErr    string
		wantAnyErr bool
	}{
		{
			name:    "writes file content",
			content: "# hello",
			want:    "# hello",
		},
		{
			name:       "file not found",
			wantAnyErr: true,
		},
		{
			name:    "write error",
			content: "data",
			writer:  &errorWriter{err: errors.New("write error")},
			wantErr: "write error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/nonexistent/file.md"
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
			err := r.renderRaw(&Page{Path: path})

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
			want:   "  create archive\n    tar cf archive.tar\n",
		},
		{
			name:   "description only no command",
			ex:     Example{Description: "just a description"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			want:   "  just a description\n",
		},
		{
			name:   "hyphens enabled",
			ex:     Example{Description: "create archive", Command: "tar cf archive.tar"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			output: config.OutputConfig{ShowHyphens: true, ExamplePrefix: "- "},
			want:   "  - create archive\n    tar cf archive.tar\n",
		},
		{
			name:   "compact mode no blank line",
			ex:     Example{Description: "create archive", Command: "tar cf archive.tar"},
			indent: config.IndentConfig{Bullet: 2, Example: 4},
			output: config.OutputConfig{Compact: true},
			want:   "  create archive    tar cf archive.tar\n",
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

func TestBuildEditURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		url  string
		want string
	}{
		{
			name: "non-empty url returned as-is",
			path: "/pages/common/tar.md",
			url:  "https://example.com",
			want: "https://example.com",
		},
		{
			name: "empty url constructs from path",
			path: "/pages/common/tar.md",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/common/tar.md",
		},
		{
			name: "linux platform extracted correctly",
			path: "/pages/linux/apt.md",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/linux/apt.md",
		},
		{
			name: "windows platform extracted correctly",
			path: "/pages/windows/dir.md",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/windows/dir.md",
		},
		{
			name: "empty path and empty url returns empty",
			path: "",
			url:  "",
			want: "",
		},
		{
			name: "path without .md extension adds .md",
			path: "/pages/common/some-page",
			want: "https://github.com/tldr-pages/tldr/edit/main/pages/common/some-page.md",
		},
		{
			name: "url takes precedence over path",
			path: "/pages/common/tar.md",
			url:  "https://custom.com/edit",
			want: "https://custom.com/edit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildEditURL(tt.path, tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
