package config

import "errors"

var (
	// ErrInvalidHex is returned
	// when a hex color string cannot be parsed,
	// e.g. "#gggggg" or "#ff".
	ErrInvalidHex = errors.New("invalid hex")

	// ErrUnknownColor is returned
	// when a color string does not match
	// any known format.
	ErrUnknownColor = errors.New("unknown color")

	// ErrInvalidColor256 is returned
	// when a 256-color palette index
	// is not a valid number.
	ErrInvalidColor256 = errors.New("invalid 256 color")

	// ErrInvalidRGB is returned
	// when an RGB color string
	// has an invalid format.
	ErrInvalidRGB = errors.New("invalid rgb")
)
