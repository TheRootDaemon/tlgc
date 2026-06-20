package upstream

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/TheRootDaemon/tlgc/format"
	"github.com/TheRootDaemon/tlgc/logger"
)

// DownloadBytes downloads the content at url into memory
// and returns it as a byte slice.
//
// The download is subject to the client's configured limits and timeouts.
//
// If sha256hex is non-empty, the downloaded content must match
// the expected SHA256 checksum or an error is returned.
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
		format.BytesFmt(
			int64(len(data)),
		),
		format.DurationFmt(time.Since(start)),
	)

	return data, nil
}
