package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/logger"
)

// renderTitle writes the page title,
// styled with r.style.Title and
// indented by r.indent.Title spaces.
func (r *Renderer) renderTitle(w io.Writer, title string) error {
	indent := strings.Repeat(" ", r.indent.Title)
	logger.Trace("title=%q indent=%d", title, r.indent.Title)

	_, err := io.WriteString(
		w,
		r.applyStyle(
			r.style.Title,
			r.wrapText(title, indent),
		),
	)
	return err
}
