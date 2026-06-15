package termcolor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected *Color
	}{
		{
			input:    "",
			expected: &Color{},
		},
		{
			input:    "default",
			expected: &Color{},
		},
		{
			input: "cyan",
			expected: &Color{
				Foreground: "cyan",
			},
		},
		{
			input: "bold cyan",
			expected: &Color{
				Foreground: "cyan",
				Effects:    []string{"bold"},
			},
		},
		{
			input: "on_blue white",
			expected: &Color{
				Foreground: "white",
				Background: "on_blue",
			},
		},
		{
			input: "bold underline cyan",
			expected: &Color{
				Foreground: "cyan",
				Effects:    []string{"bold", "underline"},
			},
		},
		{
			input: "red",
			expected: &Color{
				Foreground: "red",
			},
		},
		{
			input: "green",
			expected: &Color{
				Foreground: "green",
			},
		},
		{
			input: "yellow",
			expected: &Color{
				Foreground: "yellow",
			},
		},
		{
			input: "blue",
			expected: &Color{
				Foreground: "blue",
			},
		},
		{
			input: "magenta",
			expected: &Color{
				Foreground: "magenta",
			},
		},
		{
			input: "white",
			expected: &Color{
				Foreground: "white",
			},
		},
		{
			input: "grey",
			expected: &Color{
				Foreground: "grey",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Parse(tt.input)

			require.Equal(t, tt.expected.Foreground, got.Foreground)
			require.Equal(t, tt.expected.Background, got.Background)
			require.Equal(t, tt.expected.Effects, got.Effects)
		})
	}
}
