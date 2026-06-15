package pathutil

import (
	"path/filepath"
	"strings"
)

// PageName extracts the page name from a tldr page file path.
// For "/path/to/pages.en/common/some-page.md" it returns "some-page".
// Returns the empty string if the path has no file stem.
func PageName(path string) string {
	if path == "" {
		return ""
	}

	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext != "" {
		return strings.TrimSuffix(base, ext)
	}

	return base
}

// PagePlatform extracts the platform from a tldr page file path.
// For "/path/to/pages.en/linux/some-page.md" it returns "linux".
// Returns the empty string if the path has no parent directory with a name.
func PagePlatform(path string) string {
	if path == "" {
		return ""
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "/" || dir == "" {
		return ""
	}

	base := filepath.Base(dir)
	return base
}
