package text

import "strings"

// Wrap wraps a string s to fit within maxLen columns.
// Continuation lines are prefixed with indent.
// Words longer than maxLen are not split (they overflow).
// Empty strings return s unchanged.
func Wrap(s string, maxLen int, indent string) string {
	if s == "" || len(s) <= maxLen {
		return s
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}

	var b strings.Builder

	line := words[0]
	for _, word := range words[1:] {
		if len(line)+1+len(word) <= maxLen {
			line += " " + word
			continue
		}

		b.WriteString(line)
		b.WriteByte('\n')

		line = indent + word
	}

	b.WriteString(line)

	return b.String()
}
