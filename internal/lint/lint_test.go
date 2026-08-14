package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// specFixture couples a failing spec page with the errors it must produce,
// mirroring the reference tldr-lint's tldr-lint.spec.js.
type specFixture struct {
	name   string
	want   []string
	count  int
	subset bool
}

func TestLintSpecsFailing(t *testing.T) {
	tests := []specFixture{
		{"failing/001.md", []string{"TLDR001"}, 1, false},
		{"failing/002.md", []string{"TLDR002"}, 3, false},
		{"failing/003.md", []string{"TLDR003"}, 1, false},
		{"failing/004.md", []string{"TLDR004", "TLDR014"}, 4, true},
		{"failing/005.md", []string{"TLDR005"}, 2, false},
		{"failing/006.md", []string{"TLDR006"}, 1, false},
		{"failing/007.md", []string{"TLDR007"}, 2, false},
		{"failing/008.md", []string{"TLDR008"}, 1, false},
		{"failing/009.md", []string{"TLDR009"}, 1, false},
		{"failing/010.md", []string{"TLDR010"}, 7, false},
		{"failing/011.md", []string{"TLDR011"}, 2, false},
		{"failing/012.md", []string{"TLDR012"}, 2, false},
		{"failing/013.md", []string{"TLDR013"}, 1, false},
		{"failing/014.md", []string{"TLDR014"}, 5, false},
		{"failing/015.md", []string{"TLDR015"}, 1, false},
		{"failing/016.md", []string{"TLDR016"}, 1, false},
		{"failing/017.md", []string{"TLDR017"}, 1, false},
		{"failing/018.md", []string{"TLDR018"}, 2, false},
		{"failing/019.md", []string{"TLDR019"}, 1, false},
		{"failing/020.md", []string{"TLDR020"}, 3, false},
		{"failing/021.md", []string{"TLDR021"}, 2, false},
		{"failing/101.md", []string{"TLDR101"}, 1, false},
		{"failing/102.md", []string{"TLDR102"}, 1, false},
		{"failing/103.md", []string{"TLDR103"}, 2, false},
		{"failing/104.md", []string{"TLDR104"}, 2, false},
		{"failing/105.md", []string{"TLDR105"}, 2, false},
		{"failing/106.md", []string{"TLDR106"}, 1, false},
		{"failing/107", []string{"TLDR107"}, 1, false},
		{"failing/108 .md", []string{"TLDR108"}, 1, false},
		{"failing/109A.md", []string{"TLDR109"}, 1, false},
		{"failing/110.md", []string{"TLDR110"}, 1, false},
		{"failing/112.md", []string{"TLDR112"}, 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("specs", "pages", tt.name))
			require.NoError(t, err)
			defer func() {
				_ = f.Close()
			}()

			r, err := Lint(f)
			require.NoError(t, err)
			assertSpecErrors(t, r, tt.want, tt.count, tt.subset)
		})
	}
}

func TestLintSpecsForbiddenFilenameCharacters(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("specs", "pages", "failing", "111.md"))
	require.NoError(t, err)

	for _, char := range `<>:"/\|?*` {
		t.Run("111"+string(char), func(t *testing.T) {
			r := lint("111"+string(char)+".md", content)
			assertSpecErrors(t, r, []string{"TLDR111"}, 1, false)
		})
	}
}

func TestLintSpecsPassing(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("specs", "pages", "passing"))
	require.NoError(t, err)
	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			f, err := os.Open(filepath.Join("specs", "pages", "passing", entry.Name()))
			require.NoError(t, err)
			defer func() {
				_ = f.Close()
			}()

			r, err := Lint(f)
			require.NoError(t, err)
			assertSpecErrors(t, r, nil, 0, false)
		})
	}
}

func TestLintSpecsIgnore(t *testing.T) {
	f, err := os.Open(filepath.Join("specs", "pages", "failing", "004.md"))
	require.NoError(t, err)
	defer func() {
		_ = f.Close()
	}()

	r, err := Lint(f, "TLDR014")
	require.NoError(t, err)
	assertSpecErrors(t, r, []string{"TLDR004"}, 2, false)
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "standard error",
			err: Error{
				Code:        "TLDR001",
				Line:        10,
				Description: "missing command description",
			},
			want: "TLDR001:10 missing command description",
		},
		{
			name: "zero line",
			err: Error{
				Code:        "TLDR002",
				Line:        0,
				Description: "invalid page format",
			},
			want: "TLDR002:0 invalid page format",
		},
		{
			name: "empty fields",
			err:  Error{},
			want: ":0 ",
		},
		{
			name: "description with punctuation",
			err: Error{
				Code:        "TLDR003",
				Line:        42,
				Description: "description must end with a period.",
			},
			want: "TLDR003:42 description must end with a period.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.String()
			require.Equal(t, tt.want, got)
		})
	}
}

// assertSpecErrors verifies that a lint result
// contains the expected number of errors and error codes.
//
// The total error count must match count exactly.
// Every code in want must be present in the result.
// When subset is false, want must also contain every distinct error code
// reported by the result;
// when true, additional error codes are allowed.
func assertSpecErrors(
	t *testing.T,
	r *Result,
	want []string,
	count int,
	subset bool,
) {
	t.Helper()
	require.Equal(t, count, len(r.Errors))

	seen := make(map[string]bool, len(r.Errors))
	for _, e := range r.Errors {
		seen[e.Code] = true
	}
	for _, code := range want {
		require.True(t, seen[code])
	}
	if !subset {
		require.Len(t, seen, len(want))
	}
}
