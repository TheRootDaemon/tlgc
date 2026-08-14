package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLines(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []parsedLine
	}{
		{
			name: "empty input returns nil",
			raw:  "",
			want: nil,
		},
		{
			name: "whitespace only returns nil",
			raw:  " \n\t ",
			want: nil,
		},
		{
			name: "single line",
			raw:  "> Hello",
			want: []parsedLine{
				{kind: kindDescription, lineNumber: 1, rawLine: "> Hello", content: "Hello"},
			},
		},
		{
			name: "consecutive lines are numbered in order",
			raw:  "# App\n> D\n`c`",
			want: []parsedLine{
				{kind: kindTitle, lineNumber: 1, rawLine: "# App", content: "App"},
				{kind: kindDescription, lineNumber: 2, rawLine: "> D", content: "D"},
				{kind: kindCommand, lineNumber: 3, rawLine: "`c`", content: "c", hasClosingBacktick: true},
			},
		},
		{
			name: "trailing newline yields a trailing blank line",
			raw:  "> A\n",
			want: []parsedLine{
				{kind: kindDescription, lineNumber: 1, rawLine: "> A", content: "A"},
				{kind: kindBlank, lineNumber: 2, rawLine: ""},
			},
		},
		{
			name: "crlf endings classify correctly but keep raw",
			raw:  "> A\r\n`B`\r\n",
			want: []parsedLine{
				{kind: kindDescription, lineNumber: 1, rawLine: "> A\r", content: "A"},
				{kind: kindCommand, lineNumber: 2, rawLine: "`B`\r", content: "B", hasClosingBacktick: true},
				{kind: kindBlank, lineNumber: 3, rawLine: ""},
			},
		},
		{
			name: "leading blank line is preserved",
			raw:  "\n# T",
			want: []parsedLine{
				{kind: kindBlank, lineNumber: 1, rawLine: ""},
				{kind: kindTitle, lineNumber: 2, rawLine: "# T", content: "T"},
			},
		},
		{
			name: "empty middle line is a blank",
			raw:  "# T\n\n> D",
			want: []parsedLine{
				{kind: kindTitle, lineNumber: 1, rawLine: "# T", content: "T"},
				{kind: kindBlank, lineNumber: 2, rawLine: ""},
				{kind: kindDescription, lineNumber: 3, rawLine: "> D", content: "D"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLines(tt.raw)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name       string
		lineNumber int
		rawLine    string
		want       parsedLine
	}{
		{
			name:       "empty line is blank",
			lineNumber: 3,
			rawLine:    "",
			want:       parsedLine{kind: kindBlank, lineNumber: 3, rawLine: ""},
		},
		{
			name:       "whitespace only line is blank",
			lineNumber: 0,
			rawLine:    " \t ",
			want:       parsedLine{kind: kindBlank, lineNumber: 0, rawLine: " \t "},
		},
		{
			name:       "title content is trimmed",
			lineNumber: 0,
			rawLine:    "#   App  ",
			want:       parsedLine{kind: kindTitle, lineNumber: 0, rawLine: "#   App  ", content: "App"},
		},
		{
			name:       "bare hash is a title with empty content",
			lineNumber: 0,
			rawLine:    "#",
			want:       parsedLine{kind: kindTitle, lineNumber: 0, rawLine: "#", content: ""},
		},
		{
			name:       "description strips marker and one space",
			lineNumber: 2,
			rawLine:    "> Brief description.",
			want:       parsedLine{kind: kindDescription, lineNumber: 2, rawLine: "> Brief description.", content: "Brief description."},
		},
		{
			name:       "description without marker space",
			lineNumber: 0,
			rawLine:    ">No space",
			want:       parsedLine{kind: kindDescription, lineNumber: 0, rawLine: ">No space", content: "No space"},
		},
		{
			name:       "command keeps content and closing flag",
			lineNumber: 4,
			rawLine:    "`cp file file.bak`",
			want:       parsedLine{kind: kindCommand, lineNumber: 4, rawLine: "`cp file file.bak`", content: "cp file file.bak", hasClosingBacktick: true},
		},
		{
			name:       "command with missing closing backtick",
			lineNumber: 0,
			rawLine:    "`ls -la",
			want:       parsedLine{kind: kindCommand, lineNumber: 0, rawLine: "`ls -la", content: "ls -la", hasClosingBacktick: false},
		},
		{
			name:       "single backtick is an empty unterminated command",
			lineNumber: 0,
			rawLine:    "`",
			want:       parsedLine{kind: kindCommand, lineNumber: 0, rawLine: "`", content: "", hasClosingBacktick: false},
		},
		{
			name:       "crlf ending is trimmed for classification but kept in rawLine",
			lineNumber: 1,
			rawLine:    "> Hello\r",
			want:       parsedLine{kind: kindDescription, lineNumber: 1, rawLine: "> Hello\r", content: "Hello"},
		},
		{
			name:       "leading whitespace defeats marker classification",
			lineNumber: 0,
			rawLine:    "  `ls`",
			want:       parsedLine{kind: kindText, lineNumber: 0, rawLine: "  `ls`", content: "  `ls`"},
		},
		{
			name:       "text line is right trimmed",
			lineNumber: 0,
			rawLine:    "free text  \t",
			want:       parsedLine{kind: kindText, lineNumber: 0, rawLine: "free text  \t", content: "free text"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLine(tt.lineNumber, tt.rawLine)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseCommandContent(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		wantContent string
		wantClosing bool
	}{
		{
			name:        "single backtick only",
			in:          "`",
			wantContent: "",
			wantClosing: false,
		},
		{
			name:        "closed command",
			in:          "`ls -la`",
			wantContent: "ls -la",
			wantClosing: true,
		},
		{
			name:        "missing closing backtick",
			in:          "`ls -la",
			wantContent: "ls -la",
			wantClosing: false,
		},
		{
			name:        "empty closed command",
			in:          "``",
			wantContent: "",
			wantClosing: true,
		},
		{
			name:        "stops at the first closing backtick",
			in:          "`echo `a`",
			wantContent: "echo ",
			wantClosing: true,
		},
		{
			name:        "no leading backtick returns the input untouched",
			in:          "ls -la",
			wantContent: "ls -la",
			wantClosing: false,
		},
		{
			name:        "content after a closing backtick is dropped",
			in:          "`a`b",
			wantContent: "a",
			wantClosing: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, closing := parseCommandContent(tt.in)
			require.Equal(t, tt.wantContent, content)
			require.Equal(t, tt.wantClosing, closing)
		})
	}
}

func TestStripMarker(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "marker plus one space",
			in:   "> foo",
			want: "foo",
		},
		{
			name: "marker without space",
			in:   ">foo",
			want: "foo",
		},
		{
			name: "only a single space is removed",
			in:   ">  foo",
			want: " foo",
		},
		{
			name: "trailing spaces and tabs are trimmed",
			in:   "> foo \t",
			want: "foo",
		},
		{
			name: "bare marker",
			in:   ">",
			want: "",
		},
		{
			name: "marker plus only a space",
			in:   "> ",
			want: "",
		},
		{
			name: "dash marker",
			in:   "- foo",
			want: "foo",
		},
		{
			name: "tab after marker is not a separator",
			in:   ">\tfoo",
			want: "\tfoo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripMarker(tt.in)
			require.Equal(t, tt.want, got)
		})
	}
}
