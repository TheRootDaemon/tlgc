package cache

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/fs"
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
) (int, int, error) {
	logger.Debug("extracting '%s'... ", languageDirectory)

	targetDirectory := filepath.Join(c.dir, languageDirectory)

	preExisting, err := existingFiles(targetDirectory)
	if err != nil {
		return 0, 0, err
	}

	if err := recreateDirectory(targetDirectory); err != nil {
		return 0, 0, err
	}

	root, err := os.OpenRoot(targetDirectory)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		_ = root.Close()
	}()

	return extractZip(root, data, preExisting)
}

// existingFiles walks dir and returns a set of relative file paths.
// If dir does not exist, it returns an empty set and no error.
func existingFiles(dir string) (map[string]struct{}, error) {
	files := make(map[string]struct{})
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		rel, relErr := filepath.Rel(dir, path)
		if relErr != nil {
			return nil
		}
		files[rel] = struct{}{}
		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return files, nil
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
// and returns the total number of files extracted
// and how many are new (not present in preExisting).
func extractZip(
	root *os.Root,
	data []byte,
	preExisting map[string]struct{},
) (int, int, error) {
	zr, err := zip.NewReader(
		bytes.NewReader(data),
		int64(len(data)),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("reading zip archive: %w", err)
	}

	var extracted, newCount int
	for _, f := range zr.File {
		extractedFile, newFile, err := extractZipEntry(root, f, preExisting)
		if err != nil {
			return 0, 0, err
		}

		if extractedFile {
			extracted++
		}

		if newFile {
			newCount++
		}
	}
	return extracted, newCount, nil
}

// extractZipEntry extracts a single zip entry into root.
// It returns whether the entry was a file (not a directory)
// and whether it is new (not in preExisting).
func extractZipEntry(
	root *os.Root,
	f *zip.File,
	preExisting map[string]struct{},
) (bool, bool, error) {
	if !filepath.IsLocal(f.Name) {
		logger.Debug(
			"skipping unsafe zip entry: %s",
			f.Name,
		)
		return false, false, nil
	}

	name := filepath.Clean(f.Name)

	if f.FileInfo().IsDir() {
		return false, false, root.MkdirAll(name, 0o750)
	}

	// pages always live under a platform subdir,
	// root-level entries like the other miscellaneous files gets skipped
	if filepath.Dir(name) == "." {
		logger.Debug(
			"skipping root-level zip entry: %s",
			f.Name,
		)
		return false, false, nil
	}

	if err := root.MkdirAll(filepath.Dir(name), 0o750); err != nil {
		return false, false, fmt.Errorf("creating directory for %s: %w", f.Name, err)
	}

	if err := extractFile(root, f); err != nil {
		return false, false, err
	}

	_, ok := preExisting[name]
	return true, !ok, nil
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
