package upstream

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var openSHA256 = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

type errReader struct{ err error }

func (r *errReader) Read(p []byte) (int, error) { return 0, r.err }

func TestTransfer(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("io error")
	tests := []struct {
		name           string
		maxBodySize    int64
		source         string
		expectedSHA256 string
		wantN          int64
		wantErr        bool
		errContains    string
	}{
		{
			name:   "simple_copy",
			source: "hello world",
			wantN:  11,
		},
		{
			name:           "matching_sha256",
			source:         "hello world",
			expectedSHA256: openSHA256,
			wantN:          11,
		},
		{
			name:           "mismatching_sha256",
			source:         "hello world",
			expectedSHA256: "0000000000000000000000000000000000000000000000000000000000000000",
			wantErr:        true,
		},
		{
			name:        "body_exceeds_maxBodySize",
			maxBodySize: 5,
			source:      "hello world",
			wantErr:     true, errContains: "exceeds limit",
		},
		{
			name:        "body_within_maxBodySize",
			maxBodySize: 20,
			source:      "hello world",
			wantN:       11,
		},
		{
			name:        "source_error",
			wantErr:     true,
			errContains: "copy:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &Client{
				maxBodySize: tt.maxBodySize,
			}

			var dst bytes.Buffer
			var src io.Reader

			if tt.name == "source_error" {
				src = &errReader{
					err: errExpected,
				}
			} else {
				src = strings.NewReader(tt.source)
			}

			n, err := c.transfer(&dst, src, tt.expectedSHA256)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantN, n)

				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}

				if tt.name == "mismatching_sha256" {
					assert.Equal(t, tt.source, dst.String())
				}

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantN, n)
			assert.Equal(t, tt.source, dst.String())
		})
	}
}

func TestDownloadFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		sha256hex   string
		wantErr     bool
		wantContent string
	}{
		{
			name: "creates_file_with_correct_content",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("hello world"))
			},
			sha256hex:   openSHA256,
			wantContent: "hello world",
		},
		{
			name: "sha256_mismatch_cleans_up_file",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("hello world"))
			},
			sha256hex: "0000000000000000000000000000000000000000000000000000000000000000",
			wantErr:   true,
		},
		{
			name: "http_error_does_not_create_file",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			dest := filepath.Join(t.TempDir(), "test.out")

			c := New(WithHTTPClient(ts.Client()))

			ctx := context.Background()
			err := c.DownloadFile(ctx, ts.URL, tt.sha256hex, dest)

			if tt.wantErr {
				require.Error(t, err)
				_, statErr := os.Stat(dest)
				assert.True(t, os.IsNotExist(statErr))
				return
			}

			require.NoError(t, err)
			data, err := os.ReadFile(dest)
			require.NoError(t, err)
			assert.Equal(t, tt.wantContent, string(data))
		})
	}
}
