package upstream

import (
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"time"

	"github.com/TheRootDaemon/tlgc/format"
	"github.com/TheRootDaemon/tlgc/logger"
)

// DownloadFile downloads the content at url and writes it to destination.
//
// The download is subject to the client's configured limits and timeouts.
//
// If sha256hex is non-empty,
// the downloaded content must match the expected SHA256 checksum
// or an error is returned.
//
// If the download fails,
// the checksum does not match,
// or the size limit is exceeded,
// any partially written file is removed.
func (c *Client) DownloadFile(ctx context.Context, url, sha256hex, destination string) error {
	logger.Info("downloading from %s...", url)
	start := time.Now()

	resp, err := c.execute(ctx, url)
	if err != nil {
		logger.InfoEnd("failed: %s", err)
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// #nosec G304 TODO (Use os.OpenRoot with the cache directory, after the cache package is implemented)
	f, err := os.OpenFile(
		destination,
		os.O_CREATE|os.O_WRONLY|os.O_EXCL,
		0o600,
	)
	if err != nil {
		logger.InfoEnd("failed: %s", err)
		return err
	}

	cleanup := true
	defer func() {
		_ = f.Close()
		if cleanup {
			_ = os.Remove(destination)
		}
	}()

	n, err := c.transfer(
		f,
		resp.Body,
		sha256hex,
	)
	if err != nil {
		return err
	}

	cleanup = false

	logger.InfoEnd(
		"done (%s in %s)",
		format.BytesFmt(
			int64(n),
		),
		format.DurationFmt(time.Since(start)),
	)

	return nil
}

// transfer copies data from source to destination.
//
// If expectedSHA256 is non-empty, a SHA256 checksum is computed
// while copying and verified against the expected value after the copy completes.
//
// When a maximum body size is configured,
// transfer returns an error
// if the copied data exceeds the limit.
func (c *Client) transfer(
	destination io.Writer,
	source io.Reader,
	expectedSHA256 string,
) (int64, error) {
	writer := destination
	var hasher hash.Hash

	if expectedSHA256 != "" {
		hasher = sha256.New()
		writer = io.MultiWriter(destination, hasher)
	}

	var (
		n   int64
		err error
	)
	if c.maxBodySize > 0 {
		n, err = io.Copy(
			writer,
			io.LimitReader(source, c.maxBodySize+1),
		)
	} else {
		n, err = io.Copy(writer, source)
	}

	if err != nil {
		return 0, fmt.Errorf("copy: %w", err)
	}

	if c.maxBodySize > 0 && n > c.maxBodySize {
		return 0, fmt.Errorf("body %d bytes exceeds limit %d", n, c.maxBodySize)
	}

	if expectedSHA256 != "" && hasher != nil {
		if err := verifySHA256Hash(
			hasher, expectedSHA256,
		); err != nil {
			return 0, err
		}
	}

	return n, nil
}
