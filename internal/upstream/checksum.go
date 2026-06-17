package upstream

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// verifySHA256 verifies data against a hex-encoded SHA256 hash.
// It accepts raw 64-char hex and the "<hash>  filename" format.
func verifySHA256(data []byte, expected string) error {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		return nil
	}

	sum := sha256.Sum256(data)
	got := fmt.Sprintf("%x", sum)

	if len(expected) == 64 {
		if strings.EqualFold(got, expected) {
			return nil
		}

		return fmt.Errorf("sha256 mismatch: expected %s, got %s", expected, got)
	}

	parts := strings.Fields(expected)
	if len(parts) >= 1 {
		hash := parts[0]
		if strings.EqualFold(got, hash) {
			return nil
		}
	}

	return fmt.Errorf("sha256 mismatch: expected %s, got %s", expected, got)
}

// ParseChecksumLine parses a line from a .sha256 file.
// Format: "<hash>  <filename>" (two spaces, as produced by sha256sum).
// Returns the hash and filename.
// Returns an error for empty lines.
func ParseCheckSum(line string) (hash, filename string, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", fmt.Errorf("empty checksum line")
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid checksum line: %s", line)
	}

	hash = parts[0]
	filename = strings.Join(parts[1:], " ")
	filename = strings.TrimLeft(filename, "* ")

	return hash, filename, nil
}
