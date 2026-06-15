package termcolor

import "fmt"

// Sprint returns the given text
// wrapped with ANSI color codes
// based on the provided style string.
//
// The style string is parsed using Parse.
// If the resulting style produces no ANSI sequence,
// the original text is returned unchanged.
//
//	Sprint("bold red", "error") => "\x1b[1;31merror\x1b[0m"
func Sprint(styleString, text string) string {
	c := Parse(styleString)
	if c.String() == "" {
		return text
	}
	return c.String() + text + Reset()
}

// Fprintf formats the given arguments
// according to the format specifier
// and applies the provided style string
// to the resulting text.
//
// It is equivalent to calling fmt.Sprintf followed by Sprint.
//
//	Fprintf("green", "hello %s", "world") => "\x1b[32mhello world\x1b[0m"
func Fprintf(styleString, format string, args ...any) string {
	return Sprint(styleString, fmt.Sprintf(format, args...))
}
