package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TheRootDaemon/tlgc/internal/upstream"
)

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	// pre-compute a valid ZIP for english and its hash.
	zipData := createTestZip(t, map[string]string{"common/git.md": ""})
	h := sha256.Sum256(zipData)
	correctHash := hex.EncodeToString(h[:])

	// pre-compute a valid ZIP for german and its hash.
	zipDataDe := createTestZip(t, map[string]string{"common/apt.md": ""})
	hDe := sha256.Sum256(zipDataDe)
	correctHashDe := hex.EncodeToString(hDe[:])

	tests := []struct {
		name      string
		languages []string
		preExist  func(t *testing.T, c *Cache)
		handler   http.HandlerFunc
		wantErr   bool
		check     func(t *testing.T, c *Cache)
	}{
		{
			name:      "fresh_update",
			languages: []string{"en"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/" + checksumFile:
					_, _ = fmt.Fprintf(w, "%s  %s\n", correctHash, "tldr-pages.en.zip")
				case "/tldr-pages.en.zip":
					_, _ = w.Write(zipData)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			check: func(t *testing.T, c *Cache) {
				assert.DirExists(t, filepath.Join(c.dir, "pages.en", "common"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))

				data, err := os.ReadFile(filepath.Join(c.dir, checksumFile))
				require.NoError(t, err)
				assert.Contains(t, string(data), correctHash)
			},
		},
		{
			name:      "already_up_to_date",
			languages: []string{"en"},
			preExist: func(t *testing.T, c *Cache) {
				require.NoError(t, os.MkdirAll(filepath.Join(c.dir, "pages.en", "common"), 0o750))
				require.NoError(t, c.saveChecksums(map[string]string{"tldr-pages.en.zip": correctHash}))
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/" + checksumFile:
					_, _ = fmt.Fprintf(w, "%s  %s\n", correctHash, "tldr-pages.en.zip")
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			check: func(t *testing.T, c *Cache) {
				assert.DirExists(t, filepath.Join(c.dir, "pages.en", "common"))
			},
		},
		{
			name:      "partial_update",
			languages: []string{"en", "de"},
			preExist: func(t *testing.T, c *Cache) {
				require.NoError(t, os.MkdirAll(filepath.Join(c.dir, "pages.en", "common"), 0o750))
				require.NoError(t, os.WriteFile(
					filepath.Join(c.dir, "pages.en", "common", "git.md"),
					[]byte("# git\n"), 0o640,
				))
				require.NoError(t, c.saveChecksums(map[string]string{"tldr-pages.en.zip": correctHash}))
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/" + checksumFile:
					_, _ = fmt.Fprintf(
						w, "%s  %s\n%s  %s\n",
						correctHash, "tldr-pages.en.zip",
						correctHashDe, "tldr-pages.de.zip",
					)
				case "/tldr-pages.de.zip":
					_, _ = w.Write(zipDataDe)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			},
			check: func(t *testing.T, c *Cache) {
				assert.DirExists(t, filepath.Join(c.dir, "pages.en", "common"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))

				assert.DirExists(t, filepath.Join(c.dir, "pages.de", "common"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.de", "common", "apt.md"))

				data, err := os.ReadFile(filepath.Join(c.dir, checksumFile))
				require.NoError(t, err)
				assert.Contains(t, string(data), correctHash)
				assert.Contains(t, string(data), correctHashDe)
			},
		},
		{
			name:      "checksum_download_fails",
			languages: []string{"en"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			cacheDir := t.TempDir()
			defer setupConfig(t, cacheDir, ts.URL)()

			if tt.preExist != nil {
				tt.preExist(t, &Cache{dir: cacheDir})
			}

			c := &Cache{dir: cacheDir}
			client := upstream.New(
				upstream.WithHTTPClient(ts.Client()),
			)

			err := c.Update(ctx, tt.languages, client)

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

func TestUpdateLanguage(t *testing.T) {
	ctx := context.Background()

	// pre-compute a valid ZIP and its SHA256 hash.
	zipData := createTestZip(t, map[string]string{"common/git.md": ""})
	h := sha256.Sum256(zipData)
	correctHash := hex.EncodeToString(h[:])

	// pre-compute garbage data and its hash for the invalid-zip case.
	invalidZipData := []byte("not a valid zip file")
	hi := sha256.Sum256(invalidZipData)
	invalidHash := hex.EncodeToString(hi[:])

	tests := []struct {
		name         string
		language     string
		preExist     func(t *testing.T, c *Cache)
		oldChecksums map[string]string
		newChecksums map[string]string
		handler      http.HandlerFunc
		wantUpdated  bool
		wantErr      bool
	}{
		{
			name:         "needs_update",
			language:     "en",
			oldChecksums: nil,
			newChecksums: map[string]string{"tldr-pages.en.zip": correctHash},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(zipData)
			},
			wantUpdated: true,
		},
		{
			name:     "up_to_date",
			language: "en",
			preExist: func(t *testing.T, c *Cache) {
				require.NoError(t, os.MkdirAll(filepath.Join(c.dir, "pages.en", "common"), 0o750))
			},
			oldChecksums: map[string]string{"tldr-pages.en.zip": correctHash},
			newChecksums: map[string]string{"tldr-pages.en.zip": correctHash},
			wantUpdated:  false,
		},
		{
			name:         "not_in_new_checksums",
			language:     "en",
			oldChecksums: nil,
			newChecksums: map[string]string{},
			wantUpdated:  false,
		},
		{
			name:     "hash_changed",
			language: "en",
			preExist: func(t *testing.T, c *Cache) {
				require.NoError(t, os.MkdirAll(filepath.Join(c.dir, "pages.en", "common"), 0o750))
			},
			oldChecksums: map[string]string{"tldr-pages.en.zip": "oldhash"},
			newChecksums: map[string]string{"tldr-pages.en.zip": correctHash},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(zipData)
			},
			wantUpdated: true,
		},
		{
			name:         "checksum_mismatch",
			language:     "en",
			oldChecksums: nil,
			newChecksums: map[string]string{"tldr-pages.en.zip": "0000000000000000000000000000000000000000000000000000000000000000"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(zipData)
			},
			wantUpdated: false,
			wantErr:     true,
		},
		{
			name:         "download_fails",
			language:     "en",
			oldChecksums: nil,
			newChecksums: map[string]string{"tldr-pages.en.zip": "irrelevant"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantUpdated: false,
			wantErr:     true,
		},
		{
			name:         "invalid_zip",
			language:     "en",
			oldChecksums: nil,
			newChecksums: map[string]string{"tldr-pages.en.zip": invalidHash},
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(invalidZipData)
			},
			wantUpdated: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			cacheDir := t.TempDir()
			defer setupConfig(t, cacheDir, ts.URL)()

			if tt.preExist != nil {
				tt.preExist(t, &Cache{dir: cacheDir})
			}

			c := &Cache{dir: cacheDir}
			client := upstream.New(upstream.WithHTTPClient(ts.Client()))

			gotUpdated, err := c.updateLanguage(
				ctx, client, tt.language,
				tt.oldChecksums, tt.newChecksums,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, gotUpdated)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUpdated, gotUpdated)
		})
	}
}

func TestNeedsUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		exists       bool
		archive      string
		oldChecksums map[string]string
		newChecksums map[string]string
		want         bool
	}{
		{
			name:         "archive_not_in_new_checksums",
			archive:      "tldr-pages.en.zip",
			oldChecksums: nil,
			newChecksums: map[string]string{},
			want:         false,
		},
		{
			name:         "no_old_hash",
			exists:       true,
			archive:      "tldr-pages.en.zip",
			oldChecksums: map[string]string{},
			newChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			want:         true,
		},
		{
			name:         "directory_missing",
			exists:       false,
			archive:      "tldr-pages.en.zip",
			oldChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			newChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			want:         true,
		},
		{
			name:         "hash_changed",
			exists:       true,
			archive:      "tldr-pages.en.zip",
			oldChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			newChecksums: map[string]string{"tldr-pages.en.zip": "def"},
			want:         true,
		},
		{
			name:         "up_to_date",
			exists:       true,
			archive:      "tldr-pages.en.zip",
			oldChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			newChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			want:         false,
		},
		{
			name:         "empty_new_checksums",
			exists:       true,
			archive:      "tldr-pages.en.zip",
			oldChecksums: map[string]string{"tldr-pages.en.zip": "abc"},
			newChecksums: map[string]string{},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsUpdate(tt.exists, tt.archive, tt.oldChecksums, tt.newChecksums)
			assert.Equal(t, tt.want, got)
		})
	}
}
