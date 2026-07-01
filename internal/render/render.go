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

// Renderer renders parsed tldr pages to a writer
// with optional ANSI color and configurable styles, indentation, and output settings.
type Renderer struct {
	useColor bool                // whether ANSI color output is enabled
	w        io.Writer           // destination for rendered output
	style    config.StyleConfig  // style configuration for each page element
	output   config.OutputConfig // output visibility and formatting options
	indent   config.IndentConfig // indentation per section
}

// RenderOption configures a Renderer.
type RenderOption func(*Renderer)

// WithColor enables or disables ANSI color output.
func WithColor(enabled bool) RenderOption {
	return func(r *Renderer) {
		r.useColor = enabled
	}
}

// WithWriter sets the output writer for the Renderer.
func WithWriter(w io.Writer) RenderOption {
	return func(r *Renderer) {
		r.w = w
	}
}

// WithStyle replaces the default style configuration for all page elements.
func WithStyle(style config.StyleConfig) RenderOption {
	return func(r *Renderer) {
		r.style = style
	}
}

// WithOutput replaces the default output configuration
// (title visibility, hyphens, edit link, line length, etc.).
func WithOutput(output config.OutputConfig) RenderOption {
	return func(r *Renderer) {
		r.output = output
	}
}

// WithIndent replaces the default indentation configuration
// for each section (title, description, bullet, example).
func WithIndent(indent config.IndentConfig) RenderOption {
	return func(r *Renderer) {
		r.indent = indent
	}
}

// New creates a Renderer that writes to w.
//
// Defaults from the active config are used for style, output, and indentation;
// options may override any of these.
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

// Render writes a formatted tldr page to the Renderer's writer.
// platform is used only when PlatformTitle is enabled.
// Nil pages are silently ignored.
func (r *Renderer) Render(platform string, p *Page) error {
	if p == nil {
		return nil
	}

	if r.output.EditLink {
		if url := buildEditURL(p.Path, p.URL); url != "" {
			if err := r.renderEditLink(r.w, url); err != nil {
				return err
			}
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

		if err := r.renderTitle(r.w, title); err != nil {
			return err
		}

		_, err := io.WriteString(r.w, "\n\n")
		if err != nil {
			return nil
		}
	}

	if err := r.renderDescriptions(r.w, p.Description, p.URL); err != nil {
		return err
	}

	for i, ex := range p.Examples {
		if i > 0 && !r.output.Compact {
			_, err := io.WriteString(r.w, "\n")
			if err != nil {
				return err
			}
		}

		if err := r.renderExample(r.w, ex); err != nil {
			return err
		}
	}

	return nil
}

// renderEditLink writes the edit URL to w, with a trailing newline
// unless output is in compact mode.
func (r *Renderer) renderEditLink(w io.Writer, url string) error {
	logger.Info("edit this page on GitHub")
	_, err := io.WriteString(w, url)
	if err != nil {
		return err
	}

	if !r.output.Compact {
		_, err := io.WriteString(w, "\n")
		return err
	}

	return nil
}

// renderRaw reads the raw markdown file at p.Path and
// writes it to the Renderer's writer.
func (r *Renderer) renderRaw(p *Page) error {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return err
	}

	_, err = r.w.Write(data)
	return err
}

// renderTitle writes the page title,
// styled with r.style.Title and
// indented by r.indent.Title spaces.
func (r *Renderer) renderTitle(w io.Writer, title string) error {
	indent := strings.Repeat(" ", r.indent.Title)

	_, err := io.WriteString(
		w,
		r.applyStyle(
			r.style.Title,
			r.wrapText(title, indent),
		),
	)
	return err
}

// renderDescriptionLine writes one description line,
// styled with r.style.Description, indented,
// and followed by a newline.
func (r *Renderer) renderDescriptionLine(w io.Writer, text, indent string) error {
	_, err := io.WriteString(
		w,
		r.applyStyle(
			r.style.Description,
			r.wrapText(text, indent),
		),
	)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "\n")
	return err
}

// renderDescriptions writes all description lines
// followed by the "More information" URL (if set),
// each indented by r.indent.Description.
// Writes a trailing blank line after descriptions.
func (r *Renderer) renderDescriptions(w io.Writer, descs []string, url string) error {
	if len(descs) == 0 && url == "" {
		return nil
	}

	indent := strings.Repeat(" ", r.indent.Description)

	for _, d := range descs {
		if err := r.renderDescriptionLine(w, d, indent); err != nil {
			return err
		}
	}

	if url != "" {
		if err := r.renderDescriptionLine(w, "More information: "+url+".", indent); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "\n")
	return err
}

// renderBulletLine writes one bullet item line,
// styled with r.style.Bullet and indented.
// No trailing newline is added.
func (r *Renderer) renderBulletLine(w io.Writer, text, indent string) error {
	_, err := io.WriteString(
		w,
		r.applyStyle(
			r.style.Bullet,
			r.wrapText(text, indent),
		),
	)

	return err
}

// renderExample writes one example,
// a bullet line for the description (prefixed with ExamplePrefix when ShowHyphens is set)
// followed by the styled command text on the next line.
// In compact mode the blank line between bullet and command is omitted.
func (r *Renderer) renderExample(w io.Writer, ex Example) error {
	indent := strings.Repeat(" ", r.indent.Bullet)
	desc := ex.Description

	if r.output.ShowHyphens {
		desc = r.output.ExamplePrefix + desc
	}

	if err := r.renderBulletLine(w, desc, indent); err != nil {
		return err
	}

	if !r.output.Compact {
		_, err := io.WriteString(w, "\n")
		if err != nil {
			return err
		}
	}

	if ex.Command != "" {
		segments := ParseCommand(ex.Command)
		if err := r.renderCommand(w, segments); err != nil {
			return err
		}
	}

	return nil
}

// buildEditURL returns the GitHub edit URL for a tldr page.
// If url is non-empty it is returned as-is (the page has a custom source).
// Otherwise the URL is constructed from the page's file path.
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

// wrapText applies text wrapping to s and prepends indent to every line.
// Wrapping is controlled by r.output.LineLength;
// when LineLength ≤ 0 or s is empty,
// only the indent is prepended without wrapping.
func (r *Renderer) wrapText(s, indent string) string {
	if r.output.LineLength <= 0 || s == "" {
		return indent + s
	}

	wrapped := text.Wrap(s, r.output.LineLength, indent)
	return indent + wrapped
}
