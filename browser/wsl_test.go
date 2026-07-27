package browser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWSL(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "wsl2_marker",
			content: "Linux version 5.10.16.3-microsoft-standard-WSL2",
			want:    true,
		},
		{
			name:    "microsoft_mixed_case",
			content: "Linux version 5.15.0 (Microsoft@Microsoft.com) (gcc 11.2.0) #1 SMP",
			want:    true,
		},
		{
			name:    "no_microsoft_marker",
			content: "Linux version 6.1.0-23-amd64 (debian-kernel@lists.debian.org)",
			want:    false,
		},
		{
			name:    "empty_content",
			content: "",
			want:    false,
		},
		{
			name:    "partial_word_micah",
			content: "Linux version 5.15.0 (micah@kernel.org)",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reader := func(_ string) ([]byte, error) {
				return []byte(tt.content), nil
			}
			assert.Equal(t, tt.want, isWSL(reader))
		})
	}
}

func TestIsWSLReadError(t *testing.T) {
	t.Parallel()

	reader := func(_ string) ([]byte, error) {
		return nil, fmt.Errorf("permission denied")
	}

	assert.False(t, isWSL(reader))
}

func TestRunningOnWSL(t *testing.T) {
	t.Parallel()

	// smoke test: verify runningOnWSL doesn't panic and returns a bool.
	result := runningOnWSL()
	assert.IsType(t, false, result)
}
