package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultColor(t *testing.T) {
	c := DefaultColor()
	assert.Equal(t, ColorKindNamed, c.Kind)
	assert.Equal(t, ColorDefault, c.Named)
}

func TestOutputColorMarshalText(t *testing.T) {
	tests := []struct {
		name    string
		color   OutputColor
		want    string
		wantErr bool
	}{
		{
			name:  "named color",
			color: OutputColor{Kind: ColorKindNamed, Named: ColorMagenta},
			want:  "magenta",
		},
		{
			name:  "bright named color",
			color: OutputColor{Kind: ColorKindNamed, Named: ColorBrightRed},
			want:  "bright-red",
		},
		{
			name:  "default color",
			color: OutputColor{Kind: ColorKindNamed, Named: ColorDefault},
			want:  "default",
		},
		{
			name:  "empty named marshals as default",
			color: OutputColor{Kind: ColorKindNamed},
			want:  "default",
		},
		{
			name:  "256-color",
			color: OutputColor{Kind: ColorKindColor256, Color256: 208},
			want:  "color256:208",
		},
		{
			name:  "RGB color",
			color: OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{255, 0, 0}},
			want:  "rgb:255,0,0",
		},
		{
			name:  "default kind falls back to default",
			color: OutputColor{Kind: 99},
			want:  "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.color.MarshalText()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestOutputColorUnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    OutputColor
		wantErr bool
	}{
		{
			name:  "named color - magenta",
			input: "magenta",
			want:  OutputColor{Kind: ColorKindNamed, Named: ColorMagenta},
		},
		{
			name:  "named color - red",
			input: "red",
			want:  OutputColor{Kind: ColorKindNamed, Named: ColorRed},
		},
		{
			name:  "named color - bright cyan",
			input: "bright-cyan",
			want:  OutputColor{Kind: ColorKindNamed, Named: ColorBrightCyan},
		},
		{
			name:  "named color - default",
			input: "default",
			want:  OutputColor{Kind: ColorKindNamed, Named: ColorDefault},
		},
		{
			name:  "named color - grey",
			input: "grey",
			want:  OutputColor{Kind: ColorKindNamed, Named: ColorName("grey")},
		},
		{
			name:  "hex color - 6 digit",
			input: "#ff0000",
			want:  OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{255, 0, 0}},
		},
		{
			name:  "hex color - 3 digit",
			input: "#f00",
			want:  OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{255, 0, 0}},
		},
		{
			name:  "hex color - green",
			input: "#00ff00",
			want:  OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{0, 255, 0}},
		},
		{
			name:  "256 color - color256 prefix",
			input: "color256:208",
			want:  OutputColor{Kind: ColorKindColor256, Color256: 208},
		},
		{
			name:  "256 color - short prefix",
			input: "256:42",
			want:  OutputColor{Kind: ColorKindColor256, Color256: 42},
		},
		{
			name:  "256 color - zero",
			input: "color256:0",
			want:  OutputColor{Kind: ColorKindColor256, Color256: 0},
		},
		{
			name:  "256 color - max",
			input: "color256:255",
			want:  OutputColor{Kind: ColorKindColor256, Color256: 255},
		},
		{
			name:  "RGB color",
			input: "rgb:255,0,0",
			want:  OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{255, 0, 0}},
		},
		{
			name:  "RGB color with spaces",
			input: "rgb:0, 255, 128",
			want:  OutputColor{Kind: ColorKindRGB, RGB: [3]uint8{0, 255, 128}},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "unknown color name",
			input:   "unknown",
			wantErr: true,
		},
		{
			name:    "invalid hex - too short",
			input:   "#ff",
			wantErr: true,
		},
		{
			name:    "invalid hex - bad chars",
			input:   "#gggggg",
			wantErr: true,
		},
		{
			name:    "invalid 256 - out of range",
			input:   "color256:256",
			wantErr: true,
		},
		{
			name:    "invalid 256 - not a number",
			input:   "color256:abc",
			wantErr: true,
		},
		{
			name:    "invalid RGB - wrong parts count",
			input:   "rgb:255,0",
			wantErr: true,
		},
		{
			name:    "invalid RGB - not a number",
			input:   "rgb:255,abc,0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c OutputColor
			err := c.UnmarshalText([]byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, c)
		})
	}
}

func TestOutputColorRoundTrip(t *testing.T) {
	colors := []OutputColor{
		{Kind: ColorKindNamed, Named: ColorRed},
		{Kind: ColorKindNamed, Named: ColorMagenta},
		{Kind: ColorKindNamed, Named: ColorBrightGreen},
		{Kind: ColorKindNamed, Named: ColorDefault},
		{Kind: ColorKindColor256, Color256: 208},
		{Kind: ColorKindColor256, Color256: 0},
		{Kind: ColorKindColor256, Color256: 255},
		{Kind: ColorKindRGB, RGB: [3]uint8{255, 0, 0}},
		{Kind: ColorKindRGB, RGB: [3]uint8{0, 255, 0}},
		{Kind: ColorKindRGB, RGB: [3]uint8{128, 128, 128}},
	}

	for _, original := range colors {
		text, err := original.MarshalText()
		require.NoError(t, err)

		t.Run(string(text), func(t *testing.T) {
			var decoded OutputColor
			err = decoded.UnmarshalText(text)
			require.NoError(t, err)

			assert.Equal(t, original, decoded)
		})
	}
}

func TestParseHexToRGB(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [3]uint8
		wantErr bool
	}{
		{
			name:  "6-digit hex with hash",
			input: "#ff0000",
			want:  [3]uint8{255, 0, 0},
		},
		{
			name:  "6-digit hex without hash",
			input: "00ff00",
			want:  [3]uint8{0, 255, 0},
		},
		{
			name:  "3-digit hex with hash",
			input: "#f00",
			want:  [3]uint8{255, 0, 0},
		},
		{
			name:  "3-digit hex without hash",
			input: "0f0",
			want:  [3]uint8{0, 255, 0},
		},
		{
			name:  "white",
			input: "#ffffff",
			want:  [3]uint8{255, 255, 255},
		},
		{
			name:  "black",
			input: "#000000",
			want:  [3]uint8{0, 0, 0},
		},
		{
			name:    "too short",
			input:   "#ff",
			wantErr: true,
		},
		{
			name:    "bad characters",
			input:   "#gggggg",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHexToRGB(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
