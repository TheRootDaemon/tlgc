package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
)

// inlineSegment is one piece of inline content
// within a description or bullet line.
// code is true when the text was wrapped in backticks
// and should be styled with the InlineCode style.
type inlineSegment struct {
	code bool
	text string
}

// renderStyledInline writes text,
// splitting backtick-delimited code spans
// so plain segments use baseStyle and
// code segments use codeStyle.
// When no backticks are present,
// it falls through to the single-segment fast path (applyStyle + wrapText).
func (r *Renderer) renderStyledInline(
	w io.Writer,
	text,
	indent string,
	baseStyle,
	codeStyle config.OutputStyle,
) error {
	segments := parseInline(text)
	if len(segments) == 1 {
		_, err := io.WriteString(
			w,
			r.applyStyle(
				baseStyle,
				r.wrapText(text, indent),
			),
		)
		return err
	}

	wrappedLines := r.wrapInlineSegments(segments, indent)

	for i, line := range wrappedLines {
		if i > 0 {
			_, err := io.WriteString(w, "\n")
			if err != nil {
				return err
			}
		}

		if err := r.renderStyledLine(
			w,
			line,
			segments,
			baseStyle,
			codeStyle,
		); err != nil {
			return err
		}
	}

	return nil
}

// wrapInlineSegments joins all segment text into a single string,
// wraps it via wrapText, and returns the individual lines.
func (r *Renderer) wrapInlineSegments(
	segments []inlineSegment,
	indent string,
) []string {
	var plainText strings.Builder
	for _, segment := range segments {
		plainText.WriteString(segment.text)
	}

	wrapped := r.wrapText(plainText.String(), indent)
	return strings.Split(wrapped, "\n")
}

// renderStyledLine writes one wrapped line,
// applying baseStyle or codeStyle to each segment
// based on its code field.
// Segments are matched to the line text
// by position to avoid false matches from repeated words.
func (r *Renderer) renderStyledLine(
	w io.Writer,
	line string,
	segments []inlineSegment,
	baseStyle,
	codeStyle config.OutputStyle,
) error {
	if line == "" {
		return nil
	}

	remaining := line
	for _, segment := range segments {
		if segment.text == "" {
			continue
		}

		if remaining == "" {
			break
		}

		pos := strings.Index(remaining, segment.text)
		if pos < 0 {
			break
		}
		if pos > 0 {
			before := remaining[:pos]
			remaining = remaining[pos:]
			_, err := io.WriteString(
				w,
				r.applyStyle(baseStyle, before),
			)
			if err != nil {
				return err
			}
		}
		if len(segment.text) > len(remaining) {
			break
		}
		piece := remaining[:len(segment.text)]
		remaining = remaining[len(segment.text):]
		if segment.code {
			_, err := io.WriteString(w, r.applyStyle(codeStyle, piece))
			if err != nil {
				return err
			}
		} else {
			_, err := io.WriteString(w, r.applyStyle(baseStyle, piece))
			if err != nil {
				return err
			}
		}
	}
	if remaining != "" {
		_, err := io.WriteString(w, r.applyStyle(baseStyle, remaining))
		if err != nil {
			return err
		}
	}

	return nil
}

// parseInline splits s by backtick-delimited code spans.
// Even-index parts are plain text, odd-index parts are code spans.
// Returns at minimum one segment (the whole text as non-code).
func parseInline(s string) []inlineSegment {
	if !strings.Contains(s, "`") {
		return []inlineSegment{
			{
				text: s,
			},
		}
	}

	parts := strings.Split(s, "`")
	segments := make([]inlineSegment, 0, len(parts))

	for i, p := range parts {
		if p == "" && i%2 == 0 {
			continue
		}
		segments = append(
			segments,
			inlineSegment{
				text: p,
				code: i%2 == 1,
			},
		)
	}

	// drop trailing empty text segments
	for len(segments) > 0 &&
		!segments[len(segments)-1].code &&
		segments[len(segments)-1].text == "" {
		segments = segments[:len(segments)-1]
	}

	if len(segments) == 0 {
		segments = []inlineSegment{
			{
				text: s,
			},
		}
	}

	return segments
}
