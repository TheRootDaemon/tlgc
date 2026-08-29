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
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TheRootDaemon/tlgc/internal/upstream"
)

// contentHash returns the hex-encoded sha256 of s,
// matching the hashes produced by the cache package.
func contentHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

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
			name:              "skips_root_level_files",
			languageDirectory: "pages.en",
			buildZip: func(t *testing.T) []byte {
				return createTestZip(t, map[string]string{
					"LICENSE.md":    "license",
					"common/git.md": "",
				})
			},
			check: func(t *testing.T, c *Cache) {
				assert.NoFileExists(t, filepath.Join(c.dir, "pages.en", "LICENSE.md"))
				assert.FileExists(t, filepath.Join(c.dir, "pages.en", "common", "git.md"))
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
			_, err := c.extractArchive(tt.languageDirectory, zipData)

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
			hash, err := extractFile(root, f)
			require.NoError(t, err)
			assert.Equal(t, contentHash(tt.content), hash)

			got, err := os.ReadFile(filepath.Join(dir, "test.md"))
			require.NoError(t, err)
			assert.Equal(t, tt.content, string(got))
		})
	}
}

func TestExtractArchive_MkdirAllError(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "notadir")
	require.NoError(t, os.WriteFile(filePath, nil, 0o644))

	c := &Cache{dir: filePath}
	_, err := c.extractArchive("pages.en", createEmptyZip(t))
	assert.Error(t, err)
}

func TestExtractFile_OpenFileError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod 0o500 does not prevent file creation on Windows")
	}

	dir := t.TempDir()
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	defer func() {
		_ = root.Close()
	}()

	zipData := createTestZip(t, map[string]string{"test.md": "content"})
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	require.NoError(t, err)
	require.Len(t, zipReader.File, 1)

	f := zipReader.File[0]
	_, err = extractFile(root, f)
	assert.Error(t, err)
}

func TestHashPages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T, dir string)
		want  map[string]string
	}{
		{
			name:  "nonexistent_directory",
			setup: nil,
			want:  map[string]string{},
		},
		{
			name: "empty_directory",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(dir, 0o750))
			},
			want: map[string]string{},
		},
		{
			name: "hashes_nested_pages",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "common"), 0o750))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "common", "git.md"), []byte("# git\n"), 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "common", "ls.md"), nil, 0o644))
			},
			want: map[string]string{
				filepath.Join("common", "git.md"): contentHash("# git\n"),
				filepath.Join("common", "ls.md"):  contentHash(""),
			},
		},
		{
			name: "ignores_non_pages",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "common"), 0o750))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "LICENSE.md"), []byte("license"), 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "common", "notes.txt"), []byte("notes"), 0o644))
			},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "subdir")
			if tt.setup != nil {
				tt.setup(t, dir)
			}

			got, err := hashPages(dir)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractZip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		files        map[string]string
		rawData      []byte            // when set, used instead of building a zip from files
		oldPages     map[string]string // relative page path -> old content
		wantTotal    int
		wantAdded    int
		wantModified int
		wantRemoved  int
		wantErr      bool
	}{
		{
			name: "all_new",
			files: map[string]string{
				"common/git.md": "",
				"common/ls.md":  "",
			},
			wantTotal: 2,
			wantAdded: 2,
		},
		{
			name:      "all_existing",
			files:     map[string]string{"common/git.md": ""},
			oldPages:  map[string]string{"common/git.md": ""},
			wantTotal: 1,
		},
		{
			name: "modified_content",
			files: map[string]string{
				"common/git.md": "# git v2\n",
			},
			oldPages:     map[string]string{"common/git.md": "# git v1\n"},
			wantTotal:    1,
			wantModified: 1,
		},
		{
			name: "added_and_removed",
			files: map[string]string{
				"common/ls.md": "",
			},
			oldPages:    map[string]string{"common/git.md": ""},
			wantTotal:   1,
			wantAdded:   1,
			wantRemoved: 1,
		},
		{
			name: "mixed_changes",
			files: map[string]string{
				"common/a.md": "a",
				"common/b.md": "b2",
				"common/c.md": "c",
			},
			oldPages: map[string]string{
				"common/b.md": "b1",
				"common/c.md": "c",
				"common/d.md": "d",
			},
			wantTotal:    3,
			wantAdded:    1,
			wantModified: 1,
			wantRemoved:  1,
		},
		{
			name: "directory_entries_not_counted",
			files: map[string]string{
				"common/":       "",
				"common/git.md": "",
			},
			wantTotal: 1,
			wantAdded: 1,
		},
		{
			name: "skips_unsafe_entries",
			files: map[string]string{
				"../escape.md":  "EVIL",
				"common/git.md": "",
			},
			wantTotal: 1,
			wantAdded: 1,
		},
		{
			name: "skips_root_level_files",
			files: map[string]string{
				"LICENSE.md":    "",
				"common/git.md": "",
			},
			wantTotal: 1,
			wantAdded: 1,
		},
		{
			name: "skips_nested_non_page_files",
			files: map[string]string{
				"common/notes.txt": "notes",
				"common/git.md":    "",
			},
			wantTotal: 1,
			wantAdded: 1,
		},
		{
			name:        "empty_archive_removes_all",
			oldPages:    map[string]string{"common/git.md": ""},
			wantRemoved: 1,
		},
		{
			name:    "invalid_zip",
			rawData: []byte("not a zip file"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			root, err := os.OpenRoot(dir)
			require.NoError(t, err)
			defer func() { _ = root.Close() }()

			oldPages := make(map[string]string, len(tt.oldPages))
			for name, content := range tt.oldPages {
				oldPages[name] = contentHash(content)
			}

			var zipData []byte
			if tt.rawData != nil {
				zipData = tt.rawData
			} else if tt.files != nil {
				zipData = createTestZip(t, tt.files)
			} else {
				zipData = createEmptyZip(t)
			}

			st, err := extractZip(root, zipData, oldPages)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, st.totalPages)
			assert.Equal(t, tt.wantAdded, st.added)
			assert.Equal(t, tt.wantModified, st.modified)
			assert.Equal(t, tt.wantRemoved, st.removed)
		})
	}
}

func TestRecreateDirectory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T, dir string)
		check func(t *testing.T, dir string)
	}{
		{
			name: "creates_new_directory",
			check: func(t *testing.T, dir string) {
				assert.DirExists(t, dir)
			},
		},
		{
			name: "replaces_existing_directory",
			setup: func(t *testing.T, dir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(dir, "sub"), 0o750))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "old.md"), nil, 0o644))
			},
			check: func(t *testing.T, dir string) {
				assert.DirExists(t, dir)
				entries, err := os.ReadDir(dir)
				require.NoError(t, err)
				assert.Empty(t, entries)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "dir")
			if tt.setup != nil {
				tt.setup(t, dir)
			}
			require.NoError(t, recreateDirectory(dir))
			tt.check(t, dir)
		})
	}
}

func TestExtractZipEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    map[string]string
		wantPage bool
		wantHash string
		check    func(t *testing.T, dir string)
	}{
		{
			name:     "extracts_file",
			files:    map[string]string{"common/git.md": "# git\n"},
			wantPage: true,
			wantHash: contentHash("# git\n"),
			check: func(t *testing.T, dir string) {
				got, err := os.ReadFile(filepath.Join(dir, "common", "git.md"))
				require.NoError(t, err)
				assert.Equal(t, "# git\n", string(got))
			},
		},
		{
			name:     "directory_entry",
			files:    map[string]string{"common/": ""},
			wantPage: false,
			check: func(t *testing.T, dir string) {
				assert.DirExists(t, filepath.Join(dir, "common"))
			},
		},
		{
			name:     "unsafe_entry",
			files:    map[string]string{"../escape.md": "EVIL"},
			wantPage: false,
		},
		{
			name:     "skips_root_level_file",
			files:    map[string]string{"LICENSE.md": "license"},
			wantPage: false,
			check: func(t *testing.T, dir string) {
				assert.NoFileExists(t, filepath.Join(dir, "LICENSE.md"))
			},
		},
		{
			name:     "skips_nested_non_page_file",
			files:    map[string]string{"common/notes.txt": "notes"},
			wantPage: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			root, err := os.OpenRoot(dir)
			require.NoError(t, err)
			defer func() { _ = root.Close() }()

			zipData := createTestZip(t, tt.files)
			zr, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
			require.NoError(t, err)
			require.Len(t, zr.File, 1)

			isPage, hash, err := extractZipEntry(root, zr.File[0])
			require.NoError(t, err)
			assert.Equal(t, tt.wantPage, isPage)
			assert.Equal(t, tt.wantHash, hash)

			if tt.check != nil {
				tt.check(t, dir)
			}
		})
	}
}
