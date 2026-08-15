package cache

import (
	"context"
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

func TestLoadChecksums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		want     map[string]string
		wantNil  bool
	}{
		{
			name: "file_exists_valid",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(
					filepath.Join(dir, checksumFile),
					[]byte("abc111  en.zip\ndef222  de.zip\n"),
					0o600,
				)
				require.NoError(t, err)
				return dir
			},
			want: map[string]string{
				"en.zip": "abc111",
				"de.zip": "def222",
			},
		},
		{
			name: "file_does_not_exist",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			wantNil: true,
		},
		{
			name: "dir_does_not_exist",
			setupDir: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			wantNil: true,
		},
		{
			name: "file_empty",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(
					filepath.Join(dir, checksumFile),
					[]byte(""),
					0o600,
				)
				require.NoError(t, err)
				return dir
			},
			want: map[string]string{},
		},
		{
			name: "file_whitespace_only",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(
					filepath.Join(dir, checksumFile),
					[]byte("\n\n  \n"),
					0o600,
				)
				require.NoError(t, err)
				return dir
			},
			want: map[string]string{},
		},
		{
			name: "file_mixed_valid_invalid",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(
					filepath.Join(dir, checksumFile),
					[]byte("abc111  good.zip\nbadline\nabc222  ok.zip\n"),
					0o600,
				)
				require.NoError(t, err)
				return dir
			},
			want: map[string]string{
				"good.zip": "abc111",
				"ok.zip":   "abc222",
			},
		},
		{
			name: "full_sha256_hash",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(
					filepath.Join(dir, checksumFile),
					[]byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  tldr-pages.en.zip\n"),
					0o600,
				)
				require.NoError(t, err)
				return dir
			},
			want: map[string]string{
				"tldr-pages.en.zip": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}
			got := c.loadChecksums()

			if tt.wantNil {
				assert.Nil(t, got)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSaveChecksums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setupDir func(t *testing.T) string
		data     map[string]string
		wantFile string
	}{
		{
			name: "saves_single_entry",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			data: map[string]string{
				"en.zip": "abc",
			},
			wantFile: "abc  en.zip\n",
		},
		{
			name: "empty_map",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			data:     map[string]string{},
			wantFile: "",
		},
		{
			name: "special_chars_in_filename",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			data: map[string]string{
				"f!@#.zip": "abc123",
			},
			wantFile: "abc123  f!@#.zip\n",
		},
		{
			name: "overwrites_existing_file",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				err := os.WriteFile(
					filepath.Join(dir, checksumFile),
					[]byte("oldhash  old.zip\n"),
					0o600,
				)
				require.NoError(t, err)
				return dir
			},
			data: map[string]string{
				"new.zip": "newhash",
			},
			wantFile: "newhash  new.zip\n",
		},
		{
			name: "creates_directory",
			setupDir: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "sub", "dir")
			},
			data: map[string]string{
				"a.zip": "h",
			},
			wantFile: "h  a.zip\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupDir(t)
			c := &Cache{dir: dir}

			err := c.saveChecksums(tt.data)
			require.NoError(t, err)

			data, err := os.ReadFile(filepath.Join(dir, checksumFile))
			require.NoError(t, err)
			assert.Equal(t, tt.wantFile, string(data))
		})
	}
}

func TestSaveChecksumsRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data map[string]string
	}{
		{
			name: "multiple_entries",
			data: map[string]string{
				"en.zip": "abc",
				"de.zip": "def",
				"zh.zip": "ghi",
			},
		},
		{
			name: "empty_map",
			data: map[string]string{},
		},
		{
			name: "large_map",
			data: func() map[string]string {
				data := make(map[string]string)
				for i := range 20 {
					name := fmt.Sprintf("tldr-pages.%d.zip", i)
					hash := fmt.Sprintf("%064d", i)
					data[name] = hash
				}
				return data
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			c := &Cache{dir: dir}

			require.NoError(t, c.saveChecksums(tt.data))

			got := c.loadChecksums()
			assert.Equal(t, tt.data, got)
		})
	}
}

func TestDownloadChecksum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		want    string
	}{
		{
			name: "successful_download",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/"+checksumFile, r.URL.Path)
				_, _ = w.Write([]byte("hash  en.zip\n"))
			},
			want: "hash  en.zip\n",
		},
		{
			name: "empty_response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(nil)
			},
			want: "",
		},
		{
			name: "server_error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "not_found",
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

			got, err := downloadChecksum(
				context.Background(),
				client,
				ts.URL,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestParseChecksum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []byte
		want  map[string]string
	}{
		{
			name:  "single_valid_line",
			input: []byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  tldr-pages.en.zip\n"),
			want: map[string]string{
				"tldr-pages.en.zip": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
		{
			name: "multiple_languages",
			input: []byte(
				"abc111  tldr-pages.en.zip\n" +
					"def222  tldr-pages.de.zip\n" +
					"ghi333  tldr-pages.zh.zip\n",
			),
			want: map[string]string{
				"tldr-pages.en.zip": "abc111",
				"tldr-pages.de.zip": "def222",
				"tldr-pages.zh.zip": "ghi333",
			},
		},
		{
			name: "skips_empty_lines",
			input: []byte(
				"abc111  a.zip\n\n\n\nabc222  b.zip\n",
			),
			want: map[string]string{
				"a.zip": "abc111",
				"b.zip": "abc222",
			},
		},
		{
			name: "skips_invalid_lines",
			input: []byte(
				"abc111  good.zip\n" +
					"badline\n" +
					"short\n",
			),
			want: map[string]string{
				"good.zip": "abc111",
			},
		},
		{
			name:  "all_lines_invalid",
			input: []byte("garbage\nshort  \n"),
			want:  map[string]string{},
		},
		{
			name:  "empty_input",
			input: []byte(""),
			want:  map[string]string{},
		},
		{
			name:  "only_whitespace",
			input: []byte("  \n\t\n  \n"),
			want:  map[string]string{},
		},
		{
			name:  "binary_mode_strips_star",
			input: []byte("abc123  *tldr-pages.en.zip\n"),
			want: map[string]string{
				"tldr-pages.en.zip": "abc123",
			},
		},
		{
			name:  "filename_with_path",
			input: []byte("abc123  sub/dir/file.txt\n"),
			want: map[string]string{
				"sub/dir/file.txt": "abc123",
			},
		},
		{
			name: "trailing_newline",
			input: []byte(
				"abc111  a.zip\n" +
					"abc222  b.zip\n",
			),
			want: map[string]string{
				"a.zip": "abc111",
				"b.zip": "abc222",
			},
		},
		{
			name: "duplicate_filename_last_wins",
			input: []byte(
				"abc111  f.zip\n" +
					"abc222  f.zip\n",
			),
			want: map[string]string{
				"f.zip": "abc222",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseChecksum(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
