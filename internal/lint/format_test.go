package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "empty",
			content: "",
			want:    "",
		},
		{
			name:    "whitespace_only",
			content: " \n\t\n",
			want:    "",
		},
		{
			name: "full_page",
			content: `# test

> Test description.
> More information: <https://example.com>.

- Example:

` + "`grep {{pattern}} {{file}}`" + `

- Second example:

` + "`echo hello`" + `
`,
			want: `# test

> Test description.
> More information: <https://example.com>.

- Example:

` + "`grep {{pattern}} {{file}}`" + `

- Second example:

` + "`echo hello`" + `
`,
		},
		{
			name: "lowercase_description_is_capitalized_and_punctuated",
			content: `# test

> hello world

- Example:

` + "`echo hello`" + `
`,
			want: `# test

> Hello world.

- Example:

` + "`echo hello`" + `
`,
		},
		{
			name: "lowercase_example_desc_colon_added",
			content: `# test

> Test description.

- example

` + "`echo hello`" + `
`,
			want: `# test

> Test description.

- Example:

` + "`echo hello`" + `
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormat_PreservesCommandText(t *testing.T) {
	t.Parallel()

	page := `# test

> Test description.
> More information: <https://example.com>.

- Example:

` + "`grep {{pattern}} {{file}}`" + `
`

	formatted := Format(page)
	require.Contains(t, formatted, "`grep {{pattern}} {{file}}`")
	require.NotContains(t, formatted, "undefined")
}

func TestFormat_PreservesCommandWithoutPlaceholders(t *testing.T) {
	t.Parallel()

	page := `# test

> Test description.
> More information: <https://example.com>.

- Example:

` + "`echo hello`" + `
`

	formatted := Format(page)
	require.Contains(t, formatted, "`echo hello`")
	require.NotContains(t, formatted, "undefined")
}

func TestFormatDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "hello world", want: "Hello world."},
		{in: "Hello world", want: "Hello world."},
		{in: "already punctuated.", want: "Already punctuated."},
		{in: "ends with colon:", want: "Ends with colon:."},
		{in: "ελληνικά", want: "Ελληνικά."},
		{in: "", want: "."},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, formatDescription(tt.in))
		})
	}
}

func TestFormatExampleDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "example", want: "Example:"},
		{in: "Example", want: "Example:"},
		{in: "already punctuated:", want: "Already punctuated:"},
		{in: "ends with period.", want: "Ends with period:"},
		{in: "", want: ":"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, formatExampleDescription(tt.in))
		})
	}
}

func TestUpperFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "", want: ""},
		{in: "hello", want: "Hello"},
		{in: "Hello", want: "Hello"},
		{in: "123abc", want: "123abc"},
		{in: "ελληνικά", want: "Ελληνικά"},
		{in: "hello WORLD", want: "Hello WORLD"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, upperFirst(tt.in))
		})
	}
}

func TestStripTrailing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		cut  string
		want string
	}{
		{in: "hello", cut: ".,;!?", want: "hello"},
		{in: "hello.", cut: ".,;!?", want: "hello"},
		{in: "hello:", cut: ".,;!?", want: "hello:"},
		{in: "hello...", cut: ".,;!?", want: "hello.."},
		{in: "", cut: ".,;!?", want: ""},
		{in: "hello:", cut: ".:,;", want: "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, stripTrailing(tt.in, tt.cut))
		})
	}
}

func TestAngleBracketedURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{in: "More information: <https://example.com>.", want: "<https://example.com>"},
		{in: "no brackets", want: ""},
		{in: "open < only", want: ""},
		{in: "a<b>c<d>", want: "<b>"},
		{in: "bare <>", want: "<>"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, angleBracketedURL(tt.in))
		})
	}
}
