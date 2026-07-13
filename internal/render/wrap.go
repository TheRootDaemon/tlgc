package render

import (
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/text"
)

// wrapText applies text wrapping to s and prepends indent to every line.
// Wrapping is controlled by r.output.LineLength;
// when LineLength ≤ 0 or s is empty,
// only the indent is prepended without wrapping.
func (r *Renderer) wrapText(s, indent string) string {
	if r.output.LineLength <= 0 || s == "" {
		logger.Trace(
			"no wrap (len=%d ll=%d)",
			len(s),
			r.output.LineLength,
		)
		return indent + s
	}

	wrapped := text.Wrap(s, r.output.LineLength, indent)
	logger.Trace(
		"input=%d ll=%d -> %d chars",
		len(s),
		r.output.LineLength,
		len(wrapped),
	)
	return indent + wrapped
}
