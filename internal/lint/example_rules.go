package lint

import (
	"regexp"
	"strings"
	"unicode"
)

// infinitiveTensePattern matches an example description
// that starts with a gerund ("...ing ")
// or a third-person present-tense verb ("...s ").
var infinitiveTensePattern = regexp.MustCompile(`(^[A-Za-z]{3,}ing )|(^[A-Za-z]+[^usy]s )`)

// checkExampleDescriptionEndsWithColon enforces TLDR005.
//
// It reports an error if an example description does not end with a colon.
func checkExampleDescriptionEndsWithColon(p *parsedPage, r *Result) {
	for _, section := range p.exampleSections {
		val := strings.TrimRight(section.description, " \t")
		if val == "" {
			continue
		}

		runes := []rune(val)
		if runes[len(runes)-1] != ':' {
			addError(r, "TLDR005", section.descriptionLineNumber)
		}
	}
}

// checkExampleDescriptionSurroundedByBlankLines enforces TLDR007.
//
// It reports an error if an example description is not separated
// from the surrounding content by blank lines.
func checkExampleDescriptionSurroundedByBlankLines(p *parsedPage, r *Result) {
	for _, section := range p.exampleSections {
		// find the position of the description in p.lines.
		descriptionIdx := lineIndex(p.lines, section.descriptionLineNumber)
		if descriptionIdx < 0 {
			continue
		}

		// check for blank line before description.
		if descriptionIdx == 0 || p.lines[descriptionIdx-1].kind != kindBlank {
			addError(r, "TLDR007", section.descriptionLineNumber)
		}

		if len(section.commands) > 0 {
			// find the position of the first command.
			firstCommandIdx := lineIndex(p.lines, section.commands[0].lineNumber)
			if firstCommandIdx >= 0 && firstCommandIdx > descriptionIdx {
				// the line before the first command must be blank.
				if p.lines[firstCommandIdx-1].kind != kindBlank {
					addError(r, "TLDR007", section.commands[0].lineNumber)
				}
			}
		}
	}
}

// checkExampleDescriptionStartsWithCapital enforces TLDR015.
//
// It reports an error if an example description starts with a lower-case letter.
// A leading '[' (a placeholder) is allowed.
func checkExampleDescriptionStartsWithCapital(p *parsedPage, r *Result) {
	for _, section := range p.exampleSections {
		val := strings.TrimSpace(section.description)
		if val == "" {
			continue
		}

		// allowed: uppercase letter, or '['.
		runes := []rune(val)
		if runes[0] == '[' {
			continue
		}

		if unicode.IsLetter(runes[0]) && unicode.IsLower(runes[0]) {
			addError(r, "TLDR015", section.descriptionLineNumber)
		}
	}
}

// checkMaximumExampleCount enforces TLDR019.
//
// It reports an error if the page contains more than 8 examples.
func checkMaximumExampleCount(p *parsedPage, r *Result) {
	if len(p.exampleSections) > 8 {
		addError(r, "TLDR019", 0)
	}
}

// checkInfinitiveTense enforces TLDR104.
//
// It reports an error if an example description uses the gerund or
// present tense instead of the infinitive tense.
func checkInfinitiveTense(p *parsedPage, r *Result) {
	for _, section := range p.exampleSections {
		if infinitiveTensePattern.MatchString(section.description) {
			addError(r, "TLDR104", section.descriptionLineNumber)
		}
	}
}

// checkSingleCommandPerExample enforces TLDR105.
//
// It reports an error for every command in an example that has more than one command.
func checkSingleCommandPerExample(p *parsedPage, r *Result) {
	for _, section := range p.exampleSections {
		if len(section.commands) > 1 {
			for _, cmd := range section.commands {
				addError(r, "TLDR105", cmd.lineNumber)
			}
		}
	}
}

// lineIndex returns the index of the line
// with the given line number.
// It returns -1 if no such line exists.
func lineIndex(lines []parsedLine, lineNumber int) int {
	for i, l := range lines {
		if l.lineNumber == lineNumber {
			return i
		}
	}

	return -1
}
