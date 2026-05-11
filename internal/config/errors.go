package config

import "errors"

var (
	ErrInvalidHex      = errors.New("invalid hex")
	ErrUnknownColor    = errors.New("unknown color")
	ErrInvalidColor256 = errors.New("invalid 256 color")
	ErrInvalidRGB      = errors.New("invalid rgb")
)
