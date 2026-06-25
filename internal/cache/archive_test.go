package cache

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TheRootDaemon/tlgc/internal/upstream"
)

func TestDownloadArchive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		archiveName string
		hash        string
		handler     http.HandlerFunc
		wantErr     bool
		wantData    string
	}{
		{
			name:        "successful_download",
			archiveName: "tldr-pages.en.zip",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/tldr-pages.en.zip", r.URL.Path)
				_, _ = w.Write([]byte("zip-content"))
			},
			wantData: "zip-content",
		},
		{
			name:        "empty_hash_matches_any_content",
			archiveName: "archive.zip",
			hash:        "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("anything"))
			},
			wantData: "anything",
		},
		{
			name:        "hash_matches",
			archiveName: "tldr-pages.de.zip",
			hash: func() string {
				h := sha256.Sum256([]byte("de-content"))
				return hex.EncodeToString(h[:])
			}(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("de-content"))
			},
			wantData: "de-content",
		},
		{
			name:        "hash_mismatch",
			archiveName: "data.zip",
			hash:        "0000000000000000000000000000000000000000000000000000000000000000",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("actual-data"))
			},
			wantErr: true,
		},
		{
			name:        "server_error",
			archiveName: "missing.zip",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			client := upstream.New(
				upstream.WithHTTPClient(ts.Client()),
			)

			got, err := downloadArchive(
				context.Background(),
				client,
				ts.URL,
				tt.archiveName,
				tt.hash,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantData, string(got))
		})
	}
}

func TestExtractArchive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantErr           bool
		name              string
		languageDirectory string
		buildZip          func(t *testing.T) []byte
		preExist          func(t *testing.T, c *Cache)
		check             func(t *testing.T, c *Cache)
	}{
		{
			name:              "flat_structure",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"common/git.md": "",
					"common/ls.md":  "",
				})
			},
			check: func(t *testing.T, c *Cache) {
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "ls.md"))
			},
		},
		{
			name:              "nested_directories",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"common/git.md": "",
					"linux/apt.md":  "",
					"osx/brew.md":   "",
				})
			},
			check: func(t *testing.T, c *Cache) {
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "linux", "apt.md"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "osx", "brew.md"))
			},
		},
		{
			name:              "directory_entries_in_zip",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"common/":       "",
					"common/git.md": "",
				})
			},
			check: func(t *testing.T, c *Cache) {
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))
			},
		},
		{
			name:              "empty_zip",
			languageDirectory: "pages.en",
			buildZip:          createEmptyZip,
			check: func(t *testing.T, c *Cache) {
				dir := filepath.Join(c.dir, "pages.en")
				assert.DirExists(t, dir)

				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Empty(t, entries)
			},
		},
		{
			name:              "invalid_zip_data",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return []byte("not a zip file")
			},
			wantErr: true,
		},
		{
			name:              "skips_path_traversal",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"../escape.md":  "EVIL",
					"common/git.md": "",
				})
			},
			check: func(t *testing.T, c *Cache) {
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))
				assert.NoFileExists(t, filepath.Join(c.dir, "escape.md"))
			},
		},
		{
			name:              "removes_existing_directory",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"common/git.md": "",
				})
			},
			preExist: func(t *testing.T, c *Cache) {
				oldDir := filepath.Join(c.dir, "pages.en", "common")
				require.NoError(t, os.MkdirAll(oldDir, 0o750))
				require.NoError(t, os.WriteFile(
					filepath.Join(oldDir, "old.md"),
					[]byte("old"),
					0o640,
				))
			},
			check: func(t *testing.T, c *Cache) {
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))
				assert.NoFileExists(t, filepath.Join(c.dir, "pages.en", "common", "old.md"))
			},
		},
		{
			name:              "single_file_content",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"common/git.md": "# git\n",
				})
			},
			check: func(t *testing.T, c *Cache) {
				got, err := os.ReadFile(filepath.Join(c.dir, "pages.en", "common", "git.md"))
				require.NoError(t, err)
				assert.Equal(t, "# git\n", string(got))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{dir: t.TempDir()}

			if tt.preExist != nil {
				tt.preExist(t, c)
			}

			zipData := tt.buildZip(t)
			err := c.extractArchive(tt.languageDirectory, zipData)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

// TestExtractFile tests the extractFile helper.
func TestExtractFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
	}{
		{name: "writes_file_content", content: "# git\n"},
		{name: "empty_file", content: ""},
		{name: "binary_content", content: "\x00\x01\x02"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			root, err := os.OpenRoot(dir)
			require.NoError(t, err)
			defer func() {
				_ = root.Close()
			}()

			zipData := createTestZip(t, map[string]string{"test.md": tt.content})
			zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
			require.NoError(t, err)
			require.Len(t, zipReader.File, 1)

			f := zipReader.File[0]
			err = extractFile(root, f)
			require.NoError(t, err)

			got, err := os.ReadFile(filepath.Join(dir, "test.md"))
			require.NoError(t, err)
			assert.Equal(t, tt.content, string(got))
		})
	}
}
