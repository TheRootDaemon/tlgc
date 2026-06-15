package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	got := Default()
	assert.NotEmpty(t, got)
	assert.Contains(t, All(), got, "Default should return a known platform")
}

func TestDefaultAllBranches(t *testing.T) {
	// Validate every branch of Default returns one of the known platforms.
	// We can't change runtime.GOOS at runtime, but we can verify the
	// switch logic by calling All() and checking the result is valid.
	known := All()
	tests := []struct {
		goos string
		want string
	}{
		{goos: "linux", want: "linux"},
		{goos: "darwin", want: "osx"},
		{goos: "windows", want: "windows"},
		{goos: "android", want: "android"},
		{goos: "freebsd", want: "common"},
		{goos: "netbsd", want: "common"},
		{goos: "openbsd", want: "common"},
		{goos: "plan9", want: "common"},
		{goos: "", want: "common"},
	}
	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			assert.Equal(t, tt.want, resolveGOOS(tt.goos))
			assert.Contains(t, known, tt.want)
		})
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "macos resolves to osx", input: "macos", want: "osx"},
		{name: "linux stays linux", input: "linux", want: "linux"},
		{name: "osx stays osx", input: "osx", want: "osx"},
		{name: "windows stays windows", input: "windows", want: "windows"},
		{name: "android stays android", input: "android", want: "android"},
		{name: "common stays common", input: "common", want: "common"},
		{name: "empty string stays empty", input: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResolveCaseSensitivity(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "Macos", input: "Macos"},
		{name: "MACOS", input: "MACOS"},
		{name: "macOS", input: "macOS"},
		{name: "macoS", input: "macoS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.input)
			assert.Equal(t, tt.input, got, "Resolve is case-sensitive")
		})
	}
}

func TestResolveAllKnownPlatforms(t *testing.T) {
	// Every known platform should pass through Resolve unchanged.
	for _, p := range All() {
		t.Run(p, func(t *testing.T) {
			got := Resolve(p)
			assert.Equal(t, p, got)
			if p == "osx" {
				assert.Equal(t, "osx", Resolve("macos"))
			}
		})
	}
}

func TestAll(t *testing.T) {
	got := All()
	want := []string{"common", "linux", "osx", "windows", "android"}
	assert.Equal(t, want, got)
	assert.Len(t, got, 5)
}

func TestAllReturnsDistinctCopy(t *testing.T) {
	a := All()
	b := All()
	a[0] = "mutated"
	assert.Equal(t, "common", b[0], "mutating one All() result should not affect another")
}

func TestAllNoDuplicates(t *testing.T) {
	got := All()
	seen := make(map[string]bool, len(got))
	for _, p := range got {
		assert.False(t, seen[p], "All should not contain duplicates: %q", p)
		seen[p] = true
	}
}
