package format

import "fmt"

// BytesFmt formats an int64 byte count into a human-readable string using
// binary (KiB, MiB, GiB, etc.) units.
//
// Examples:
//
//	0                     -> "0 B"
//	1024                  -> "1.0 KiB"
//	1099511627776         -> "1.0 TiB"
//	math.MaxInt64         -> "8.0 EiB"
func BytesFmt(bytes int64) string {
	const unit = 1024
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(1), 0
	units := []string{
		"B",
		"KiB",
		"MiB",
		"GiB",
		"TiB",
		"PiB",
		"EiB",
	}

	for bytes >= div*unit && exp < len(units)-1 {
		div *= unit
		exp++
	}

	return fmt.Sprintf(
		"%.1f %s",
		float64(bytes)/float64(div),
		units[exp],
	)
}
