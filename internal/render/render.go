package render

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/pathutil"
	"github.com/TheRootDaemon/tlgc/termcolor"
	"github.com/TheRootDaemon/tlgc/text"
)

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

	if r.output.EditLink {
		if url := buildEditURL(p.Path, p.URL); url != "" {
			r.renderEditLink(r.w, url)
		}
	}

	if r.output.RawMarkdown {
		return r.renderRaw(p)
	}

	if r.output.ShowTitle && p.Title != "" {
		title := p.Title
		if r.output.PlatformTitle && platform != "" {
			title = platform + " (" + p.Title + ")"
		}

		r.renderTitle(r.w, title)
		io.WriteString(r.w, "\n\n")
	}

	r.renderDescriptions(r.w, p.Description, p.URL)

	for i, ex := range p.Examples {
		if i > 0 && !r.output.Compact {
			io.WriteString(r.w, "\n")
		}

		r.renderExample(r.w, ex)
	}

	return nil
}

func (r *Renderer) renderEditLink(w io.Writer, url string) {
	logger.Info("edit this page on GitHub")
	io.WriteString(w, url)

	if !r.output.Compact {
		io.WriteString(w, "\n")
	}
}

func (r *Renderer) renderRaw(p *Page) error {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return err
	}

	_, err = r.w.Write(data)
	return err
}

func (r *Renderer) renderTitle(w io.Writer, title string) {
	indent := strings.Repeat(" ", r.indent.Title)

	io.WriteString(
		w,
		r.applyStyle(
			r.style.Title,
			r.wrapText(title, indent),
		),
	)
}

func (r *Renderer) renderDescriptionLine(w io.Writer, text, indent string) {
	io.WriteString(
		w,
		r.applyStyle(
			r.style.Description,
			r.wrapText(text, indent),
		),
	)
	io.WriteString(w, "\n")
}

func (r *Renderer) renderDescriptions(w io.Writer, descs []string, url string) {
	if len(descs) == 0 && url == "" {
		return
	}

	indent := strings.Repeat(" ", r.indent.Description)

	for _, d := range descs {
		r.renderDescriptionLine(w, d, indent)
	}

	if url != "" {
		r.renderDescriptionLine(w, "More information: "+url+".", indent)
	}

	io.WriteString(w, "\n")
}

func (r *Renderer) renderBulletLine(w io.Writer, text, indent string) {
	io.WriteString(
		w,
		r.applyStyle(
			r.style.Bullet,
			r.wrapText(text, indent),
		),
	)
}

func (r *Renderer) renderExample(w io.Writer, ex Example) {
	indent := strings.Repeat(" ", r.indent.Bullet)
	desc := ex.Description

	if r.output.ShowHyphens {
		desc = r.output.ExamplePrefix + desc
	}

	r.renderBulletLine(w, desc, indent)

	if !r.output.Compact {
		io.WriteString(w, "\n")
	}

	if ex.Command != "" {
		segments := ParseCommand(ex.Command)
		r.renderCommand(w, segments)
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

func (r *Renderer) wrapText(s, indent string) string {
	if r.output.LineLength <= 0 || s == "" {
		return indent + s
	}

	wrapped := text.Wrap(s, r.output.LineLength, indent)
	return indent + wrapped
}
