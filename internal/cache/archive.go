package cache

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/upstream"
	"github.com/TheRootDaemon/tlgc/logger"
)

// downloadArchive downloads the named archive from the given mirror
// and verifies it against the provided hash.
func downloadArchive(
	ctx context.Context,
	client *upstream.Client,
	mirror string,
	archiveName, hash string,
) ([]byte, error) {
	url := mirror + "/" + archiveName
	return client.DownloadBytes(ctx, url, hash)
}

// extractArchive removes the existing language directory,
// recreates it,
// and extracts the zip archive contents into it.
func (c *Cache) extractArchive(
	languageDirectory string,
	data []byte,
) error {
	logger.InfoStart("extracting '%s'... ", languageDirectory)

	targetDirectory := filepath.Join(c.dir, languageDirectory)

	if err := os.RemoveAll(targetDirectory); err != nil {
		return err
	}
	if err := os.MkdirAll(targetDirectory, 0o750); err != nil {
		return err
	}

	zipReader, err := zip.NewReader(
		bytes.NewReader(data),
		int64(len(data)),
	)
	if err != nil {
		return fmt.Errorf("reading zip archive: %w", err)
	}

	root, err := os.OpenRoot(targetDirectory)
	if err != nil {
		return err
	}
	defer func() {
		_ = root.Close()
	}()

	var extracted int
	for _, f := range zipReader.File {
		if strings.Contains(f.Name, "..") {
			continue
		}

		if f.FileInfo().IsDir() {
			if err := root.MkdirAll(f.Name, 0o750); err != nil {
				return fmt.Errorf("creating directory %s: %w", f.Name, err)
			}

			continue
		}

		if err := root.MkdirAll(filepath.Dir(f.Name), 0o750); err != nil {
			return fmt.Errorf("creating directory for %s: %w", f.Name, err)
		}

		if err := extractFile(root, f); err != nil {
			return err
		}

		extracted++
	}

	logger.InfoEnd("%d pages", extracted)
	return nil
}

// extractFile writes a single zip entry to disk
// using the given root directory.
func extractFile(
	root *os.Root,
	f *zip.File,
) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening %s in zip: %w", f.Name, err)
	}

	out, err := root.OpenFile(
		f.Name,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		f.Mode(),
	)
	if err != nil {
		_ = rc.Close()
		return fmt.Errorf("creating %s: %w", f.Name, err)
	}

	_, err = out.ReadFrom(rc)
	_ = rc.Close()
	_ = out.Close()

	if err != nil {
		return fmt.Errorf("writing %s: %w", f.Name, err)
	}

	return nil
}
