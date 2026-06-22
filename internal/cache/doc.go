// Package cache manages local tldr-pages downloads.
//
// The cache stores tldr-pages in directories named pages.<lang>
// under a root cache directory.
//
// Each language directory contains platform subdirectories
// (common, linux, osx, windows, android) with .md page files.
package cache

const (
	// checksumFile is the name of the checksum file in the cache directory.
	checksumFile = "tldr.sha256sums"

	// englishDirectory is the name of the English pages directory.
	englishDirectory = "pages.en"
)
