package termcolor

import "strings"

// Parse parses a space-separated style string into a Color.
// It supports the standard ansi colors and effects.
//
// Only the first foreground and background are used
// and multiple effects are allowed.
//
// Unknown tokens are ignored,
// empty string or "default" returns a zero-value Color.
func Parse(s string) *Color {
	if s == "" || s == "default" {
		return &Color{}
	}

	c := &Color{}
	parts := strings.FieldsSeq(s)

	for part := range parts {
		if _, ok := foregroundCodes[part]; ok {
			if c.Foreground == "" {
				c.Foreground = part
			}
		} else if _, ok := backgroundCodes[part]; ok {
			if c.Background == "" {
				c.Background = part
			}
		} else if _, ok := effectCodes[part]; ok {
			c.Effects = append(c.Effects, part)
		}
	}

	return c
}
