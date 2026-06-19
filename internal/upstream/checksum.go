package upstream

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"strings"
)

// verifySHA256hex verifies data against a hex-encoded SHA256 hash.
// It accepts raw 64-char hex and the "<hash>  filename" format.
func verifySHA256hex(got, expected string) error {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		return nil
	}

	if len(expected) != 64 {
		parts := strings.Fields(expected)
		if len(parts) >= 1 {
			expected = parts[0]
		}
	}

	if len(expected) == 64 {
		if strings.EqualFold(got, expected) {
			return nil
		}

		return fmt.Errorf("sha256 mismatch: expected %s, got %s", expected, got)
	}

	if strings.EqualFold(got, expected) {
		return nil
	}

	return fmt.Errorf("sha256 mismatch: expected %s, got %s", expected, got)
}

func verifySHA256(data []byte, expected string) error {
	sum := sha256.Sum256(data)

	return verifySHA256hex(
		fmt.Sprintf("%x", sum),
		expected,
	)
}

func verifySHA256Hash(h hash.Hash, expected string) error {
	return verifySHA256hex(
		fmt.Sprintf("%x", h.Sum(nil)),
		expected,
	)
}

// ParseChecksum parses a line from a .sha256 file.
// Format: "<hash>  <filename>" (two spaces, as produced by sha256sum).
// Returns the hash and filename.
// Returns an error for empty lines.
func ParseChecksum(line string) (hash, filename string, err error) {
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
