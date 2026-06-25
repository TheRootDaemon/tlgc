package cache

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/require"
)

// createTestZip builds an in-memory ZIP from a path→content map.
// Entries whose path ends with "/" are treated as directory entries.
func createTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for path, content := range files {
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
