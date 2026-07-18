package render

import "github.com/TheRootDaemon/tlgc/text"

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
