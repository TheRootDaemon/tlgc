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
	root, err := os.OpenRoot(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer func() {
		_ = root.Close()
	}()

	pages := make(map[string]string)

	err = fs.WalkDir(root.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if !isPageFile(path) {
			return nil
		}

		content, err := fs.ReadFile(root.FS(), path)
		if err != nil {
			return err
		}

		sum := sha256.Sum256(content)
		pages[path] = hex.EncodeToString(sum[:])
		return nil
	})

	return pages, err
}
