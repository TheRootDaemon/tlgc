package version

import "runtime/debug"

// Version is set at build time via -ldflags.
// Falls back to "dev" for local development.
var Version = "dev"

// String returns the version string.
// Prefers debug.ReadBuildInfo (embedded by go install/pkg@version),
// then the Version ldflags override, then "dev".
func String() string {
	info, ok := debug.ReadBuildInfo()
	if ok &&
		info.Main.Version != "" &&
		info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return Version
}
