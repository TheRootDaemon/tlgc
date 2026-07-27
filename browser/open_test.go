package browser

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type call struct {
	name string
	args []string
}

func TestOpenOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("openOnLinux only runs on Linux")
	}
	if runningOnWSL() {
		t.Skip("WSL routes to explorer.exe, tested separately")
	}

	tests := []struct {
		name      string
		display   string
		wayland   string
		wantErr   bool
		wantCalls []call
	}{
		{
			name:    "with display",
			display: ":0",
			wantCalls: []call{
				{name: "xdg-open", args: []string{"https://example.com"}},
			},
		},
		{
			name:    "no display",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DISPLAY", tt.display)
			t.Setenv("WAYLAND_DISPLAY", tt.wayland)

			var calls []call
			mockBrowse(t, func(name string, args ...string) error {
				calls = append(calls, call{name, args})
				return nil
			})

			err := openOnLinux("https://example.com")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "no display server detected")
				assert.Empty(t, calls, "browse should not be called on error")
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestOpenOnWSL(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("openOnWSL only runs on Linux")
	}
	if !runningOnWSL() {
		t.Skip("not running on WSL")
	}

	var calls []call
	mockBrowse(t, func(name string, args ...string) error {
		calls = append(calls, call{name, args})
		return nil
	})

	err := openOnWSL("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, []call{
		{name: "explorer.exe", args: []string{"https://example.com"}},
	}, calls)
}

func TestOpen(t *testing.T) {
	t.Setenv("DISPLAY", "")
	t.Setenv("WAYLAND_DISPLAY", "")

	switch runtime.GOOS {
	case "linux":
		if runningOnWSL() {
			var calls []call
			mockBrowse(t, func(name string, args ...string) error {
				calls = append(calls, call{name, args})
				return nil
			})

			err := Open("https://example.com")
			assert.NoError(t, err)
			assert.Equal(
				t,
				[]call{
					{
						name: "explorer.exe",
						args: []string{"https://example.com"},
					},
				}, calls,
			)
		} else {
			err := Open("https://example.com")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "no display server detected")
		}
	case "darwin":
		var calls []call
		mockBrowse(t, func(name string, args ...string) error {
			calls = append(calls, call{name, args})
			return nil
		})

		err := Open("https://example.com")
		assert.NoError(t, err)
		assert.Equal(
			t,
			[]call{{name: "open", args: []string{"https://example.com"}}},
			calls,
		)
	case "windows":
		var calls []call
		mockBrowse(t, func(name string, args ...string) error {
			calls = append(calls, call{name, args})
			return nil
		})

		err := Open("https://example.com")
		assert.NoError(t, err)
		assert.Equal(
			t,
			[]call{{name: "explorer.exe", args: []string{"https://example.com"}}},
			calls,
		)
	}
}

func TestHasDisplay(t *testing.T) {
	tests := []struct {
		name    string
		display string
		wayland string
		want    bool
	}{
		{name: "neither set", display: "", wayland: "", want: false},
		{name: "x11 only", display: ":0", wayland: "", want: true},
		{name: "wayland only", display: "", wayland: "wayland-0", want: true},
		{name: "both set", display: ":0", wayland: "wayland-0", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DISPLAY", tt.display)
			t.Setenv("WAYLAND_DISPLAY", tt.wayland)
			assert.Equal(t, tt.want, hasDisplay())
		})
	}
}

// mockBrowse replaces the package-level browse function for the duration of t.
func mockBrowse(t *testing.T, fn func(string, ...string) error) {
	t.Helper()
	oldBrowse := browse
	browse = fn
	t.Cleanup(func() { browse = oldBrowse })
}
