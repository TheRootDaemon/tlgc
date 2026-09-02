package cache

import (
	"context"
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

	zipEn := createTestZip(t, map[string]string{"common/git.md": ""})
	hashEn := contentHash(string(zipEn))

	zipDe := createTestZip(t, map[string]string{"common/apt.md": ""})
	hashDe := contentHash(string(zipDe))

	tests := []struct {
		name      string
		languages []string
		preExist  func(t *testing.T, c *Cache)
		checksums map[string]string
		archives  map[string][]byte
		handler   http.Handler
		wantErr   bool
		check     func(t *testing.T, c *Cache)
	}{
		{
			name:      "fresh_update",
			languages: []string{"en"},
			checksums: map[string]string{"tldr-pages.en.zip": hashEn},
			archives:  map[string][]byte{"tldr-pages.en.zip": zipEn},
			check: func(t *testing.T, c *Cache) {
				assert.DirExists(t, filepath.Join(c.dir, "pages.en", "common"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))

				data, err := os.ReadFile(filepath.Join(c.dir, checksumFile))
				require.NoError(t, err)
				assert.Contains(t, string(data), hashEn)
			},
		},
		{
			name:      "already_up_to_date",
			languages: []string{"en"},
			preExist: func(t *testing.T, c *Cache) {
				require.NoError(t, os.MkdirAll(filepath.Join(c.dir, "pages.en", "common"), 0o750))
				require.NoError(t, c.saveChecksums(map[string]string{"tldr-pages.en.zip": hashEn}))
			},
			checksums: map[string]string{"tldr-pages.en.zip": hashEn},
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
				require.NoError(t, c.saveChecksums(map[string]string{"tldr-pages.en.zip": hashEn}))
			},
			checksums: map[string]string{
				"tldr-pages.en.zip": hashEn,
				"tldr-pages.de.zip": hashDe,
			},
			archives: map[string][]byte{"tldr-pages.de.zip": zipDe},
			check: func(t *testing.T, c *Cache) {
				assert.DirExists(t, filepath.Join(c.dir, "pages.en", "common"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))

				assert.DirExists(t, filepath.Join(c.dir, "pages.de", "common"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.de", "common", "apt.md"))

				data, err := os.ReadFile(filepath.Join(c.dir, checksumFile))
				require.NoError(t, err)
				assert.Contains(t, string(data), hashEn)
				assert.Contains(t, string(data), hashDe)
			},
		},
		{
			name:      "checksum_download_fails",
			languages: []string{"en"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.handler
			if handler == nil {
				handler = archiveServer(t, tt.checksums, tt.archives)
			}

			ts := httptest.NewServer(handler)
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

	zipData := createTestZip(t, map[string]string{"common/git.md": ""})
	correctHash := contentHash(string(zipData))

	invalidZipData := []byte("not a valid zip file")
	invalidHash := contentHash(string(invalidZipData))

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

			gotUpdated, _, err := c.updateLanguage(
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
