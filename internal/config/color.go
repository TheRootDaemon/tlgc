package config

import (
	"fmt"
	"strconv"
	"strings"
)

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

func DefaultColor() OutputColor {
	return OutputColor{
		Kind:  ColorKindNamed,
		Named: ColorDefault,
	}
}

func (c OutputColor) MarshalText() ([]byte, error) {
	switch c.Kind {
	case ColorKindNamed:
		if c.Named == "" || c.Named == ColorDefault {
			return []byte("default"), nil
		}
		return []byte(c.Named), nil
	case ColorKindColor256:
		return fmt.Appendf(nil, "color256:%d", c.Color256), nil
	case ColorKindRGB:
		return fmt.Appendf(nil, "rgb:%d,%d,%d", c.RGB[0], c.RGB[1], c.RGB[2]), nil
	default:
		return []byte("default"), nil
	}
}

func (c *OutputColor) UnmarshalText(text []byte) error {
	s := string(text)

	// named ANSI colors
	switch s {
	case "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"bright-black", "bright-red", "bright-green", "bright-yellow",
		"bright-blue", "bright-magenta", "bright-cyan", "bright-white",
		"default", "grey":
		c.Kind = ColorKindNamed
		c.Named = ColorName(s)
		return nil
	}

	// hex colors
	if strings.HasPrefix(s, "#") {
		c.Kind = ColorKindRGB
		rgb, err := parseHexToRGB(s)
		if err != nil {
			return err
		}

		c.RGB = rgb
		return nil
	}

	// 256-color ANSI palette
	if strings.HasPrefix(s, "color256:") || strings.HasPrefix(s, "256:") {
		prefix := "color256:"
		if strings.HasPrefix(s, "256:") {
			prefix = "256:"
		}

		val, err := strconv.ParseUint(s[len(prefix):], 10, 8)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidColor256, err)
		}

		c.Kind = ColorKindColor256
		c.Color256 = uint8(val)
		return nil
	}

	// explicit RGB colors
	if strings.HasPrefix(s, "rgb:") {
		parts := strings.Split(s[4:], ",")
		if len(parts) != 3 {
			return fmt.Errorf("%w: %v", ErrInvalidRGB, s)
		}

		for i, p := range parts {
			val, err := strconv.ParseUint(strings.TrimSpace(p), 10, 8)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrInvalidRGB, err)
			}
			c.RGB[i] = uint8(val)
		}

		c.Kind = ColorKindRGB
		return nil
	}

	return fmt.Errorf("%w: %v", ErrUnknownColor, s)
}

func parseHexToRGB(hex string) ([3]uint8, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		hex = string(hex[0]) + string(hex[0]) +
			string(hex[1]) + string(hex[1]) +
			string(hex[2]) + string(hex[2])
	}

	if len(hex) != 6 {
		return [3]uint8{}, fmt.Errorf("%w: %v", ErrInvalidHex, hex)
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return [3]uint8{}, fmt.Errorf("%w: %v", ErrInvalidHex, hex)
	}

	return [3]uint8{r, g, b}, nil
}
