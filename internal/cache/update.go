package cache

import (
	"context"
	"fmt"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/internal/upstream"
	"github.com/TheRootDaemon/tlgc/logger"
)

// Update downloads the latest tldr-pages archives
// for the given languages
// and extracts them into the cache.
func (c *Cache) Update(
	ctx context.Context,
	languages []string,
	client *upstream.Client,
) error {
	checksums, err := downloadChecksum(ctx, client)
	if err != nil {
		return fmt.Errorf("downloading checksum: %s", err)
	}

	oldChecksums := c.loadChecksums()
	newChecksums := parseChecksum(checksums)

	var downloaded int
	for _, language := range languages {
		updated, err := c.updateLanguage(
			ctx,
			client,
			language,
			oldChecksums,
			newChecksums,
		)
		if err != nil {
			return err
		}

		if updated {
			downloaded++
		}
	}

	if err := c.saveChecksums(newChecksums); err != nil {
		return fmt.Errorf("saving checksums: %w", err)
	}

	if downloaded == 0 {
		logger.Info("pages are up to date")
		return nil
	}

	c.platforms.Store([]string(nil))
	return nil
}

// updateLanguage downloads and extracts a single language archive,
// if the checksum has changed or the directory is missing.
// It returns true when an archive was actually downloaded.
func (c *Cache) updateLanguage(
	ctx context.Context,
	client *upstream.Client,
	language string,
	oldChecksums,
	newChecksums map[string]string,
) (bool, error) {
	languageDirectory := "pages." + language
	archiveName := fmt.Sprintf("tldr-pages.%s.zip", language)
	if !needsUpdate(
		c.subDirExists(languageDirectory),
		archiveName,
		oldChecksums,
		newChecksums,
	) {
		return false, nil
	}

	hash := newChecksums[archiveName]
	data, err := downloadArchive(
		ctx,
		client,
		config.Cache().Mirror,
		archiveName,
		hash,
	)
	if err != nil {
		return false, fmt.Errorf("downloading %s: %w", archiveName, err)
	}

	if err := c.extractArchive(languageDirectory, data); err != nil {
		return false, fmt.Errorf("extracting %s: %w", languageDirectory, err)
	}

	return true, nil
}

// needsUpdate reports whether an archive should be downloaded,
// based on whether it exists in the new checksums,
// whether we already have the same hash,
// and whether the language directory exists on disk.
func needsUpdate(
	exists bool,
	archive string,
	oldChecksums map[string]string,
	newChecksums map[string]string,
) bool {
	newHash, ok := newChecksums[archive]
	if !ok {
		return false
	}

	oldHash, hadOld := oldChecksums[archive]

	return !hadOld ||
		!exists ||
		oldHash != newHash
}
