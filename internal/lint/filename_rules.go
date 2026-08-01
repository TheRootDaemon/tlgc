package lint

import (
	"path/filepath"
	"strings"
)

// checkFileExtension reports an error if filename
// does not use the appropriate extension (.md)
// required for TLDR pages.
func checkFileExtension(filename string, r *Result) {
	if filepath.Ext(filename) != ".md" {
		addError(r, "TLDR107", 0)
	}
}

// checkFilenameWhitespace reports an error if the filename
// contains spaces or tab characters.
func checkFilenameWhitespace(filename string, r *Result) {
	base := filepath.Base(filename)
	if strings.ContainsAny(base, " \t") {
		addError(r, "TLDR108", 0)
	}
}

// checkFilenameLowercase reports an error if the filename
// contains uppercase letters.
func checkFilenameLowercase(filename string, r *Result) {
	base := filepath.Base(filename)
	if base != strings.ToLower(base) {
		addError(r, "TLDR109", 0)
	}
}

// checkForbiddenFilenameCharacters reports an error if filename
// contains characters that are invalid on Windows filesystems.
func checkForbiddenFilenameCharacters(filename string, r *Result) {
	base := filepath.Base(filename)
	if strings.ContainsAny(base, `<>:"/\|?*`) {
		addError(r, "TLDR111", 0)
	}
}
