package duration

import (
	"fmt"
	"time"
)

// DurationFmt formats a time.Duration into a human-readable string.
//
// Examples:
//
//	1s        -> "1s"
//	1m        -> "1min"
//	1m1s      -> "1min, 1s"
//	1h        -> "1h"
//	1h1m      -> "1h, 1min"
//	1d        -> "1d"
//	1d1h      -> "1d, 1h"
func DurationFmt(d time.Duration) string {
	total := uint64(d.Seconds())

	const (
		minute uint64 = 60
		hour          = 60 * minute
		day           = 24 * hour
	)

	days := total / day
	hours := (total % day) / hour
	minutes := (total % hour) / minute
	seconds := total % minute

	switch {
	case days > 0:
		if hours > 0 {
			return fmt.Sprintf("%dd, %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)

	case hours > 0:
		if minutes > 0 {
			return fmt.Sprintf("%dh, %dmin", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)

	case minutes > 0:
		if seconds > 0 {
			return fmt.Sprintf("%dmin, %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dmin", minutes)

	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
