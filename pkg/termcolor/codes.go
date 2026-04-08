package termcolor

var (
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

type Color struct {
	Foreground string
	Background string
	Effects    []string
}
