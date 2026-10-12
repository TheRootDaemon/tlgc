package render

import (
	"io"
	"strings"
)

// renderTitle writes the page title,
// styled with r.style.Title and
// indented by r.indent.Title spaces.
func (r *Renderer) renderTitle(title string) error {
	indent := strings.Repeat(" ", r.indent.Title)
	_, err := io.WriteString(
		r.w,
		r.applyStyle(
			r.style.Title,
			r.wrapText(title, indent),
		),
	)
	return err
}

// renderPageTitle renders the page title section for p.
//
// If title rendering is disabled
// or the page has no title, it does nothing.
// When PlatformTitle is enabled,
// the title is prefixed with "<platform>/".
// After rendering the title,
// it writes the appropriate trailing newline
// based on the configured output mode.
func (r *Renderer) renderPageTitle(platform string, p *Page) error {
	if !r.output.ShowTitle || p.Title == "" {
		return nil
	}

	title := p.Title
	if r.output.PlatformTitle && platform != "" {
		title = platform + "/" + title
	}

	if err := r.renderTitle(title); err != nil {
		return err
	}

	if r.output.Compact {
		return r.writeNewline()
	}

	return r.writeBlankLine()
}
