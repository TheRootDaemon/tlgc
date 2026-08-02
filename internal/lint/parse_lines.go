package lint

import "strings"

// parseLines splits raw content at '\n'
// and classifies each line.
//
// It returns nil if the content is empty or whitespace-only.
func parseLines(raw string) []parsedLine {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, "\n")
	lines := make([]parsedLine, 0, len(parts))
	for lineNumber, part := range parts {
		lines = append(
			lines,
			parseLine(lineNumber, part),
		)
	}

	return lines
}

// parseLine classifies a single source line
// and extracts any relevant content from it.
//
// Command lines record whether a closing backtick was present.
func parseLine(lineNumber int, rawLine string) parsedLine {
	pl := parsedLine{
		lineNumber: lineNumber,
		rawLine:    rawLine,
	}

	// classification ignores a possible trailing '\r' (DOS line endings).
	s := strings.TrimRight(rawLine, "\r")

	if s == "" || strings.TrimSpace(s) == "" {
		pl.kind = kindBlank
		return pl
	}

	switch s[0] {
	case '`':
		pl.kind = kindCommand
		pl.content, pl.hasClosingBacktick = parseCommandContent(s)

	case '#':
		pl.kind = kindTitle
		pl.content = strings.TrimSpace(s[1:])

	case '>':
		pl.kind = kindDescription
		pl.content = stripMarker(s)

	case '-':
		pl.kind = kindExampleDesc
		pl.content = stripMarker(s)

	default:
		pl.kind = kindText
		pl.content = strings.TrimRight(s, " \t")
	}

	return pl
}

// parseCommandContent extracts the value between backticks.
//
// It returns the extracted command text
// and reports whether a closing backtick was found.
func parseCommandContent(s string) (string, bool) {
	if s == "`" {
		return "", false
	}

	if !strings.HasPrefix(s, "`") {
		return s, false
	}

	inner := s[1:]
	if before, _, ok := strings.Cut(inner, "`"); ok {
		return before, true
	}

	return inner, false
}

// stripMarker removes the leading marker character
// and any immediately following space
// from a description or example-description line.
func stripMarker(s string) string {
	if len(s) > 1 && s[1] == ' ' {
		return strings.TrimRight(s[2:], " \t")
	}

	return strings.TrimRight(s[1:], " \t")
}
