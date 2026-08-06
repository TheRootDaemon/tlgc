package lint

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// descriptionCapitalExceptions lists lowercase words that are allowed
	// to start a description despite the TLDR003 capitalization rule.
	descriptionCapitalExceptions = map[string]bool{
		"npm":  true,
		"pnpm": true,
	}

	// infoLinkLabelPattern matches the "More information:" label case-insensitively.
	infoLinkLabelPattern = regexp.MustCompile(`(?i)^(more\s+info(?:rmation)?:?\s*)`)

	// noteLabelPattern matches an incorrectly formatted "Note:" label.
	noteLabelPattern = regexp.MustCompile(`\b(note|NOTE): `)

	// urlPattern matches a URL wrapped in angle brackets.
	urlPattern = regexp.MustCompile(`<[^>]+>`)

	// standardTermPattern matches the standard terms
	// that must be wrapped in backticks.
	// Word boundaries are checked in checkValueForStandardTerms
	// because RE2 does not support lookaround assertions.
	standardTermPattern = regexp.MustCompile(`(?i)(stdout|stdin|stderr|regex|regular\s+expression|standard\s+(?:input|in|output|out|error|err))`)
)

// checkDescriptionStartsWithCapital enforces TLDR003.
//
// It reports an error if a description starts with a lower-case
// letter and the first word is not an allowed exception.
func checkDescriptionStartsWithCapital(p *parsedPage, r *Result) {
	for _, d := range p.descriptions {
		val := strings.TrimSpace(d.content)
		if val == "" {
			continue
		}

		// extract first word.
		firstWord := val
		if idx := strings.IndexAny(val, " \t"); idx >= 0 {
			firstWord = val[:idx]
		}
		if descriptionCapitalExceptions[firstWord] {
			continue
		}

		// check if first rune is a lower-case letter.
		runes := []rune(val)
		if len(runes) > 0 &&
			unicode.IsLetter(runes[0]) &&
			unicode.IsLower(runes[0]) {
			addError(r, "TLDR003", d.lineNumber)
		}
	}
}

// checkDescriptionEndsWithPeriod enforces TLDR004.
//
// It reports an error if a description does not end with a period.
// Information link lines are exempt.
func checkDescriptionEndsWithPeriod(p *parsedPage, r *Result) {
	for _, d := range p.descriptions {
		if isInfoLink(d.content) {
			continue
		}
		val := strings.TrimRight(d.content, " \t")
		if val == "" {
			continue
		}

		runes := []rune(val)
		if runes[len(runes)-1] != '.' {
			addError(r, "TLDR004", d.lineNumber)
		}
	}
}

// checkInformationLinkLabel enforces TLDR016.
//
// It reports an error if an information link uses any label
// other than the exact string "More information: ".
func checkInformationLinkLabel(p *parsedPage, r *Result) {
	for _, d := range p.descriptions {
		m := infoLinkLabelPattern.FindStringSubmatch(d.content)
		if m != nil && m[1] != "More information: " {
			addError(r, "TLDR016", d.lineNumber)
		}
	}
}

// checkInformationLinkBrackets enforces TLDR017.
//
// It reports an error if an information link URL is not surrounded by angle brackets.
func checkInformationLinkBrackets(p *parsedPage, r *Result) {
	for _, l := range p.infoLinks {
		val := l.content
		if !strings.Contains(val, "<") || !strings.Contains(val, ">") {
			addError(r, "TLDR017", l.lineNumber)
		}
	}
}

// checkSingleInformationLink enforces TLDR018.
//
// It reports an error for every information link after the first.
func checkSingleInformationLink(p *parsedPage, r *Result) {
	if len(p.infoLinks) > 1 {
		for _, l := range p.infoLinks[1:] {
			addError(r, "TLDR018", l.lineNumber)
		}
	}
}

// checkNoteLabelFormat enforces TLDR020.
//
// It reports an error if a description or example description
// contains a "Note:" label that is not exactly "Note: ".
func checkNoteLabelFormat(p *parsedPage, r *Result) {
	for _, d := range p.descriptions {
		if noteLabelPattern.MatchString(d.content) {
			addError(r, "TLDR020", d.lineNumber)
		}
	}
	for _, sec := range p.exampleSections {
		if noteLabelPattern.MatchString(sec.description) {
			addError(r, "TLDR020", sec.descriptionLineNumber)
		}
	}
}

// checkStandardTermsInBackticks enforces TLDR112.
//
// It reports an error when a standard term
// (stdin, stdout, stderr, regex)
// appears outside backticks in a description or example description.
func checkStandardTermsInBackticks(p *parsedPage, r *Result) {
	for _, d := range p.descriptions {
		checkValueForStandardTerms(
			d.content,
			d.lineNumber,
			r,
		)
	}
	for _, sec := range p.exampleSections {
		checkValueForStandardTerms(
			sec.description,
			sec.descriptionLineNumber,
			r,
		)
	}
}

// checkValueForStandardTerms checks one text value
// for standard terms that are not wrapped in backticks,
// reporting at most one error per line.
func checkValueForStandardTerms(val string, line int, r *Result) {
	// split on backticks:
	// even parts are outside backticks,
	// odd parts are inside.
	for i, part := range strings.Split(val, "`") {
		if i%2 == 1 {
			continue
		}

		part = urlPattern.ReplaceAllString(part, "")

		for _, m := range standardTermPattern.FindAllStringIndex(part, -1) {
			if m[0] > 0 &&
				(part[m[0]-1] == '-' || isWordCharacter(part[m[0]-1])) {
				continue
			}

			if m[1] < len(part) &&
				(part[m[1]] == '-' || isWordCharacter(part[m[1]])) {
				continue
			}

			addError(r, "TLDR112", line)
			return // at most one error per line
		}
	}
}

// isWordCharacter reports whether b is an ASCII letter, digit, or underscore.
func isWordCharacter(b byte) bool {
	return b == '_' ||
		(b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9')
}
