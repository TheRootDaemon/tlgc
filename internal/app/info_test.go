package app

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestFormatPlatformBreakdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		platforms []cache.PlatformInfo
		want      string
	}{
		{
			name:      "empty",
			platforms: nil,
			want:      "",
		},
		{
			name: "single_platform",
			platforms: []cache.PlatformInfo{
				{Name: "common", Pages: 3000},
			},
			want: " (common: 3000)",
		},
		{
			name: "multiple_platforms",
			platforms: []cache.PlatformInfo{
				{Name: "common", Pages: 3000},
				{Name: "linux", Pages: 2500},
			},
			want: " (common: 3000, linux: 2500)",
		},
		{
			name: "three_platforms",
			platforms: []cache.PlatformInfo{
				{Name: "common", Pages: 3000},
				{Name: "linux", Pages: 2500},
				{Name: "osx", Pages: 2100},
			},
			want: " (common: 3000, linux: 2500, osx: 2100)",
		},
		{
			name: "zero_pages",
			platforms: []cache.PlatformInfo{
				{Name: "windows", Pages: 0},
			},
			want: " (windows: 0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripANSI(formatPlatformBreakdown(tt.platforms))
			assert.Equal(t, tt.want, got)
		})
	}
}
