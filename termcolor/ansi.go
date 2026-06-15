package termcolor

import (
	"fmt"
	"strings"
)

// formatCodes converts a slice of ANSI numeric codes
// into a semicolon-separated string
// suitable to use in an ANSI escape sequence.
//
// It returns an empty string if the input slice is empty.
//
//	[]int{1, 31} => "1;31"
func formatCodes(codes []int) string {
	if len(codes) == 0 {
		return ""
	}

	if len(codes) == 1 {
		return fmt.Sprintf("%d", codes[0])
	}

	var code strings.Builder
	fmt.Fprintf(&code, "%d", codes[0])
	for i := 1; i < len(codes); i++ {
		fmt.Fprintf(&code, ";%d", codes[i])
	}

	return code.String()
}

// String returns the ANSI escape sequence corresponding to the Color.
//
// It combines effects, foreground, and background
// into a single escape sequence of the form: "\x1B[<codes>m"
//
// where <codes> is a semicolon-separated list of numeric ANSI codes.
//
// The order of codes is:
//   - effects (in the order they appear)
//   - foreground color (FGParams if set, otherwise named foreground)
//   - background color (BGParams if set, otherwise named background)
//
// If the Color has no valid foreground, background, effects,
// FGParams, or BGParams, an empty string is returned.
func (c *Color) String() string {
	if c.Foreground == "" && c.Background == "" && len(c.Effects) == 0 && len(c.FGParams) == 0 && len(c.BGParams) == 0 {
		return ""
	}

	var codes []int
	for _, effect := range c.Effects {
		if code, ok := effectCodes[effect]; ok {
			codes = append(codes, code)
		}
	}

	if len(c.FGParams) > 0 {
		codes = append(codes, c.FGParams...)
	} else if fg, ok := foregroundCodes[c.Foreground]; ok {
		codes = append(codes, fg)
	}

	if len(c.BGParams) > 0 {
		codes = append(codes, c.BGParams...)
	} else if bg, ok := backgroundCodes[c.Background]; ok {
		codes = append(codes, bg)
	}

	if len(codes) == 0 {
		return ""
	}

	return fmt.Sprintf("\x1B[%sm", formatCodes(codes))
}

// Reset returns the ANSI escape sequence that resets all terminal formatting.
//
// This is equivalent to "\x1B[0m"
// and is typically used to clear any applied colors or text effects.
func Reset() string {
	return "\x1B[0m"
}
