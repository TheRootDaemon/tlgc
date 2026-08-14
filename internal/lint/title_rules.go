package lint

import (
	"strings"
	"unicode"
)

// checkTitleDescriptionSeparator enforces TLDR006.
//
// It reports an error if the page title and the first command
// description are not separated by at least one blank line.
func checkTitleDescriptionSeparator(p *parsedPage, r *Result) {
	if len(p.descriptions) == 0 {
		return
	}

	firstDescriptionNumber := p.descriptions[0].lineNumber

	// find the position of the title and first description in p.lines.
	titleIndex := -1
	descriptionIndex := -1
	for i, l := range p.lines {
		if l.lineNumber == p.titleLineNumber {
			titleIndex = i
		}
		if l.lineNumber == firstDescriptionNumber {
			descriptionIndex = i
		}
	}

	if titleIndex < 0 || descriptionIndex < 0 {
		return
	}

	// there should be at least one blank line between them
	// (check the line immediately before the description).
	if descriptionIndex-titleIndex < 2 || p.lines[descriptionIndex-1].kind != kindBlank {
		addError(r, "TLDR006", p.titleLineNumber)
	}
}

// checkTitleCharsacters enforces TLDR013.
//
// It reports an error if the page title contains characters
// outside the allowed character set or ends with a period,
// except for the special titles "." and " .".
func checkTitleCharacters(p *parsedPage, r *Result) {
	if p.title == "" {
		return
	}

	// check for characters outside the allowed set.
	for _, ch := range p.title {
		if !isValidTitleRune(ch) {
			addError(r, "TLDR013", p.titleLineNumber)
			return
		}
	}

	// title should not end with '.' unless it is '.' or ' .'
	if strings.HasSuffix(p.title, ".") &&
		p.title != "." &&
		p.title != " ." {
		addError(r, "TLDR013", p.titleLineNumber)
	}
}

// checkTitleHash enforces TLDR106.
//
// It reports an error if the page does not contain a title line,
// that is, a line beginning with '#'.
func checkTitleHash(p *parsedPage, r *Result) {
	// if the page has no title, error at the first non-blank line.
	for _, l := range p.lines {
		if l.kind == kindTitle {
			return
		}
	}
	for _, l := range p.lines {
		if l.kind != kindBlank {
			addError(r, "TLDR106", l.lineNumber)
			return
		}
	}
}

// isValidTitleRune reports whether ch is permitted in a page title.
//
// Valid title characters include Unicode letters and digits,
// underscores, spaces, and the punctuation
// permitted by the TLDR page format.
func isValidTitleRune(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsDigit(ch) ||
		ch == '_' ||
		strings.ContainsRune("+[]{}()!%,^~$:><|?.- ", ch)
}
