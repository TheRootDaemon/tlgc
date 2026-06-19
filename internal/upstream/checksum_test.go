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

func TestVerifySHA256hex(t *testing.T) {
	validHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	tests := []struct {
		name     string
		got      string
		expected string
		wantErr  bool
	}{
		{
			name:     "matching hex",
			got:      validHash,
			expected: validHash,
			wantErr:  false,
		},
		{
			name:     "mismatching hex",
			got:      validHash,
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
			wantErr:  true,
		},
		{
			name:     "empty expected",
			got:      validHash,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "whitespace only expected",
			got:      validHash,
			expected: "  \t  ",
			wantErr:  false,
		},
		{
			name:     "filename format matching",
			got:      validHash,
			expected: validHash + "  filename.txt",
			wantErr:  false,
		},
		{
			name:     "binary mode matching",
			got:      validHash,
			expected: validHash + " *binary.bin",
			wantErr:  false,
		},
		{
			name:     "case insensitive uppercase expected",
			got:      validHash,
			expected: "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
			wantErr:  false,
		},
		{
			name:     "case insensitive uppercase got",
			got:      "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
			expected: validHash,
			wantErr:  false,
		},
		{
			name:     "filename format wrong hash",
			got:      validHash,
			expected: "0000000000000000000000000000000000000000000000000000000000000000  file.txt",
			wantErr:  true,
		},
		{
			name:     "short string direct compare",
			got:      "short",
			expected: "short",
			wantErr:  false,
		},
		{
			name:     "short string mismatch",
			got:      "short",
			expected: "different",
			wantErr:  true,
		},
		{
			name:     "filename with trailing spaces",
			got:      validHash,
			expected: validHash + "  file.txt  ",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySHA256hex(tt.got, tt.expected)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestVerifySHA256Hash(t *testing.T) {
	validHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	data := []byte("hello world")

	tests := []struct {
		name     string
		data     []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "matching hash",
			data:     data,
			expected: validHash,
			wantErr:  false,
		},
		{
			name:     "mismatching hash",
			data:     data,
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
			wantErr:  true,
		},
		{
			name:     "filename format",
			data:     data,
			expected: validHash + "  myfile.txt",
			wantErr:  false,
		},
		{
			name:     "empty expected",
			data:     data,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "case insensitive",
			data:     data,
			expected: "B94D27B9934D3E08A52E52D7DA7DABFAC484EFE37A5380EE9088F7ACE2EFCDE9",
			wantErr:  false,
		},
		{
			name:     "binary mode",
			data:     data,
			expected: validHash + " *binary.bin",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := sha256.New()
			_, err := h.Write(tt.data)
			require.NoError(t, err)

			err = verifySHA256Hash(h, tt.expected)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestParseChecksum(t *testing.T) {
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
			hash, filename, err := ParseChecksum(tt.line)
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
