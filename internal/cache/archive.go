package cache

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
) (pageStats, error) {
	logger.Debug("extracting '%s'... ", languageDirectory)

	targetDirectory := filepath.Join(c.dir, languageDirectory)

	oldPages, err := hashPages(targetDirectory)
	if err != nil {
		return pageStats{}, err
	}

	if err := recreateDirectory(targetDirectory); err != nil {
		return pageStats{}, err
	}

	root, err := os.OpenRoot(targetDirectory)
	if err != nil {
		return pageStats{}, err
	}
	defer func() {
		_ = root.Close()
	}()

	return extractZip(root, data, oldPages)
}

// recreateDirectory removes dir if it exists
// and recreates it with 0o750 permissions.
func recreateDirectory(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, 0o750)
}

// extractZip extracts the zip data into root
// and classifies every extracted page against oldPages:
// added when previously absent,
// modified when its content hash differs.
// Pages in oldPages that the archive does not contain count as removed.
func extractZip(
	root *os.Root,
	data []byte,
	oldPages map[string]string,
) (pageStats, error) {
	zr, err := zip.NewReader(
		bytes.NewReader(data),
		int64(len(data)),
	)
	if err != nil {
		return pageStats{}, fmt.Errorf("reading zip archive: %w", err)
	}

	var st pageStats
	extracted := make(map[string]struct{}, len(oldPages))

	for _, f := range zr.File {
		isPage, hash, err := extractZipEntry(root, f)
		if err != nil {
			return pageStats{}, err
		}
		if !isPage {
			continue
		}

		name := filepath.Clean(f.Name)
		st.totalPages++
		extracted[name] = struct{}{}

		oldHash, ok := oldPages[name]
		switch {
		case !ok:
			st.added++
		case oldHash != hash:
			st.modified++
		}
	}

	for name := range oldPages {
		if _, ok := extracted[name]; !ok {
			st.removed++
		}
	}

	return st, nil
}

// extractZipEntry extracts a single zip entry into root.
func extractZipEntry(
	root *os.Root,
	f *zip.File,
) (bool, string, error) {
	if !filepath.IsLocal(f.Name) {
		logger.Debug(
			"skipping unsafe zip entry: %s",
			f.Name,
		)
		return false, "", nil
	}

	name := filepath.Clean(f.Name)

	if f.FileInfo().IsDir() {
		return false, "", root.MkdirAll(name, 0o750)
	}

	// pages always live under a platform subdir,
	// root-level entries like the other miscellaneous files gets skipped
	if !isPageFile(name) {
		logger.Debug(
			"skipping root-level zip entry: %s",
			f.Name,
		)
		return false, "", nil
	}

	if err := root.MkdirAll(filepath.Dir(name), 0o750); err != nil {
		return false, "", fmt.Errorf("creating directory for %s: %w", f.Name, err)
	}

	hash, err := extractFile(root, f)
	if err != nil {
		return false, "", err
	}

	return true, hash, nil
}

// extractFile writes a single zip entry to disk
// using the given root directory
// and returns the sha256 hash of its content.
func extractFile(
	root *os.Root,
	f *zip.File,
) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", fmt.Errorf("opening %s in zip: %w", f.Name, err)
	}

	hasher := sha256.New()

	out, err := root.OpenFile(
		f.Name,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		f.Mode(),
	)
	if err != nil {
		_ = rc.Close()
		return "", fmt.Errorf("creating %s: %w", f.Name, err)
	}

	_, err = out.ReadFrom(io.TeeReader(rc, hasher))
	_ = rc.Close()
	_ = out.Close()

	if err != nil {
		return "", fmt.Errorf("writing %s: %w", f.Name, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
