package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "typical path", path: "/path/to/pages.en/common/some-page.md", want: "some-page"},
		{name: "no extension", path: "/path/to/pages.en/common/some-page", want: "some-page"},
		{name: "empty path", path: "", want: ""},
		{name: "root path", path: "/", want: "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PageName(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPagePlatform(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "linux platform", path: "/path/to/pages.en/linux/some-page.md", want: "linux"},
		{name: "common platform", path: "/path/to/pages.en/common/some-page.md", want: "common"},
		{name: "no parent directory", path: "some-page.md", want: ""},
		{name: "empty path", path: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PagePlatform(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}
