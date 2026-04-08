package termcolor

import (
	"os"

	"golang.org/x/term"
)

var (
	// getenv is a wrapper around os.Getenv,
	getenv = os.Getenv

	// stdoutFd returns the file descriptor for standard output.
	stdoutFd = func() uintptr { return os.Stdout.Fd() }

	// isTerminal checks whether a file descriptor is a terminal.
	isTerminal = func(fd uintptr) bool {
		return term.IsTerminal(int(fd))
	}
)

// SupportsColor reports whether the current environment
// likely supports ANSI color output.
//
// It returns false in the following cases:
//   - The NO_COLOR environment variable is set (per https://no-color.org).
//   - The TERM environment variable is set to "dumb".
//   - Stdout is not connected to a terminal.
//
// Otherwise, it returns true.
func SupportsColor() bool {
	if getenv("NO_COLOR") != "" {
		return false
	}

	if getenv("TERM") == "dumb" {
		return false
	}

	return isTerminal(stdoutFd())
}
