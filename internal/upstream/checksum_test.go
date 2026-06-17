package upstream

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifySHA256(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "matching hash",
			data:     []byte("hello world"),
			expected: fmt.Sprintf("%x", sha256.Sum256([]byte("hello world"))),
			wantErr:  false,
		},
		{
			name:     "mismatching hash",
			data:     []byte("hello world"),
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
			wantErr:  true,
		},
		{
			name:     "with filename format",
			data:     []byte("test data"),
			expected: fmt.Sprintf("%x", sha256.Sum256([]byte("test data"))) + "  myfile.txt",
			wantErr:  false,
		},
		{
			name:     "with binary mode",
			data:     []byte("test data"),
			expected: fmt.Sprintf("%x", sha256.Sum256([]byte("test data"))) + " *myfile.bin",
			wantErr:  false,
		},
		{
			name:     "empty expected",
			data:     []byte("some data"),
			expected: "",
			wantErr:  false,
		},

		{
			name:     "whitespace only",
			data:     []byte("some data"),
			expected: "  \t  ",
			wantErr:  false,
		},

		{
			name:     "case insensitive",
			data:     []byte("case test"),
			expected: fmt.Sprintf("%X", sha256.Sum256([]byte("case test"))),
			wantErr:  false,
		},
		{
			name:     "filename format wrong hash",
			data:     []byte("hello world"),
			expected: "0000000000000000000000000000000000000000000000000000000000000000  file.txt",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySHA256(tt.data, tt.expected)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestParseCheckSum(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		wantHash     string
		wantFilename string
		wantErr      bool
		wantErrMsg   string
	}{
		{
			name:         "valid line",
			line:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  filename.txt",
			wantHash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantFilename: "filename.txt",
		},
		{
			name:         "binary mode",
			line:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  *binary.bin",
			wantHash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantFilename: "binary.bin",
		},
		{
			name:         "single space",
			line:         "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 file.txt",
			wantHash:     "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			wantFilename: "file.txt",
		},
		{
			name:         "with path",
			line:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  sub/dir/file.txt",
			wantHash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantFilename: "sub/dir/file.txt",
		},
		{
			name:       "empty line",
			line:       "",
			wantErr:    true,
			wantErrMsg: "empty",
		},
		{
			name:       "whitespace only",
			line:       "   ",
			wantErr:    true,
			wantErrMsg: "empty",
		},
		{
			name:       "invalid format",
			line:       "justoneword",
			wantErr:    true,
			wantErrMsg: "invalid checksum line",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, filename, err := ParseCheckSum(tt.line)
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantHash, hash)
			assert.Equal(t, tt.wantFilename, filename)
		})
	}
}
