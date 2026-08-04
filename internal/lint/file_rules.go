package lint

import "strings"

// checkLeadingWhitespace enforces TLDR001.
//
// It reports an error if the first non-blank line of the page
// starts with a space or a tab,
// or if the page starts with blank lines at all.
func checkLeadingWhitespace(p *parsedPage, r *Result) {
	for i, l := range p.lines {
		if l.kind == kindBlank {
			continue
		}

		// first non-blank line: check for leading space/tab on the line itself.
		if len(l.rawLine) > 0 &&
			(l.rawLine[0] == ' ' || l.rawLine[0] == '\t') {
			addError(r, "TLDR001", l.lineNumber)
		} else if i > 0 {
			// leading blank line triggers TLDR001.
			addError(r, "TLDR001", 0)
		}

		break
	}
}

// checkSpaceAfterPrefix enforces TLDR002.
//
// It reports an error if a title, description, or example description
// line does not have exactly one space
// after its marker ('#', '>', '-').
func checkSpaceAfterPrefix(p *parsedPage, r *Result) {
	for _, l := range p.lines {
		switch l.kind {
		case kindTitle, kindDescription, kindExampleDesc:
			if len(l.rawLine) > 1 && l.rawLine[1] != ' ' {
				addError(r, "TLDR002", l.lineNumber)
			}
		}
	}
}

// checkNoTrailingWhitespaceAtEOF enforces TLDR008.
//
// It reports an error if the page ends with trailing whitespace,
// that is, with more than the single terminating newline.
func checkNoTrailingWhitespaceAtEOF(p *parsedPage, r *Result) {
	trimmed := strings.TrimRight(p.rawContent, " \t\r\n")
	trailing := p.rawContent[len(trimmed):]

	if before, ok := strings.CutSuffix(trailing, "\n"); ok {
		trailing = before
		trailing = strings.TrimSuffix(trailing, "\r")
	}

	if strings.TrimRight(trailing, " \t") == "" {
		return
	}

	line := 0
	for _, l := range p.lines {
		if l.kind != kindBlank {
			line = l.lineNumber
		}
	}

	addError(r, "TLDR008", line+1)
}

// checkEndsWithNewline enforces TLDR009.
//
// It reports an error if the page does not end with a newline.
func checkEndsWithNewline(p *parsedPage, r *Result) {
	if !strings.HasSuffix(p.rawContent, "\n") {
		addError(r, "TLDR009", 0)
	}
}

// checkUnixLineEndings enforces TLDR010.
//
// It reports an error if the page contains carriage returns,
// that is, non-Unix (CRLF or CR) line endings.
func checkUnixLineEndings(p *parsedPage, r *Result) {
	if strings.Contains(p.rawContent, "\r") {
		addError(r, "TLDR010", 0)
	}
}

// checkConsecutiveBlankLines enforces TLDR011.
//
// It reports an error for every blank line that directly
// follows another blank line.
func checkConsecutiveBlankLines(p *parsedPage, r *Result) {
	count := 0
	for _, l := range p.lines {
		if l.kind == kindBlank {
			count++
			if count > 1 {
				addError(r, "TLDR011", l.lineNumber)
			}
		} else {
			count = 0
		}
	}
}

// checkNoTabs enforces TLDR012.
//
// It reports an error if the page contains any tab character.
func checkNoTabs(p *parsedPage, r *Result) {
	if strings.Contains(p.rawContent, "\t") {
		addError(r, "TLDR012", 0)
	}
}

// checkTrailingWhitespace enforces TLDR014.
//
// It reports an error for every line that ends with a space or a tab.
func checkTrailingWhitespace(p *parsedPage, r *Result) {
	for _, l := range p.lines {
		s := strings.TrimRight(l.rawLine, "\r")
		if len(s) > 0 &&
			(s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
			addError(r, "TLDR014", l.lineNumber)
		}
	}
}
