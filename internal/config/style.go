package config

import "github.com/TheRootDaemon/tlgc/pkg/termcolor"

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
		if o.Color.Named != ColorDefault {
			c.Background = string(o.Color.Named)
		}
	case ColorKindColor256:
		c.BGParams = []int{38, 5, int(o.Color.Color256)}
	case ColorKindRGB:
		c.BGParams = []int{38, 2, int(o.Color.RGB[0]), int(o.Color.RGB[1]), int(o.Color.RGB[2])}
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
