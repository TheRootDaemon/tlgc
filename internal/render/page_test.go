package render

import (
	"testing"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    *Page
	}{
		{
			name:    "full page with title description URL and examples",
			content: "# tar\n\n> archive utility.\n> More information: <https://example.org/tar>.\n\n- create an archive:\n\n`tar cf {{archive.tar}} {{files}}`\n\n- extract an archive:\n\n`tar xf {{archive.tar}}`\n",
			want: &Page{
				Title:       "tar",
				Description: []string{"archive utility."},
				URL:         "https://example.org/tar",
				Examples: []Example{
					{Description: "create an archive", Command: "tar cf {{archive.tar}} {{files}}"},
					{Description: "extract an archive", Command: "tar xf {{archive.tar}}"},
				},
			},
		},
		{
			name:    "title only",
			content: "# title-only\n",
			want: &Page{
				Title: "title-only",
			},
		},
		{
			name:    "description without URL",
			content: "# desc-test\n\n> just a description.\n",
			want: &Page{
				Title:       "desc-test",
				Description: []string{"just a description."},
			},
		},
		{
			name:    "multiple description lines merged with space",
			content: "# multi-desc\n\n> line one.\n> line two.\n> line three.\n",
			want: &Page{
				Title:       "multi-desc",
				Description: []string{"line one.", "line two.", "line three."},
			},
		},
		{
			name:    "example description without trailing colon",
			content: "# no-colon\n\n- description without colon\n\n`echo hello`\n",
			want: &Page{
				Title: "no-colon",
				Examples: []Example{
					{Description: "description without colon", Command: "echo hello"},
				},
			},
		},
		{
			name:    "command with placeholders",
			content: "# placeholders\n\n- example:\n\n`tar cf {{archive.tar}} {{files}}`\n",
			want: &Page{
				Title: "placeholders",
				Examples: []Example{
					{Description: "example", Command: "tar cf {{archive.tar}} {{files}}"},
				},
			},
		},
		{
			name:    "command without placeholder",
			content: "# plain\n\n- example:\n\n`ls -la`\n",
			want: &Page{
				Title: "plain",
				Examples: []Example{
					{Description: "example", Command: "ls -la"},
				},
			},
		},
		{
			name:    "URL extracted from More information line",
			content: "# with-url\n\n> A description.\n> More information: <https://example.org/custom>.\n",
			want: &Page{
				Title:       "with-url",
				Description: []string{"A description."},
				URL:         "https://example.org/custom",
			},
		},
		{
			name:    "empty content",
			content: "",
			want:    &Page{},
		},
		{
			name:    "example without command",
			content: "# no-command\n\n- description only:\n",
			want: &Page{
				Title: "no-command",
				Examples: []Example{
					{Description: "description only"},
				},
			},
		},
		{
			name:    "command without preceding description",
			content: "# bare-command\n\n`echo hello`\n",
			want: &Page{
				Title: "bare-command",
				Examples: []Example{
					{Command: "echo hello"},
				},
			},
		},
		{
			name:    "real-world content from fixture",
			content: "# test page\n\n> This is a test page.\n> More information: <https://example.org>.\n\n- This is a description of a `command` example:\n\n`command --opt1 --opt2 {{placeholder}}`\n\n- Another one:\n\n`command --opt1 {{placeholder1 placeholder2 ...}}`\n",
			want: &Page{
				Title:       "test page",
				Description: []string{"This is a test page."},
				URL:         "https://example.org",
				Examples: []Example{
					{Description: "This is a description of a `command` example", Command: "command --opt1 --opt2 {{placeholder}}"},
					{Description: "Another one", Command: "command --opt1 {{placeholder1 placeholder2 ...}}"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.content)
			require.NotNil(t, got)
			assert.Equal(t, tt.want.Title, got.Title, "Title mismatch")
			assert.Equal(t, tt.want.Description, got.Description, "Description mismatch")
			assert.Equal(t, tt.want.URL, got.URL, "URL mismatch")
			assert.Equal(t, tt.want.Examples, got.Examples, "Examples mismatch")
		})
	}
}

func TestParseCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want []Segment
	}{
		{
			name: "no placeholders",
			raw:  "tar cf archive.tar",
			want: []Segment{
				{Kind: Text, Text: "tar cf archive.tar"},
			},
		},
		{
			name: "single placeholder",
			raw:  "tar cf {{archive.tar}}",
			want: []Segment{
				{Kind: Text, Text: "tar cf "},
				{Kind: Placeholder, Text: "archive.tar"},
			},
		},
		{
			name: "multiple placeholders",
			raw:  "mv {{source}} {{destination}}",
			want: []Segment{
				{Kind: Text, Text: "mv "},
				{Kind: Placeholder, Text: "source"},
				{Kind: Text, Text: " "},
				{Kind: Placeholder, Text: "destination"},
			},
		},
		{
			name: "option placeholder short long",
			raw:  "cmd {{[-s|--long]}}",
			want: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-s", Long: "--long"},
			},
		},
		{
			name: "option placeholder long short reversed",
			raw:  "cmd {{[--long|-s]}}",
			want: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-s", Long: "--long"},
			},
		},
		{
			name: "mixed text placeholder and option",
			raw:  "cmd {{[-s|--long]}} {{file}}",
			want: []Segment{
				{Kind: Text, Text: "cmd "},
				{Kind: Option, Short: "-s", Long: "--long"},
				{Kind: Text, Text: " "},
				{Kind: Placeholder, Text: "file"},
			},
		},
		{
			name: "empty string",
			raw:  "",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCommand(tt.raw)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDisplayText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		segment Segment
		style   config.OptionStyle
		want    string
	}{
		{
			name:    "SegmentText with short style",
			segment: Segment{Kind: Text, Text: "hello"},
			style:   config.OptionStyleShort,
			want:    "hello",
		},
		{
			name:    "SegmentText with long style",
			segment: Segment{Kind: Text, Text: "hello"},
			style:   config.OptionStyleLong,
			want:    "hello",
		},
		{
			name:    "SegmentText with both style",
			segment: Segment{Kind: Text, Text: "hello"},
			style:   config.OptionStyleCombined,
			want:    "hello",
		},
		{
			name:    "SegmentPlaceholder with short style",
			segment: Segment{Kind: Placeholder, Text: "file"},
			style:   config.OptionStyleShort,
			want:    "file",
		},
		{
			name:    "SegmentPlaceholder with long style",
			segment: Segment{Kind: Placeholder, Text: "file"},
			style:   config.OptionStyleLong,
			want:    "file",
		},
		{
			name:    "SegmentPlaceholder with both style",
			segment: Segment{Kind: Placeholder, Text: "file"},
			style:   config.OptionStyleCombined,
			want:    "file",
		},
		{
			name:    "SegmentOption with short style",
			segment: Segment{Kind: Option, Short: "-s", Long: "--long"},
			style:   config.OptionStyleShort,
			want:    "-s",
		},
		{
			name:    "SegmentOption with long style",
			segment: Segment{Kind: Option, Short: "-s", Long: "--long"},
			style:   config.OptionStyleLong,
			want:    "--long",
		},
		{
			name:    "SegmentOption with both style",
			segment: Segment{Kind: Option, Short: "-s", Long: "--long"},
			style:   config.OptionStyleCombined,
			want:    "[-s|--long]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.segment.DisplayText(tt.style)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "valid URL with trailing text", body: "More information: <https://example.org>.", want: "https://example.org"},
		{name: "valid URL without trailing text", body: "More information: <https://example.org>", want: "https://example.org"},
		{name: "no URL marker", body: "Some description text.", want: ""},
		{name: "missing opening bracket", body: "More information: https://example.org>.", want: ""},
		{name: "missing closing bracket", body: "More information: <https://example.org", want: ""},
		{name: "empty brackets", body: "More information: <>.", want: ""},
		{name: "empty body", body: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseURL(tt.body)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInnerPlaceholders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		inner string
		want  Segment
	}{
		{name: "option short long", inner: "[-s|--long]", want: Segment{Kind: Option, Short: "-s", Long: "--long"}},
		{name: "option long short reversed", inner: "[--long|-s]", want: Segment{Kind: Option, Short: "-s", Long: "--long"}},
		{name: "plain placeholder", inner: "file", want: Segment{Kind: Placeholder, Text: "file"}},
		{name: "path placeholder", inner: "path/to/file", want: Segment{Kind: Placeholder, Text: "path/to/file"}},
		{name: "empty inner", inner: "", want: Segment{Kind: Placeholder, Text: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseInnerPlaceholders(tt.inner)
			assert.Equal(t, tt.want, got)
		})
	}
}
