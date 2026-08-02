package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantCode string
	}{
		{
			name:     "md extension passes",
			filename: "tldr.md",
			wantCode: "",
		},
		{
			name:     "path with md extension passes",
			filename: "pages/common/tldr.md",
			wantCode: "",
		},
		{
			name:     "dotfile with md extension passes",
			filename: ".md",
			wantCode: "",
		},
		{
			name:     "non-md extension fails",
			filename: "tldr.txt",
			wantCode: "TLDR107",
		},
		{
			name:     "no extension fails",
			filename: "tldr",
			wantCode: "TLDR107",
		},
		{
			name:     "wrong final extension fails",
			filename: "tldr.md.txt",
			wantCode: "TLDR107",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				r := &Result{}
				checkFileExtension(tt.filename, r)
				require.Equal(t, tt.wantCode, errorCode(r))
			},
		)
	}
}

func TestCheckFilenameWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantCode string
	}{
		{
			name:     "clean name passes",
			filename: "tldr.md",
			wantCode: "",
		},
		{
			name:     "space in directory only passes",
			filename: "my dir/tldr.md",
			wantCode: "",
		},
		{
			name:     "space in name fails",
			filename: "tldr page.md",
			wantCode: "TLDR108",
		},
		{
			name:     "tab in name fails",
			filename: "tldr\tpage.md",
			wantCode: "TLDR108",
		},
		{
			name:     "leading space fails",
			filename: " tldr.md",
			wantCode: "TLDR108",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				r := &Result{}
				checkFilenameWhitespace(tt.filename, r)
				require.Equal(t, tt.wantCode, errorCode(r))
			},
		)
	}
}

func TestCheckFilenameLowercase(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantCode string
	}{
		{
			name:     "lowercase name passes",
			filename: "tldr.md",
			wantCode: "",
		},
		{
			name:     "path with lowercase name passes",
			filename: "pages/tldr.md",
			wantCode: "",
		},
		{
			name:     "digits pass",
			filename: "tldr1.md",
			wantCode: "",
		},
		{
			name:     "uppercase first letter fails",
			filename: "Tldr.md",
			wantCode: "TLDR109",
		},
		{
			name:     "all uppercase fails",
			filename: "TLDR.MD",
			wantCode: "TLDR109",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkFilenameLowercase(tt.filename, r)
			require.Equal(t, tt.wantCode, errorCode(r))
		})
	}
}

func TestCheckForbiddenFilenameCharacters(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantCode string
	}{
		{
			name:     "clean name passes",
			filename: "tldr.md",
			wantCode: "",
		},
		{
			name:     "directory separator in path passes",
			filename: "dir/tldr.md",
			wantCode: "",
		},
		{
			name:     "angle bracket fails",
			filename: "tldr<page.md",
			wantCode: "TLDR111",
		},
		{
			name:     "question mark fails",
			filename: "tldr?page.md",
			wantCode: "TLDR111",
		},
		{
			name:     "colon fails",
			filename: "tldr:page.md",
			wantCode: "TLDR111",
		},
		{
			name:     "backslash fails",
			filename: "tldr\\page.md",
			wantCode: "TLDR111",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				r := &Result{}
				checkForbiddenFilenameCharacters(tt.filename, r)
				require.Equal(t, tt.wantCode, errorCode(r))
			},
		)
	}
}

// errorCode returns the code of the first reported error,
// or "" if the result is clean.
// Rules report at most one error, always at line 0.
func errorCode(r *Result) string {
	if len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].Code
}
