package config

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/termcolor"
	"github.com/stretchr/testify/assert"
)

func TestOutputStyleToTermColor_namedColors(t *testing.T) {
	magenta := OutputColor{Kind: ColorKindNamed, Named: ColorMagenta}
	green := OutputColor{Kind: ColorKindNamed, Named: ColorGreen}
	cyan := OutputColor{Kind: ColorKindNamed, Named: ColorCyan}
	red := OutputColor{Kind: ColorKindNamed, Named: ColorRed}
	yellow := OutputColor{Kind: ColorKindNamed, Named: ColorYellow}
	def := DefaultColor()

	tests := []struct {
		name  string
		style OutputStyle
		check func(t *testing.T, c *termcolor.Color)
	}{
		{
			name: "title - magenta bold",
			style: OutputStyle{
				Color: magenta,
				Bold:  true,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "magenta", c.Foreground)
				assert.Contains(t, c.Effects, "bold")
			},
		},
		{
			name: "description - magenta",
			style: OutputStyle{
				Color: magenta,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "magenta", c.Foreground)
				assert.Empty(t, c.Effects)
			},
		},
		{
			name: "bullet - green",
			style: OutputStyle{
				Color: green,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "green", c.Foreground)
			},
		},
		{
			name: "example - cyan",
			style: OutputStyle{
				Color: cyan,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "cyan", c.Foreground)
			},
		},
		{
			name: "url - red italic",
			style: OutputStyle{
				Color:  red,
				Italic: true,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "red", c.Foreground)
				assert.Contains(t, c.Effects, "italic")
			},
		},
		{
			name: "inline_code - yellow italic",
			style: OutputStyle{
				Color:  yellow,
				Italic: true,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "yellow", c.Foreground)
				assert.Contains(t, c.Effects, "italic")
			},
		},
		{
			name: "placeholder - red italic",
			style: OutputStyle{
				Color:  red,
				Italic: true,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "red", c.Foreground)
				assert.Contains(t, c.Effects, "italic")
			},
		},
		{
			name: "default color - no foreground set",
			style: OutputStyle{
				Color: def,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Empty(t, c.Foreground)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.style.ToTermColor()
			tt.check(t, c)
		})
	}
}

func TestOutputStyleToTermColor_background(t *testing.T) {
	tests := []struct {
		name  string
		style OutputStyle
		check func(t *testing.T, c *termcolor.Color)
	}{
		{
			name: "named background",
			style: OutputStyle{
				Background: OutputColor{Kind: ColorKindNamed, Named: ColorBlue},
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, "on_blue", c.Background)
			},
		},
		{
			name: "default background - empty",
			style: OutputStyle{
				Background: DefaultColor(),
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Empty(t, c.Background)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.style.ToTermColor()
			tt.check(t, c)
		})
	}
}

func TestOutputStyleToTermColor_extendedColors(t *testing.T) {
	tests := []struct {
		name  string
		style OutputStyle
		check func(t *testing.T, c *termcolor.Color)
	}{
		{
			name: "256-color foreground",
			style: OutputStyle{
				Color: OutputColor{Kind: ColorKindColor256, Color256: 208},
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, []int{38, 5, 208}, c.FGParams)
			},
		},
		{
			name: "RGB foreground",
			style: OutputStyle{
				Color: OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{255, 0, 0}},
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, []int{38, 2, 255, 0, 0}, c.FGParams)
			},
		},
		{
			name: "256-color background",
			style: OutputStyle{
				Background: OutputColor{Kind: ColorKindColor256, Color256: 42},
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, []int{48, 5, 42}, c.BGParams)
			},
		},
		{
			name: "RGB background",
			style: OutputStyle{
				Background: OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{0, 255, 0}},
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Equal(t, []int{48, 2, 0, 255, 0}, c.BGParams)
			},
		},
		{
			name: "all effects",
			style: OutputStyle{
				Bold:          true,
				Italic:        true,
				Underline:     true,
				Dim:           true,
				Strikethrough: true,
			},
			check: func(t *testing.T, c *termcolor.Color) {
				assert.Contains(t, c.Effects, "bold")
				assert.Contains(t, c.Effects, "italic")
				assert.Contains(t, c.Effects, "underline")
				assert.Contains(t, c.Effects, "dim")
				assert.Contains(t, c.Effects, "strikethrough")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.style.ToTermColor()
			tt.check(t, c)
		})
	}
}

func TestOutputStyleToTermColor_emptyStyle(t *testing.T) {
	style := OutputStyle{}
	c := style.ToTermColor()

	assert.Empty(t, c.Foreground)
	assert.Empty(t, c.Background)
	assert.Empty(t, c.Effects)
	assert.Empty(t, c.FGParams)
	assert.Empty(t, c.BGParams)
}
