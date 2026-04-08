package termcolor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupportsColor(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		isTerminal bool
		expected   bool
	}{
		{
			name: "NO_COLOR disables color",
			env: map[string]string{
				"NO_COLOR": "1",
			},
			isTerminal: true,
			expected:   false,
		},
		{
			name: "TERM dumb disables color",
			env: map[string]string{
				"TERM": "dumb",
			},
			isTerminal: true,
			expected:   false,
		},
		{
			name:       "not a terminal disables color",
			env:        map[string]string{},
			isTerminal: false,
			expected:   false,
		},
		{
			name:       "valid terminal supports color",
			env:        map[string]string{},
			isTerminal: true,
			expected:   true,
		},
	}

	oldGetenv := getenv
	oldIsTerminal := isTerminal
	t.Cleanup(func() {
		getenv = oldGetenv
		isTerminal = oldIsTerminal
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mocks getenv, isTerminal
			getenv = func(key string) string {
				return tt.env[key]
			}
			isTerminal = func(fd uintptr) bool {
				return tt.isTerminal
			}

			got := SupportsColor()
			require.Equal(t, tt.expected, got)
		})
	}
}
