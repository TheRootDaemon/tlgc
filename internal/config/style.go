package config

type ColorKind int

const (
	ColorKindNamed ColorKind = iota
	ColorKindColor256
	ColorKindRGB
)

type ColorName string

const (
	ColorBlack         ColorName = "black"
	ColorRed           ColorName = "red"
	ColorGreen         ColorName = "green"
	ColorYellow        ColorName = "yellow"
	ColorBlue          ColorName = "blue"
	ColorMagenta       ColorName = "magenta"
	ColorCyan          ColorName = "cyan"
	ColorWhite         ColorName = "white"
	ColorBrightBlack   ColorName = "bright-black"
	ColorBrightRed     ColorName = "bright-red"
	ColorBrightGreen   ColorName = "bright-green"
	ColorBrightYellow  ColorName = "bright-yellow"
	ColorBrightBlue    ColorName = "bright-blue"
	ColorBrightMagenta ColorName = "bright-magenta"
	ColorBrightCyan    ColorName = "bright-cyan"
	ColorBrightWhite   ColorName = "bright-white"
	ColorDefault       ColorName = "default"
)

type OutputColor struct {
	Kind     ColorKind
	Named    ColorName
	Color256 uint8
	RGB      [3]uint8
}

type OutputStyle struct {
	Color         OutputColor
	Background    OutputColor
	Bold          bool
	Underline     bool
	Italic        bool
	Dim           bool
	Strikethrough bool
}

type StyleConfig struct {
	Title       OutputStyle
	Description OutputStyle
	Bullet      OutputStyle
	Example     OutputStyle
	URL         OutputStyle
	InlineCode  OutputStyle
	Placeholder OutputStyle
}
