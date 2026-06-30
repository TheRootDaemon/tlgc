package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/text"
)

type mappedWord struct {
	text         string
	segmentIndex int
}

func (r *Renderer) renderCommand(w io.Writer, segments []Segment) {
	mappedWords := mapWords(segments, r.output.OptionStyle)
	if len(mappedWords) == 0 {
		return
	}

	displayText := displayText(mappedWords)
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

		r.renderCommandLine(
			w,
			words,
			mappedWords,
			segments,
			exampleIndent,
			&wordOffset,
		)
	}
}

func (r *Renderer) renderCommandLine(
	w io.Writer,
	words []string,
	mappedWords []mappedWord,
	segments []Segment,
	indent string,
	wordOffset *int,
) {
	io.WriteString(w, indent)

	for j, word := range words {
		if *wordOffset >= len(mappedWords) {
			break
		}

		mapped := mappedWords[*wordOffset]
		segment := segments[mapped.segmentIndex]

		io.WriteString(
			w,
			r.applyStyle(
				r.styleForSegment(&segment),
				word,
			),
		)

		if j < len(words)-1 {
			io.WriteString(w, " ")
		}

		*wordOffset++
	}

	io.WriteString(w, "\n")
}

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

func displayText(words []mappedWord) string {
	var b strings.Builder

	for i, word := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(word.text)
	}

	return b.String()
}

func wrapLines(
	width int,
	indent,
	displayTest string,
) []string {
	var wrapped string
	if width <= 0 {
		return []string{displayTest}
	}

	wrapped = text.Wrap(displayTest, width, indent)
	return strings.Split(wrapped, "\n")
}
