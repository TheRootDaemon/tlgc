package termcolor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatCodes(t *testing.T) {
	tests := []struct {
		name     string
		codes    []int
		expected string
	}{
		{
			name:     "empty slice",
			codes:    []int{},
			expected: "",
		},
		{
			name:     "single code",
			codes:    []int{1},
			expected: "1",
		},
		{
			name:     "two codes",
			codes:    []int{1, 31},
			expected: "1;31",
		},
		{
			name:     "multiple codes",
			codes:    []int{1, 31, 44},
			expected: "1;31;44",
		},
		{
			name:     "single large code",
			codes:    []int{256},
			expected: "256",
		},
		{
			name:     "negative codes",
			codes:    []int{-1, -2},
			expected: "-1;-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatCodes(tt.codes)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestColorString(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected string
	}{
		{
			name:     "empty color",
			color:    Color{},
			expected: "",
		},
		{
			name:     "empty color with empty effects",
			color:    Color{Foreground: "", Background: "", Effects: []string{}},
			expected: "",
		},
		{
			name:     "foreground only",
			color:    Color{Foreground: "red"},
			expected: "\x1B[31m",
		},
		{
			name:     "foreground only - black",
			color:    Color{Foreground: "black"},
			expected: "\x1B[30m",
		},
		{
			name:     "foreground only - grey",
			color:    Color{Foreground: "grey"},
			expected: "\x1B[90m",
		},
		{
			name:     "background only",
			color:    Color{Background: "on_red"},
			expected: "\x1B[41m",
		},
		{
			name:     "background only - on_black",
			color:    Color{Background: "on_black"},
			expected: "\x1B[40m",
		},
		{
			name:     "single effect - bold",
			color:    Color{Effects: []string{"bold"}},
			expected: "\x1B[1m",
		},
		{
			name:     "single effect - underline",
			color:    Color{Effects: []string{"underline"}},
			expected: "\x1B[4m",
		},
		{
			name:     "single effect - strikethrough",
			color:    Color{Effects: []string{"strikethrough"}},
			expected: "\x1B[9m",
		},
		{
			name:     "multiple effects",
			color:    Color{Effects: []string{"bold", "underline", "reverse"}},
			expected: "\x1B[1;4;7m",
		},
		{
			name:     "all effects",
			color:    Color{Effects: []string{"bold", "dim", "italic", "underline", "reverse", "blink", "hidden", "strikethrough"}},
			expected: "\x1B[1;2;3;4;7;5;8;9m",
		},
		{
			name:     "foreground and background",
			color:    Color{Foreground: "green", Background: "on_blue"},
			expected: "\x1B[32;44m",
		},
		{
			name:     "effects and foreground",
			color:    Color{Effects: []string{"bold", "italic"}, Foreground: "cyan"},
			expected: "\x1B[1;3;36m",
		},
		{
			name:     "effects and background",
			color:    Color{Effects: []string{"underline"}, Background: "on_green"},
			expected: "\x1B[4;42m",
		},
		{
			name:     "effects foreground background",
			color:    Color{Effects: []string{"bold"}, Foreground: "red", Background: "on_white"},
			expected: "\x1B[1;31;47m",
		},
		{
			name:     "full combination",
			color:    Color{Effects: []string{"bold", "underline"}, Foreground: "yellow", Background: "on_magenta"},
			expected: "\x1B[1;4;33;45m",
		},
		{
			name:     "unknown foreground",
			color:    Color{Foreground: "unknown"},
			expected: "",
		},
		{
			name:     "unknown background",
			color:    Color{Background: "unknown"},
			expected: "",
		},
		{
			name:     "unknown effects only",
			color:    Color{Effects: []string{"unknown"}},
			expected: "",
		},
		{
			name:     "mixed valid and invalid effects",
			color:    Color{Effects: []string{"bold", "invalid", "underline"}},
			expected: "\x1B[1;4m",
		},
		{
			name:     "unknown foreground with valid background",
			color:    Color{Foreground: "unknown", Background: "on_red"},
			expected: "\x1B[41m",
		},
		{
			name:     "valid foreground with unknown background",
			color:    Color{Foreground: "blue", Background: "unknown"},
			expected: "\x1B[34m",
		},
		{
			name:     "unknown effects with valid foreground",
			color:    Color{Effects: []string{"invalid"}, Foreground: "green"},
			expected: "\x1B[32m",
		},
		{
			name:     "unknown effects with valid background",
			color:    Color{Effects: []string{"invalid"}, Background: "on_cyan"},
			expected: "\x1B[46m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.String()
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestReset(t *testing.T) {
	got := Reset()
	require.Equal(t, "\x1B[0m", got)
}
