package render

import (
	"io"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/termcolor"
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

	logger.Debug("useColor=%t", r.useColor)

	return r
}

// Render writes a formatted tldr page to the Renderer's writer.
// platform is used only when PlatformTitle is enabled.
// Nil pages are silently ignored.
func (r *Renderer) Render(platform string, p *Page) error {
	if p == nil {
		logger.Trace("nil page, skipped")
		return nil
	}

	logger.Debug(
		"title=%q platform=%q raw=%t compact=%t edit=%d descs=%d examples=%d",
		p.Title,
		platform,
		r.output.RawMarkdown,
		r.output.Compact,
		r.output.EditLink,
		len(p.Description),
		len(p.Examples),
	)

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
			return err
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
