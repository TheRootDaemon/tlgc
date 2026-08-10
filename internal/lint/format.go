package lint

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Format renders a canonical tldr page from content.
// It mirrors the reference linter.format function in tldr-lint:
// titles are kept as-is,
// descriptions are capitalized
// and forced to end in a period,
// example descriptions are capitalized
// and forced to end in a colon,
// and commands are re-emitted verbatim.
//
// It returns the empty string for empty or whitespace-only content.
func Format(content string) string {
	page := parse(content)
	if page == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString("# ")
	b.WriteString(page.title)
	b.WriteString("\n\n")

	for _, line := range page.descriptions {
		if isInfoLink(line.content) {
			b.WriteString("> More information: ")
			b.WriteString(angleBracketedURL(line.content))
			b.WriteString(".\n")
			continue
		}
		b.WriteString("> ")
		b.WriteString(formatDescription(line.content))
		b.WriteString("\n")
	}

	for _, section := range page.exampleSections {
		b.WriteString("\n- ")
		b.WriteString(formatExampleDescription(section.description))
		b.WriteString("\n\n")
		for _, command := range section.commands {
			b.WriteString("`")
			b.WriteString(command.content)
			b.WriteString("`\n")
		}
	}

	return b.String()
}

// formatDescription capitalizes the first character
// and forces a trailing period,
// mirroring the reference lexer which strips
// a single trailing punctuation character first.
func formatDescription(s string) string {
	s = stripTrailing(s, ".,;!?")
	return upperFirst(s) + "."
}

// formatExampleDescription capitalizes the first character
// and forces a trailing colon,
// mirroring the reference lexer which strips
// a single trailing punctuation character first.
func formatExampleDescription(s string) string {
	s = stripTrailing(s, ".:,;")
	return upperFirst(s) + ":"
}

// upperFirst returns s with its first character uppercased.
func upperFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// stripTrailing removes a single trailing character
// of s when it is one of cut.
func stripTrailing(s, cut string) string {
	if n := len(s); n > 0 &&
		strings.ContainsRune(cut, rune(s[n-1])) {
		return s[:n-1]
	}
	return s
}

// angleBracketedURL returns the <...> substring
// of an information link line,
// or the empty string if absent.
func angleBracketedURL(line string) string {
	openIndex := strings.Index(line, "<")
	if openIndex < 0 ||
		!strings.Contains(line[openIndex:], ">") {
		return ""
	}

	closeIndex := strings.Index(line[openIndex:], ">")
	return line[openIndex : openIndex+closeIndex+1]
}
