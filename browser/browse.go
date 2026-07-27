package browser

import "os/exec"

// browse starts the named program with the given arguments
// and waits for it to complete.
var browse = func(name string, args ...string) error {
	// #nosec G204
	// executable names are hardcoded by this package
	// and never originate from user input.
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}
