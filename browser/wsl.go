package browser

import (
	"os"
	"strings"
)

// runningOnWSL reports whether the current process is running
// inside Windows Subsystem for Linux (WSL).
func runningOnWSL() bool {
	return isWSL(os.ReadFile)
}

// isWSL reports whether the system is running inside Windows Subsystem for Linux (WSL)
// by inspecting the Linux kernel version information.
//
// The readFile function is injected to allow the detection logic to be
// tested without reading from the real filesystem.
func isWSL(readFile func(string) ([]byte, error)) bool {
	data, err := readFile("/proc/version")
	if err != nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(string(data)),
		"microsoft",
	)
}
