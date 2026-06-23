package cache

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/internal/upstream"
)

// loadChecksums reads the cached checksum file from disk
// and returns it as a filename-to-hash map.
func (c *Cache) loadChecksums() map[string]string {
	root, err := os.OpenRoot(c.dir)
	if err != nil {
		return nil
	}
	defer func() {
		_ = root.Close()
	}()

	checksumBytes, err := root.ReadFile(checksumFile)
	if err != nil {
		return nil
	}

	return parseChecksum(checksumBytes)
}

// saveChecksums writes the checksum map to disk in sha256sum format.
func (c *Cache) saveChecksums(checksums map[string]string) error {
	if err := os.MkdirAll(c.dir, 0o750); err != nil {
		return fmt.Errorf("creating cache directory: %s", err)
	}

	var sb strings.Builder
	for name, hash := range checksums {
		fmt.Fprintf(&sb, "%s  %s\n", hash, name)
	}

	root, err := os.OpenRoot(c.dir)
	if err != nil {
		return err
	}
	defer func() {
		_ = root.Close()
	}()

	if err := root.WriteFile(
		checksumFile,
		[]byte(
			sb.String(),
		),
		0o600,
	); err != nil {
		return fmt.Errorf("writing checksums: %w", err)
	}

	return nil
}

// downloadChecksum fetches the tldr-pages checksum file from the configured mirror.
func downloadChecksum(
	ctx context.Context,
	client *upstream.Client,
) ([]byte, error) {
	mirror := config.Cache().Mirror
	checksumURL := mirror + "/" + checksumFile
	return client.DownloadBytes(ctx, checksumURL, "")
}

// parseChecksum parses sha256sum-formatted data into a filename-to-hash map.
func parseChecksum(checksum []byte) map[string]string {
	structuredChecksum := make(map[string]string)
	for line := range strings.SplitSeq(string(checksum), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		hash, filename, err := upstream.ParseChecksum(line)
		if err != nil {
			continue
		}

		structuredChecksum[filename] = hash
	}

	return structuredChecksum
}
