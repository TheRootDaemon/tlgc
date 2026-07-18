package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/text"
)

// mappedWord pairs a single word from a command
// with the index of its originating Segment,
// so that the segment's style can be applied
// during line-by-line rendering.
//
// spaceAfter reports whether a space originally
// followed this word in the source,
// so commandText can reconstruct
// the display string faithfully
// without guessing punctuation spacing.
type mappedWord struct {
	followedBySpace bool
	segmentIndex    int
	text            string
}

// renderCommand writes a styled, wrapped command to w.
//
// It decomposes segments into word-level mappings,
// wraps the combined text to fit r.output.LineLength,
// and renders each wrapped line with per-word
// segment styling via renderCommandLine.
func (r *Renderer) renderCommand(w io.Writer, segments []Segment) error {
	mappedWords := mapWords(segments, r.output.OptionStyle)
	if len(mappedWords) == 0 {
		logger.Trace("no mapped words, skipped")
		return nil
	}

	displayText := commandText(mappedWords)
	exampleIndent := strings.Repeat(" ", r.indent.Example)
	lines := wrapLines(
		r.output.LineLength,
		exampleIndent,
		displayText,
	)

	logger.Trace(
		"segments=%d mappedWords=%d lines=%d",
		len(segments), len(mappedWords), len(lines),
	)

	wordOffset := 0

	for _, line := range lines {
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}

		if err := r.renderCommandLine(
			w,
			words,
			mappedWords,
			segments,
			exampleIndent,
			&wordOffset,
		); err != nil {
			return err
		}
	}

	return nil
}

// renderCommandLine writes one indented line of a command,
// applying the style of each word's originating Segment.
// wordOffset tracks the current position in mappedWords
// across multi-line rendering.
//
// Consecutive mappedWords with spaceAfter=false form a single
// whitespace-delimited field in the wrapped display text;
// each sub-word is rendered with its own segment's style.
func (r *Renderer) renderCommandLine(
	w io.Writer,
	words []string,
	mappedWords []mappedWord,
	segments []Segment,
	indent string,
	wordOffset *int,
) error {
	logger.Trace("words=%d offset=%d", len(words), *wordOffset)
	_, err := io.WriteString(w, indent)
	if err != nil {
		return err
	}

	for j := range words {
		if *wordOffset >= len(mappedWords) {
			break
		}

		start := *wordOffset
		for *wordOffset < len(mappedWords) {
			mw := mappedWords[*wordOffset]
			*wordOffset++
			if mw.followedBySpace {
				break
			}
		}

		for k := start; k < *wordOffset; k++ {
			mw := mappedWords[k]
			seg := segments[mw.segmentIndex]
			if _, err := io.WriteString(
				w,
				r.applyStyle(
					r.styleForSegment(&seg),
					mw.text,
				),
			); err != nil {
				return err
			}
		}

		if j < len(words)-1 {
			_, err := io.WriteString(w, " ")
			if err != nil {
				return err
			}
		}
	}

	_, err = io.WriteString(w, "\n")
	return err
}

// mapWords flattens each Segment's DisplayText into individual words
// and records whether each word was followed by whitespace.
func mapWords(segments []Segment, optionStyle config.OptionStyle) []mappedWord {
	var mappedWords []mappedWord
	for i, segment := range segments {
		text := segment.DisplayText(optionStyle)
		segmentWords := strings.Fields(text)
		for j, w := range segmentWords {
			mw := mappedWord{
				text:         w,
				segmentIndex: i,
			}

			isLast := j == len(segmentWords)-1

			if isLast {
				trailingSpace := hasTrailingSpace(text)
				leadingSpace := false
				if i+1 < len(segments) {
					next := segments[i+1].DisplayText(optionStyle)
					leadingSpace = hasLeadingSpace(next)
				}
				mw.followedBySpace = trailingSpace || leadingSpace
			} else {
				mw.followedBySpace = true
			}

			mappedWords = append(mappedWords, mw)
		}
	}

	return mappedWords
}

// commandText joins the text fields of mapped words back
// into a single space-separated string, suitable for text wrapping.
// A space is inserted only where the original source had whitespace
// (recorded in spaceAfter), so punctuation spacing is preserved
// faithfully without a hard-coded punctuation list.
func commandText(words []mappedWord) string {
	var b strings.Builder
	for i, w := range words {
		if i > 0 && words[i-1].followedBySpace {
			b.WriteByte(' ')
		}
		b.WriteString(w.text)
	}

	return b.String()
}

// wrapLines wraps displayText to fit within width columns.
// Continuation lines are prefixed with indent.
// If width ≤ 0 the text is returned as a single-element slice (no wrapping).
func wrapLines(
	width int,
	indent,
	displayText string,
) []string {
	var wrapped string
	if width <= 0 {
		logger.Trace("no wrap (width=%d)", width)
		return []string{displayText}
	}

	wrapped = text.Wrap(displayText, width, indent)
	lines := strings.Split(wrapped, "\n")
	logger.Trace(
		"width=%d input=%d output=%d lines",
		width,
		len(displayText),
		len(lines),
	)
	return lines
}

// hasTrailingSpace reports whether s ends with a space,
// indicating that whitespace originally followed the segment.
func hasTrailingSpace(s string) bool {
	return len(s) > 0 && s[len(s)-1] == ' '
}

// hasLeadingSpace reports whether s begins with a space,
// indicating that whitespace originally preceded the segment.
func hasLeadingSpace(s string) bool {
	return len(s) > 0 && s[0] == ' '
}
