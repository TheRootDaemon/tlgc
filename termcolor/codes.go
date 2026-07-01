package termcolor

var (
	// foregroundCodes maps foreground color names to ANSI SGR codes.
	foregroundCodes = map[string]int{
		"black":   30,
		"red":     31,
		"green":   32,
		"yellow":  33,
		"blue":    34,
		"magenta": 35,
		"cyan":    36,
		"white":   37,
		"grey":    90,
	}

	// backgroundCodes maps background color names to ANSI SGR codes.
	backgroundCodes = map[string]int{
		"on_black":   40,
		"on_red":     41,
		"on_green":   42,
		"on_yellow":  43,
		"on_blue":    44,
		"on_magenta": 45,
		"on_cyan":    46,
		"on_white":   47,
	}

	// effectCodes maps text effect names to ANSI SGR codes.
	effectCodes = map[string]int{
		"bold":          1,
		"dim":           2,
		"italic":        3,
		"underline":     4,
		"reverse":       7,
		"blink":         5,
		"hidden":        8,
		"strikethrough": 9,
	}
)

// Color represents an ANSI terminal color and text attributes.
type Color struct {
	// Foreground is the foreground color name.
	Foreground string

	// Background is the background color name.
	Background string

	// Effects lists the text effects to apply.
	Effects []string

	// FGParams contains the ANSI SGR parameters for the foreground color.
	FGParams []int

	// BGParams contains the ANSI SGR parameters for the background color.
	BGParams []int
}
