package platform

import "runtime"

// Default returns the default platform name based on runtime.GOOS.
func Default() string {
	return resolveGOOS(runtime.GOOS)
}

// resolveGOOS maps a GOOS value to the corresponding tldr platform.
func resolveGOOS(goos string) string {
	switch goos {
	case "linux":
		return "linux"
	case "darwin":
		return "osx"
	case "windows":
		return "windows"
	case "android":
		return "android"
	default:
		return "common"
	}
}

// Resolve resolves platform aliases to their canonical name.
func Resolve(platform string) string {
	if platform == "macos" {
		return "osx"
	}

	return platform
}

// All returns a list of all known platform names.
func All() []string {
	return []string{"common", "linux", "osx", "windows", "android"}
}
