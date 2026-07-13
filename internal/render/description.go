package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/logger"
)

// renderDescriptionLine writes one description line,
// styled with r.style.Description, indented,
// and followed by a newline.
func (r *Renderer) renderDescriptionLine(w io.Writer, text, indent string) error {
	logger.Trace("text=%q", text)
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

	logger.Trace("count=%d hasURL=%t", len(descs), url != "")

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
	logger.Trace("text=%q", text)
	_, err := io.WriteString(
		w,
		r.applyStyle(
			r.style.Bullet,
			r.wrapText(text, indent),
		),
	)

	return err
}
