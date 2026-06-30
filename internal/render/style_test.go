package render

import (
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestApplyStyle(t *testing.T) {
	t.Parallel()

	someStyle := config.OutputStyle{
		Bold:  true,
		Color: config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorRed},
	}

	tests := []struct {
		name     string
		useColor bool
		style    config.OutputStyle
		input    string
	}{
		{
			name:     "useColor false returns input unchanged with empty style",
			useColor: false,
			style:    config.OutputStyle{},
			input:    "hello world",
		},
		{
			name:     "useColor false returns input unchanged with styled style",
			useColor: false,
			style:    someStyle,
			input:    "tar cf archive.tar",
		},
		{
			name:     "useColor true returns styled output not equal to input",
			useColor: true,
			style:    someStyle,
			input:    "placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				useColor: tt.useColor,
			}
			got := r.applyStyle(tt.style, tt.input)

			if !tt.useColor {
				assert.Equal(t, tt.input, got)
			} else {
				assert.NotEqual(t, tt.input, got)
				assert.Contains(t, got, tt.input)
				assert.True(t, strings.Contains(got, "\x1b["))
				assert.True(t, strings.HasSuffix(got, "\x1b[0m"))
			}
		})
	}

	t.Run("useColor true with empty style returns input unchanged", func(t *testing.T) {
		r := &Renderer{useColor: true}
		input := "some text"
		got := r.applyStyle(config.OutputStyle{}, input)
		assert.Equal(t, input, got)
	})
}

func TestStyleForSegment(t *testing.T) {
	t.Parallel()

	exampleStyle := config.OutputStyle{
		Bold:  true,
		Color: config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorCyan},
	}

	placeholderStyle := config.OutputStyle{
		Italic: true,
		Color:  config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorRed},
	}

	r := &Renderer{
		style: config.StyleConfig{
			Example:     exampleStyle,
			Placeholder: placeholderStyle,
		},
	}

	tests := []struct {
		name    string
		segment *Segment
		want    config.OutputStyle
	}{
		{
			name:    "Text segment returns Example style",
			segment: &Segment{Kind: Text},
			want:    exampleStyle,
		},
		{
			name:    "Placeholder segment returns Placeholder style",
			segment: &Segment{Kind: Placeholder},
			want:    placeholderStyle,
		},
		{
			name:    "Option segment returns Placeholder style",
			segment: &Segment{Kind: Option, Short: "-a", Long: "--all"},
			want:    placeholderStyle,
		},
		{
			name:    "Unknown kind defaults to Example style",
			segment: &Segment{Kind: Kind(99)},
			want:    exampleStyle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.styleForSegment(tt.segment)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStyleString(t *testing.T) {
	t.Parallel()

	red := config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorRed}
	blue := config.OutputColor{Kind: config.ColorKindNamed, Named: config.ColorBlue}
	def := config.DefaultColor()

	tests := []struct {
		name  string
		style config.OutputStyle
		want  string
	}{
		{
			name:  "empty style",
			style: config.OutputStyle{},
			want:  "",
		},
		{
			name:  "bold only",
			style: config.OutputStyle{Bold: true},
			want:  "bold",
		},
		{
			name:  "dim only",
			style: config.OutputStyle{Dim: true},
			want:  "dim",
		},
		{
			name:  "italic only",
			style: config.OutputStyle{Italic: true},
			want:  "italic",
		},
		{
			name:  "underline only",
			style: config.OutputStyle{Underline: true},
			want:  "underline",
		},
		{
			name:  "strikethrough only",
			style: config.OutputStyle{Strikethrough: true},
			want:  "strikethrough",
		},
		{
			name: "all effects",
			style: config.OutputStyle{
				Bold:          true,
				Dim:           true,
				Italic:        true,
				Underline:     true,
				Strikethrough: true,
			},
			want: "bold dim italic underline strikethrough",
		},
		{
			name:  "foreground color only",
			style: config.OutputStyle{Color: red},
			want:  "red",
		},
		{
			name:  "background color only",
			style: config.OutputStyle{Background: blue},
			want:  "on_blue",
		},
		{
			name:  "foreground and background combined",
			style: config.OutputStyle{Color: red, Background: blue},
			want:  "red on_blue",
		},
		{
			name:  "effect with foreground and background",
			style: config.OutputStyle{Bold: true, Color: red, Background: blue},
			want:  "bold red on_blue",
		},
		{
			name:  "default color excluded",
			style: config.OutputStyle{Color: def},
			want:  "",
		},
		{
			name:  "default background excluded",
			style: config.OutputStyle{Background: def},
			want:  "",
		},
		{
			name:  "both default color and background excluded",
			style: config.OutputStyle{Color: def, Background: def},
			want:  "",
		},
		{
			name:  "default foreground but real background still included",
			style: config.OutputStyle{Color: def, Background: blue},
			want:  "on_blue",
		},
		{
			name:  "default background but real foreground still included",
			style: config.OutputStyle{Color: red, Background: def},
			want:  "red",
		},
		{
			name: "default color and effects — effects still included",
			style: config.OutputStyle{
				Color: def,
				Bold:  true,
			},
			want: "bold",
		},
		{
			name:  "empty named color excluded",
			style: config.OutputStyle{Color: config.OutputColor{Kind: config.ColorKindNamed, Named: ""}},
			want:  "",
		},
		{
			name:  "zero-value color kind excluded (Kind=0 but Named empty)",
			style: config.OutputStyle{Color: config.OutputColor{Kind: 0, Named: ""}},
			want:  "",
		},
		{
			name:  "256-color foreground excluded from styleString",
			style: config.OutputStyle{Color: config.OutputColor{Kind: config.ColorKindColor256, Color256: 208}},
			want:  "",
		},
		{
			name:  "256-color background excluded from styleString",
			style: config.OutputStyle{Background: config.OutputColor{Kind: config.ColorKindColor256, Color256: 42}},
			want:  "",
		},
		{
			name:  "RGB foreground excluded from styleString",
			style: config.OutputStyle{Color: config.OutputColor{Kind: config.ColorKindRGB, RGB: [3]uint8{255, 0, 0}}},
			want:  "",
		},
		{
			name:  "RGB background excluded from styleString",
			style: config.OutputStyle{Background: config.OutputColor{Kind: config.ColorKindRGB, RGB: [3]uint8{0, 255, 0}}},
			want:  "",
		},
		{
			name: "multiple effects with colors",
			style: config.OutputStyle{
				Bold:       true,
				Italic:     true,
				Underline:  true,
				Color:      red,
				Background: blue,
			},
			want: "bold italic underline red on_blue",
		},
		{
			name: "effects with only background",
			style: config.OutputStyle{
				Dim:           true,
				Strikethrough: true,
				Background:    blue,
			},
			want: "dim strikethrough on_blue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := styleString(tt.style)
			assert.Equal(t, tt.want, got)
		})
	}
}
