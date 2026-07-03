package cmd

import "strings"

// stringListValue implements flag.
// Value for a string slice.
// It supports both repeated flags (-L de -L pl)
// and comma-separated values (-L de,pl).
type stringListValue struct {
	values *[]string
}

// String returns the comma-separated representation of the value.
func (v *stringListValue) String() string {
	if v.values == nil {
		return ""
	}
	return strings.Join(*v.values, ",")
}

// Set appends one or more comma-separated values to the slice
func (v *stringListValue) Set(s string) error {
	for part := range strings.SplitSeq(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			*v.values = append(*v.values, part)
		}
	}
	return nil
}
