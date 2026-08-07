package lint

import "strings"

// checkCommandWhitespace enforces TLDR021.
//
// It reports an error if a command example begins
// or ends with whitespace.
//
// An escaped space (backslash followed by a space) is allowed.
func checkCommandWhitespace(p *parsedPage, r *Result) {
	for _, l := range p.lines {
		if l.kind != kindCommand {
			continue
		}

		// leading space (but allow escaped space \ ).
		if len(l.content) > 0 &&
			l.content[0] == ' ' &&
			!strings.HasPrefix(l.content, `\ `) {
			addError(r, "TLDR021", l.lineNumber)
			continue
		}

		// trailing space (but allow escaped space \ ).
		if len(l.content) > 0 &&
			l.content[len(l.content)-1] == ' ' &&
			!strings.HasSuffix(l.content, `\ `) {
			addError(r, "TLDR021", l.lineNumber)
		}
	}
}

// checkCommandDescriptionAnnotated enforces TLDR101.
//
// It reports an error for text that appears between the title
// and the first example description but is not annotated with a '> ' prefix.
func checkCommandDescriptionAnnotated(p *parsedPage, r *Result) {
	for i, l := range p.lines {
		if l.kind != kindText {
			continue
		}

		for j := i + 1; j < len(p.lines); j++ {
			if p.lines[j].kind == kindBlank || p.lines[j].kind == kindText {
				continue
			}
			if p.lines[j].kind == kindExampleDesc {
				addError(r, "TLDR101", l.lineNumber)
			}
			break
		}
	}
}

// checkExampleDescriptionAnnotated enforces TLDR102.
//
// It reports an error when:
//   - unannotated text appears after an example description
//   - unannotated text is followed by a command,
//     indicating that the text is likely an example description
//     missing the "- " prefix.
func checkExampleDescriptionAnnotated(p *parsedPage, r *Result) {
	for i, l := range p.lines {
		if l.kind != kindText {
			continue
		}

		if i > 0 && p.lines[i-1].kind == kindExampleDesc {
			addError(r, "TLDR102", l.lineNumber)
			continue
		}

		// text before any example description but followed by a command -> TLDR102
		for j := i + 1; j < len(p.lines); j++ {
			if p.lines[j].kind == kindBlank {
				continue
			}
			if p.lines[j].kind == kindCommand {
				addError(r, "TLDR102", l.lineNumber)
			}
			break
		}
	}
}

// checkCommandClosingBacktick enforces TLDR103.
//
// It reports an error if a command example is missing its closing backtick.
func checkCommandClosingBacktick(p *parsedPage, r *Result) {
	for _, l := range p.lines {
		if l.kind == kindCommand && !l.hasClosingBacktick {
			addError(r, "TLDR103", l.lineNumber)
		}
	}
}

// checkCommandNotEmpty enforces TLDR110.
//
// It reports an error if a command example is empty,
// that is, contains nothing between the backticks.
func checkCommandNotEmpty(p *parsedPage, r *Result) {
	for _, l := range p.lines {
		if l.kind == kindCommand && l.hasClosingBacktick && l.content == "" {
			addError(r, "TLDR110", l.lineNumber)
		}
	}
}
