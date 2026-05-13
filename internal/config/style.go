package config

import "github.com/TheRootDaemon/tlgc/pkg/termcolor"

type OutputStyle struct {
	Color         OutputColor `toml:"color"`
	Background    OutputColor `toml:"background"`
	Bold          bool        `toml:"bold"`
	Underline     bool        `toml:"underline"`
	Italic        bool        `toml:"italic"`
	Dim           bool        `toml:"dim"`
	Strikethrough bool        `toml:"strikethrough"`
}

type StyleConfig struct {
	Title       OutputStyle `toml:"title"`
	Description OutputStyle `toml:"description"`
	Bullet      OutputStyle `toml:"bullet"`
	Example     OutputStyle `toml:"example"`
	URL         OutputStyle `toml:"url"`
	InlineCode  OutputStyle `toml:"inline_code"`
	Placeholder OutputStyle `toml:"placeholder"`
}

func DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		Title: OutputStyle{
			Color: OutputColor{Kind: ColorKindNamed, Named: ColorMagenta},
			Bold:  true,
		},
		Description: OutputStyle{
			Color: OutputColor{Kind: ColorKindNamed, Named: ColorMagenta},
		},
		Bullet: OutputStyle{
			Color: OutputColor{Kind: ColorKindNamed, Named: ColorGreen},
		},
		Example: OutputStyle{
			Color: OutputColor{Kind: ColorKindNamed, Named: ColorCyan},
		},
		URL: OutputStyle{
			Color:  OutputColor{Kind: ColorKindNamed, Named: ColorRed},
			Italic: true,
		},
		InlineCode: OutputStyle{
			Color:  OutputColor{Kind: ColorKindNamed, Named: ColorYellow},
			Italic: true,
		},
		Placeholder: OutputStyle{
			Color:  OutputColor{Kind: ColorKindNamed, Named: ColorRed},
			Italic: true,
		},
	}
}

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
		if o.Background.Named != ColorDefault {
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
