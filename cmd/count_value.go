package cmd

import "fmt"

// countValue implements flag.Value for a counter.
//
// Each occurrence of the associated flag increments the counter by one,
// allowing flags such as:
//
//	--verbose
//	--verbose --verbose
//
// to be represented as increasing verbosity levels.
type countValue struct {
	count *uint8
}

// String returns the current counter value as a string.
func (v *countValue) String() string {
	if v.count == nil {
		return "0"
	}

	return fmt.Sprintf("%d", *v.count)
}

// IsBoolFlag returns true so flag.Parse treats --verbose as a boolean flag
// that does not consume the next argument.
func (v *countValue) IsBoolFlag() bool {
	return true
}

// Set increments the counter each time the flag is encountered.
func (v *countValue) Set(string) error {
	*v.count++
	return nil
}
