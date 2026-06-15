package pathutil

import (
	"path/filepath"
	"strings"
)

// PageName extracts the page name from a tldr page file path.
// For "/path/to/pages.en/common/some-page.md" it returns "some-page".
// Returns the empty string if the path has no file stem.
func PageName(path string) string {
	base := filepath.Base(path)
	if base == "." || base == "/" || base == "\\" {
		return ""
	}

	return strings.TrimSuffix(base, filepath.Ext(base))
}

// PagePlatform extracts the platform from a tldr page file path.
// For "/path/to/pages.en/linux/some-page.md" it returns "linux".
// Returns the empty string if the path has no parent directory with a name.
func PagePlatform(path string) string {
	parent := filepath.Base(filepath.Dir(path))

	if parent == "." || parent == "/" || parent == "\\" {
		return ""
	}

	return parent
}
