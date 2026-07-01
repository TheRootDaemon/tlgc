package config

import "github.com/TheRootDaemon/tlgc/termcolor"

// OutputStyle defines the visual style
// for a section of a rendered page.
//
// It combines a foreground color,
// a background color,
// and optional text effects (bold, italic, etc.).
type OutputStyle struct {
	// Color is the foreground colour used for this element.
	Color OutputColor `toml:"color"`

	// Background is the background colour used for this element.
	Background OutputColor `toml:"background"`

	// Bold controls whether the text is rendered in bold.
	Bold bool `toml:"bold"`

	// Underline controls whether the text is underlined.
	Underline bool `toml:"underline"`

	// Italic controls whether the text is rendered in italic.
	Italic bool `toml:"italic"`

	// Dim controls whether the text is rendered with reduced intensity.
	Dim bool `toml:"dim"`

	// Strikethrough controls whether the text has a strikethrough line.
	Strikethrough bool `toml:"strikethrough"`
}

// StyleConfig defines the visual style
// for each semantic section of a rendered page.
type StyleConfig struct {
	// Title is the style applied to the page title.
	Title OutputStyle `toml:"title"`

	// Description is the style applied to description lines.
	Description OutputStyle `toml:"description"`

	// Bullet is the style applied to bullet items (example descriptions).
	Bullet OutputStyle `toml:"bullet"`

	// Example is the style applied to command examples.
	Example OutputStyle `toml:"example"`

	// URL is the style applied to the "More information" URL line.
	URL OutputStyle `toml:"url"`

	// InlineCode is the style applied to inline code spans inside descriptions.
	InlineCode OutputStyle `toml:"inline_code"`

	// Placeholder is the style applied to user-supplied placeholder values.
	Placeholder OutputStyle `toml:"placeholder"`
}

// DefaultStyleConfig returns
// the default style settings.
//
// These are the defaults:
//   - title:        magenta + bold
//   - description:  magenta
//   - bullet:       green
//   - example:      cyan
//   - url:          red + italic
//   - inline_code:  yellow + italic
//   - placeholder:  red + italic
func DefaultStyleConfig() StyleConfig {
	defBg := DefaultColor()

	return StyleConfig{
		Title: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorMagenta},
			Background: defBg,
			Bold:       true,
		},
		Description: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorMagenta},
			Background: defBg,
		},
		Bullet: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorGreen},
			Background: defBg,
		},
		Example: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorCyan},
			Background: defBg,
		},
		URL: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorRed},
			Background: defBg,
			Italic:     true,
		},
		InlineCode: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorYellow},
			Background: defBg,
			Italic:     true,
		},
		Placeholder: OutputStyle{
			Color:      OutputColor{Kind: ColorKindNamed, Named: ColorRed},
			Background: defBg,
			Italic:     true,
		},
	}
}

// ToTermColor converts the OutputStyle
// to a termcolor.Color value
// that can be used for ANSI rendering.
//
// Named colors are mapped directly,
// extended colors (256/RGB) are passed
// via FGParams and BGParams.
func (o OutputStyle) ToTermColor() *termcolor.Color {
	c := &termcolor.Color{}

	switch o.Color.Kind {
	case ColorKindNamed:
		if o.Color.Named != ColorDefault {
			c.Foreground = string(o.Color.Named)
		}
	case ColorKindColor256:
		c.FGParams = []int{38, 5, int(o.Color.Color256)}
	case ColorKindRGB:
		c.FGParams = []int{38, 2, int(o.Color.RGB[0]), int(o.Color.RGB[1]), int(o.Color.RGB[2])}
	}

	switch o.Background.Kind {
	case ColorKindNamed:
		if o.Background.Named != ColorDefault && o.Background.Named != "" {
			c.Background = "on_" + string(o.Background.Named)
		}
	case ColorKindColor256:
		c.BGParams = []int{48, 5, int(o.Background.Color256)}
	case ColorKindRGB:
		c.BGParams = []int{48, 2, int(o.Background.RGB[0]), int(o.Background.RGB[1]), int(o.Background.RGB[2])}
	}

	if o.Bold {
		c.Effects = append(c.Effects, "bold")
	}
	if o.Italic {
		c.Effects = append(c.Effects, "italic")
	}
	if o.Underline {
		c.Effects = append(c.Effects, "underline")
	}
	if o.Dim {
		c.Effects = append(c.Effects, "dim")
	}
	if o.Strikethrough {
		c.Effects = append(c.Effects, "strikethrough")
	}

	return c
}
