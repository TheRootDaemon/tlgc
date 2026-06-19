package upstream

import (
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"time"

	"github.com/TheRootDaemon/tlgc/duration"
	"github.com/TheRootDaemon/tlgc/logger"
)

func (c *Client) DownloadBytes(ctx context.Context, url, sha256hex string) ([]byte, error) {
	logger.Info("downloading from %s...", url)
	start := time.Now()

	resp, err := c.execute(ctx, url)
	if err != nil {
		logger.InfoEnd("failed: %s", err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %s", err)
	}

	if c.maxBodySize > 0 && int64(len(data)) > c.maxBodySize {
		return nil, fmt.Errorf("body %d exceeds limit %d", len(data), c.maxBodySize)
	}

	if sha256hex != "" {
		if err := verifySHA256(data, sha256hex); err != nil {
			return nil, err
		}
	}

	logger.InfoEnd(
		"done (%s in %s)",
		humanizeBytes(
			int64(len(data)),
		),
		duration.DurationFmt(time.Since(start)),
	)

	return data, nil
}

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
		humanizeBytes(
			int64(n),
		),
		duration.DurationFmt(time.Since(start)),
	)

	return nil
}

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

func humanizeBytes(bytes int64) string {
	const unit = 1024
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(1), 0
	units := []string{
		"B",
		"KiB",
		"MiB",
		"GiB",
		"TiB",
	}

	for bytes >= div*unit && exp < len(units) {
		div *= unit
		exp++
	}

	return fmt.Sprintf(
		"%.1f %s",
		float64(bytes)/float64(div),
		units[exp],
	)
}
