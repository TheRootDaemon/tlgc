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

	t.Run("saves_single_entry", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		err := c.saveChecksums(map[string]string{
			"en.zip": "abc",
		})
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir, checksumFile))
		require.NoError(t, err)
		assert.Equal(t, "abc  en.zip\n", string(data))
	})

	t.Run("saves_multiple_entries_and_round_trip", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		original := map[string]string{
			"en.zip": "abc",
			"de.zip": "def",
			"zh.zip": "ghi",
		}

		err := c.saveChecksums(original)
		require.NoError(t, err)

		got := c.loadChecksums()
		assert.Equal(t, original, got)
	})

	t.Run("overwrites_existing_file", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		err := os.WriteFile(
			filepath.Join(dir, checksumFile),
			[]byte("oldhash  old.zip\n"),
			0o600,
		)
		require.NoError(t, err)

		err = c.saveChecksums(map[string]string{
			"new.zip": "newhash",
		})
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir, checksumFile))
		require.NoError(t, err)
		assert.Equal(t, "newhash  new.zip\n", string(data))
	})

	t.Run("creates_directory", func(t *testing.T) {
		base := t.TempDir()
		nested := filepath.Join(base, "sub", "dir")
		c := &Cache{dir: nested}

		err := c.saveChecksums(map[string]string{
			"a.zip": "h",
		})
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(nested, checksumFile))
		require.NoError(t, err)
		assert.Equal(t, "h  a.zip\n", string(data))
	})

	t.Run("empty_map", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		err := c.saveChecksums(map[string]string{})
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir, checksumFile))
		require.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("special_chars_in_filename", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		err := c.saveChecksums(map[string]string{
			"f!@#.zip": "abc123",
		})
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir, checksumFile))
		require.NoError(t, err)
		assert.Equal(t, "abc123  f!@#.zip\n", string(data))
	})

	t.Run("round_trip_empty_map", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		err := c.saveChecksums(map[string]string{})
		require.NoError(t, err)

		got := c.loadChecksums()
		assert.Equal(t, map[string]string{}, got)
	})

	t.Run("round_trip_large_map", func(t *testing.T) {
		dir := t.TempDir()
		c := &Cache{dir: dir}

		original := make(map[string]string)
		for i := range 20 {
			name := fmt.Sprintf("tldr-pages.%d.zip", i)
			hash := fmt.Sprintf("%064d", i)
			original[name] = hash
		}

		err := c.saveChecksums(original)
		require.NoError(t, err)

		got := c.loadChecksums()
		assert.Equal(t, original, got)
	})
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
