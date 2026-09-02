package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// pageStats holds page counts of a single archive extraction.
type pageStats struct {
	// total is the number of pages present in the archive.
	totalPages int

	// added is the number of pages not present in the old snapshot.
	added int

	// removed is the number of snapshot pages absent from the archive.
	removed int

	// modified is the number of pages whose content changed.
	modified int
}

// isPageFile reports whether name refers to a page
// ie., .md file inside a subdirectory (platform directory).
func isPageFile(path string) bool {
	return filepath.Dir(path) != "." && strings.HasSuffix(path, ".md")
}

// hashPages walks directory and returns SHA-256 hashes
// of all cached pages,
// keyed by their paths relative to directory.
//
// Only page files are included (see isPageFile).
//
// A missing directory yields an empty snapshot with no error,
// so first-time extraction works.
func hashPages(directory string) (map[string]string, error) {
	pages := make(map[string]string)

	if err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		relativePath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}

		if !isPageFile(relativePath) {
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			return err
		}

		pages[relativePath] = hash
		return nil
	}); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return pages, nil
}

// hashFile reads the file at path
// and returns its SHA-256 hash encoded
// as a hexadecimal string.
func hashFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	checksum := sha256.Sum256(content)
	return hex.EncodeToString(checksum[:]), nil
}
