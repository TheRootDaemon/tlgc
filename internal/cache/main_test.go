package cache

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/require"
)

// contentHash returns the hex-encoded sha256 of s,
// matching the hashes produced by the cache package.
func contentHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// archiveServer returns a handler serving the checksum file
// and the given archives, keyed by archive name.
// Each checksum entry is "<hash>  <name>".
func archiveServer(t *testing.T, checksums map[string]string, archives map[string][]byte) http.Handler {
	t.Helper()

	var sumLines strings.Builder
	for name, hash := range checksums {
		fmt.Fprintf(&sumLines, "%s  %s\n", hash, name)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/" + checksumFile:
			_, _ = w.Write([]byte(sumLines.String()))
		default:
			data, ok := archives[strings.TrimPrefix(r.URL.Path, "/")]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, _ = w.Write(data)
		}
	})
}

// createTestZip builds an in-memory ZIP from a path→content map.
// Entries whose path ends with "/" are created as proper directory entries.
func createTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for path, content := range files {
		if strings.HasSuffix(path, "/") {
			_, err := zw.CreateHeader(&zip.FileHeader{
				Name:   path,
				Method: zip.Store,
			})
			require.NoError(t, err)
			continue
		}

		w, err := zw.Create(path)
		require.NoError(t, err)

		if _, err := w.Write([]byte(content)); err != nil {
			require.NoError(t, err)
		}
	}

	require.NoError(t, zw.Close())
	return buf.Bytes()
}

// createEmptyZip builds a valid empty ZIP archive.
func createEmptyZip(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

// setupConfig writes a temporary config file,
// sets TLGC_CONFIG,
// reinitializes the config singleton,
// and returns a cleanup function.
func setupConfig(t *testing.T, dir, mirror string) func() {
	t.Helper()

	cfgDir := t.TempDir()
	cfgPath := filepath.Join(cfgDir, "config.toml")
	content := fmt.Sprintf("[cache]\ndir = %q\nmirror = %q\n", dir, mirror)
	require.NoError(t, os.WriteFile(cfgPath, []byte(content), 0o644))

	config.ResetForTesting()
	t.Setenv("TLGC_CONFIG", cfgPath)
	require.NoError(t, config.Initialize())
	return config.ResetForTesting
}
