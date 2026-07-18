package render

import (
	"io"
	"strings"
)

// renderDescriptionLine writes one description line,
// styled with r.style.Description, indented,
// and followed by a newline.
func (r *Renderer) renderDescriptionLine(w io.Writer, text, indent string) error {
	if err := r.renderStyledInline(
		w,
		text,
		indent,
		r.style.Description,
		r.style.InlineCode,
	); err != nil {
		return err
	}

	_, err := io.WriteString(w, "\n")
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
		if err := r.renderDescriptionURL(w, url, indent); err != nil {
			return err
		}
	}

	if !r.output.Compact {
		_, err := io.WriteString(w, "\n")
		return err
	}

	return nil
}

// renderDescriptionURL writes the "More information: <url>." line,
// followed by a newline.
// The "More information: " prefix and trailing "." use r.style.Description;
// the URL itself uses r.style.URL.
// The line is indented and wrapped via wrapText.
func (r *Renderer) renderDescriptionURL(w io.Writer, url, indent string) error {
	var text strings.Builder
	text.WriteString(
		r.applyStyle(
			r.style.Description,
			"More information: ",
		),
	)
	text.WriteString(
		r.applyStyle(
			r.style.URL,
			url,
		),
	)
	text.WriteString(
		r.applyStyle(
			r.style.Description,
			".",
		),
	)

	wrappedText := r.wrapText(text.String(), indent)

	_, err := io.WriteString(
		w,
		wrappedText+"\n",
	)
	return err
}

// renderBulletLine writes one bullet item line,
// styled with r.style.Bullet and indented.
// No trailing newline is added.
func (r *Renderer) renderBulletLine(w io.Writer, text, indent string) error {
	return r.renderStyledInline(
		w,
		text,
		indent,
		r.style.Bullet,
		r.style.InlineCode,
	)
}
