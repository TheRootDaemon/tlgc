package render

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/pathutil"
	"github.com/TheRootDaemon/tlgc/termcolor"
	"github.com/TheRootDaemon/tlgc/text"
)

type mappedWord struct {
	text         string
	segmentIndex int
}

type Renderer struct {
	useColor bool
	w        io.Writer
	style    config.StyleConfig
	output   config.OutputConfig
	indent   config.IndentConfig
}

type RenderOption func(*Renderer)

func WithColor(enabled bool) RenderOption {
	return func(r *Renderer) {
		r.useColor = enabled
	}
}

func WithWriter(w io.Writer) RenderOption {
	return func(r *Renderer) {
		r.w = w
	}
}

func WithStyle(style config.StyleConfig) RenderOption {
	return func(r *Renderer) {
		r.style = style
	}
}

func WithOutput(output config.OutputConfig) RenderOption {
	return func(r *Renderer) {
		r.output = output
	}
}

func WithIndent(indent config.IndentConfig) RenderOption {
	return func(r *Renderer) {
		r.indent = indent
	}
}

func New(w io.Writer, options ...RenderOption) *Renderer {
	r := &Renderer{
		useColor: termcolor.SupportsColor(),
		w:        w,
		style:    config.Style(),
		output:   config.Output(),
		indent:   config.Indent(),
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (r *Renderer) Render(platform string, p *Page) error {
	if p == nil {
		return nil
	}

	if r.output.RawMarkdown {
		data, err := os.ReadFile(p.Path)
		if err != nil {
			return err
		}

		_, err = r.w.Write(data)
		return err
	}

	if r.output.ShowTitle && p.Title != "" {
		title := p.Title
		if r.output.PlatformTitle && platform != "" {
			title = platform + " (" + p.Title + ")"
		}

		titleIndent := strings.Repeat(" ", r.indent.Title)
		io.WriteString(
			r.w,
			r.applyStyle(
				r.style.Title,
				r.wrapText(
					title,
					titleIndent,
				),
			),
		)
		io.WriteString(r.w, "\n\n")
	}

	if len(p.Description) > 0 {
		descIndent := strings.Repeat(" ", r.indent.Description)
		for _, desc := range p.Description {
			io.WriteString(
				r.w,
				r.applyStyle(
					r.style.Description,
					r.wrapText(
						desc,
						descIndent,
					),
				),
			)
			io.WriteString(r.w, "\n")
		}
		io.WriteString(r.w, "\n")
	}

	editURL := ""
	if r.output.EditLink {
		editURL = buildEditURL(p.Path, p.URL)
	}

	if p.URL != "" || editURL != "" {
		descIndent := strings.Repeat(" ", r.indent.Description)

		if p.URL != "" {
			urlText := "More information: <" + p.URL + ">"
			io.WriteString(
				r.w,
				r.applyStyle(
					r.style.URL,
					r.wrapText(urlText, descIndent),
				),
			)
			io.WriteString(r.w, "\n")
		}

		if editURL != "" {
			io.WriteString(
				r.w,
				r.applyStyle(
					r.style.URL,
					r.wrapText("Edit: "+editURL, descIndent),
				),
			)
			io.WriteString(r.w, "\n")
		}

		io.WriteString(r.w, "\n")
	}

	for i, ex := range p.Examples {
		if i > 0 && !r.output.Compact {
			io.WriteString(r.w, "\n")
		}

		bulletIndent := strings.Repeat(" ", r.indent.Bullet)
		desc := ex.Description
		if r.output.ShowHyphens {
			desc = r.output.ExamplePrefix + desc
		}

		io.WriteString(
			r.w,
			r.applyStyle(
				r.style.Bullet,
				r.wrapText(desc, bulletIndent),
			),
		)

		io.WriteString(r.w, "\n")

		if ex.Command != "" {
			segments := ParseCommand(ex.Command)
			r.renderCommand(r.w, segments)
		}
	}

	return nil
}

func (r *Renderer) applyStyle(s config.OutputStyle, t string) string {
	if !r.useColor {
		return t
	}

	return termcolor.Sprint(styleString(s), t)
}

func styleString(s config.OutputStyle) string {
	var parts []string
	if s.Bold {
		parts = append(parts, "bold")
	}
	if s.Dim {
		parts = append(parts, "dim")
	}
	if s.Italic {
		parts = append(parts, "italic")
	}
	if s.Underline {
		parts = append(parts, "underline")
	}
	if s.Strikethrough {
		parts = append(parts, "strikethrough")
	}
	if s.Color.Kind == config.ColorKindNamed && s.Color.Named != config.ColorDefault && s.Color.Named != "" {
		parts = append(parts, string(s.Color.Named))
	}
	if s.Background.Kind == config.ColorKindNamed && s.Background.Named != config.ColorDefault && s.Background.Named != "" {
		parts = append(parts, "on_"+string(s.Background.Named))
	}
	return strings.Join(parts, " ")
}

// wrapText wraps and indents plain text. It does not apply styling.
func (r *Renderer) wrapText(s, indent string) string {
	if r.output.LineLength <= 0 || s == "" {
		return indent + s
	}

	wrapped := text.Wrap(s, r.output.LineLength, indent)
	return indent + wrapped
}

func (r *Renderer) renderCommand(w io.Writer, segments []Segment) {
	var mappedWords []mappedWord
	for i, seg := range segments {
		words := strings.Fields(seg.DisplayText(r.output.OptionStyle))
		for _, word := range words {
			mappedWords = append(
				mappedWords,
				mappedWord{text: word, segmentIndex: i},
			)
		}
	}

	if len(mappedWords) == 0 {
		return
	}

	var b strings.Builder
	for i, mw := range mappedWords {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(mw.text)
	}
	displayText := b.String()

	exIndent := strings.Repeat(" ", r.indent.Example)

	var wrapped string
	if r.output.LineLength <= 0 {
		wrapped = displayText
	} else {
		wrapped = text.Wrap(displayText, r.output.LineLength, exIndent)
	}

	lines := strings.Split(wrapped, "\n")
	wordOffset := 0

	for _, line := range lines {
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}

		io.WriteString(w, exIndent)

		for j, word := range words {
			if wordOffset >= len(mappedWords) {
				break
			}

			mw := mappedWords[wordOffset]
			seg := segments[mw.segmentIndex]
			styled := r.applyStyle(r.styleForSegment(&seg), word)
			io.WriteString(w, styled)

			if j < len(words)-1 {
				io.WriteString(w, " ")
			}

			wordOffset++
		}

		io.WriteString(w, "\n")
	}
}

func (r *Renderer) styleForSegment(s *Segment) config.OutputStyle {
	switch s.Kind {
	case Text:
		return r.style.Example
	case Placeholder:
		return r.style.Placeholder
	case Option:
		return r.style.Placeholder
	default:
		return r.style.Example
	}
}

func buildEditURL(path, url string) string {
	if url != "" {
		return url
	}

	if path != "" {
		page := pathutil.PageName(path)
		platform := pathutil.PagePlatform(path)
		return fmt.Sprintf(
			"https://github.com/tldr-pages/tldr/edit/main/pages/%s/%s.md",
			platform,
			page,
		)
	}

	return ""
}
