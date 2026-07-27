package browser

import (
	"fmt"
	"os"
	"runtime"
)

// Open opens the given URL in the system browser.
// It detects the platform and chooses the appropriate command.
//
// On macOS and Windows it checks for remote/SSH sessions first.
// On Linux (including WSL) it checks for an active display server.
// Returns an error if no display server is available.
func Open(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return browse("open", url)
	case "windows":
		return browse("explorer.exe", url)
	case "linux":
		return openOnLinux(url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// openOnLinux handles browser opening on Linux and WSL.
func openOnLinux(url string) error {
	if runningOnWSL() {
		return openOnWSL(url)
	}

	if !hasDisplay() {
		return fmt.Errorf("no display server detected")
	}

	return browse("xdg-open", url)
}

// openOnWSL handles browser opening inside Windows Subsystem for Linux.
func openOnWSL(url string) error {
	return browse("explorer.exe", url)
}

// hasDisplay checks whether a display server is available
// by testing the DISPLAY and WAYLAND_DISPLAY environment variables.
func hasDisplay() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}
