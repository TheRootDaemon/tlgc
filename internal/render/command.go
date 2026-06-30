package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/text"
)

// mappedWord pairs a single word from a command
// with the index of its originating Segment,
// so that the segment's style can be applied
// during line-by-line rendering.
type mappedWord struct {
	text         string
	segmentIndex int
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
		return nil
	}

	displayText := commandText(mappedWords)
	exampleIndent := strings.Repeat(" ", r.indent.Example)
	lines := wrapLines(
		r.output.LineLength,
		exampleIndent,
		displayText,
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
func (r *Renderer) renderCommandLine(
	w io.Writer,
	words []string,
	mappedWords []mappedWord,
	segments []Segment,
	indent string,
	wordOffset *int,
) error {
	_, err := io.WriteString(w, indent)
	if err != nil {
		return err
	}

	for j, word := range words {
		if *wordOffset >= len(mappedWords) {
			break
		}

		mapped := mappedWords[*wordOffset]
		segment := segments[mapped.segmentIndex]

		if _, err := io.WriteString(
			w,
			r.applyStyle(
				r.styleForSegment(&segment),
				word,
			),
		); err != nil {
			return err
		}

		if j < len(words)-1 {
			_, err := io.WriteString(w, " ")
			if err != nil {
				return err
			}
		}

		*wordOffset++
	}

	_, err = io.WriteString(w, "\n")
	return err
}

// mapWords flattens each Segment's DisplayText into individual words.
func mapWords(segments []Segment, optionStyle config.OptionStyle) []mappedWord {
	var mappedWords []mappedWord
	for i, segment := range segments {
		words := strings.FieldsSeq(segment.DisplayText(optionStyle))
		for word := range words {
			mappedWords = append(
				mappedWords,
				mappedWord{
					text:         word,
					segmentIndex: i,
				},
			)
		}
	}

	return mappedWords
}

// commandText joins the text fields of mapped words back
// into a single space-separated string, suitable for text wrapping.
func commandText(words []mappedWord) string {
	var b strings.Builder

	for i, word := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(word.text)
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
		return []string{displayText}
	}

	wrapped = text.Wrap(displayText, width, indent)
	return strings.Split(wrapped, "\n")
}
