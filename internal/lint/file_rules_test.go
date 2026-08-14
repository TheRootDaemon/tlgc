package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckLeadingWhitespace(t *testing.T) {
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 0,
		rawLine:    "",
	}
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 1,
		rawLine:    "# App",
		content:    "App",
	}
	leadingSpace := parsedLine{
		kind:       kindTitle,
		lineNumber: 1,
		rawLine:    " # App",
		content:    "App",
	}
	leadingTab := parsedLine{
		kind:       kindTitle,
		lineNumber: 1,
		rawLine:    "\t# App",
		content:    "App",
	}

	tests := []struct {
		name     string
		lines    []parsedLine
		wantCode string
	}{
		{name: "clean page passes", lines: []parsedLine{title}, wantCode: ""},
		{name: "leading space fails", lines: []parsedLine{leadingSpace}, wantCode: "TLDR001"},
		{name: "leading tab fails", lines: []parsedLine{leadingTab}, wantCode: "TLDR001"},
		{name: "leading blank line fails", lines: []parsedLine{blank, title}, wantCode: "TLDR001"},
		{name: "empty page passes", lines: nil, wantCode: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkLeadingWhitespace(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCode, errorCode(r))
		})
	}
}

func TestCheckSpaceAfterPrefix(t *testing.T) {
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 0,
		rawLine:    "# App",
		content:    "App",
	}
	noSpaceTitle := parsedLine{
		kind:       kindTitle,
		lineNumber: 1,
		rawLine:    "#App",
		content:    "App",
	}
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> Description.",
		content:    "Description.",
	}
	noSpaceDescription := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    ">Description.",
		content:    "Description.",
	}
	exampleDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 4,
		rawLine:    "- Example:",
		content:    "Example:",
	}
	noSpaceExampleDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 5,
		rawLine:    "-Example:",
		content:    "Example:",
	}
	command := parsedLine{
		kind:               kindCommand,
		lineNumber:         6,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{
			name:      "clean page passes",
			lines:     []parsedLine{title, description, exampleDescription, command},
			wantCodes: nil,
		},
		{
			name:      "title without space fails",
			lines:     []parsedLine{noSpaceTitle},
			wantCodes: []string{"TLDR002"},
		},
		{
			name:      "description without space fails",
			lines:     []parsedLine{noSpaceDescription},
			wantCodes: []string{"TLDR002"},
		},
		{
			name:      "example description without space fails",
			lines:     []parsedLine{noSpaceExampleDescription},
			wantCodes: []string{"TLDR002"},
		},
		{
			name:      "multiple violations fail",
			lines:     []parsedLine{noSpaceTitle, noSpaceDescription, noSpaceExampleDescription},
			wantCodes: []string{"TLDR002", "TLDR002", "TLDR002"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkSpaceAfterPrefix(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckNoTrailingWhitespaceAtEOF(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantCodes []string
		wantLine  int
	}{
		{name: "single trailing newline passes", raw: "# App\n", wantCodes: nil},
		{name: "no trailing newline passes", raw: "# App", wantCodes: nil},
		{name: "trailing space on last line passes", raw: "# App ", wantCodes: nil},
		{name: "spaces only no newline passes", raw: "# App   ", wantCodes: nil},
		{name: "trailing space then newline passes", raw: "# App \n", wantCodes: nil},
		{name: "one blank line at EOF fails", raw: "# App\n\n", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "multiple blank lines at EOF fail", raw: "# App\n\n\n\n", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "blank line with space at EOF fails", raw: "# App\n \n", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "blank line then trailing space fails", raw: "# App\n\n ", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "whitespace after final newline fails", raw: "# App\n ", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "tab after final newline fails", raw: "# App\n\t", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "spaces after final newline fail", raw: "# App\n  ", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "crlf blank line at EOF fails", raw: "# App\n\r\n", wantCodes: []string{"TLDR008"}, wantLine: 2},
		{name: "crlf then newline at EOF fails", raw: "# App\r\n\n", wantCodes: []string{"TLDR008"}, wantLine: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			lines := parseLines(tt.raw)
			checkNoTrailingWhitespaceAtEOF(&parsedPage{rawContent: tt.raw, lines: lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
			if len(r.Errors) > 0 {
				require.Equal(t, tt.wantLine, r.Errors[0].Line)
			}
		})
	}
}

func TestCheckEndsWithNewline(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantCode string
	}{
		{name: "ends with newline passes", raw: "# App\n", wantCode: ""},
		{name: "does not end with newline fails", raw: "# App", wantCode: "TLDR009"},
		{name: "empty page fails", raw: "", wantCode: "TLDR009"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkEndsWithNewline(&parsedPage{rawContent: tt.raw}, r)
			require.Equal(t, tt.wantCode, errorCode(r))
		})
	}
}

func TestCheckUnixLineEndings(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantCodes []string
		wantLines []int
	}{
		{
			name:      "unix line endings pass",
			raw:       "# App\n",
			wantCodes: nil,
		},
		{
			name:      "carriage return fails",
			raw:       "# App\r\n",
			wantCodes: []string{"TLDR010"},
			wantLines: []int{1},
		},
		{
			name:      "crlf on multiple lines fails per line",
			raw:       "# App\r\n> Description.\r\n",
			wantCodes: []string{"TLDR010", "TLDR010"},
			wantLines: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkUnixLineEndings(&parsedPage{rawContent: tt.raw, lines: parseLines(tt.raw)}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
			require.Equal(t, tt.wantLines, errorLines(r))
		})
	}
}

func TestCheckConsecutiveBlankLines(t *testing.T) {
	title := parsedLine{kind: kindTitle, lineNumber: 0, rawLine: "# App", content: "App"}
	blank1 := parsedLine{kind: kindBlank, lineNumber: 1, rawLine: ""}
	blank2 := parsedLine{kind: kindBlank, lineNumber: 2, rawLine: ""}
	blank3 := parsedLine{kind: kindBlank, lineNumber: 3, rawLine: ""}
	description := parsedLine{kind: kindDescription, lineNumber: 4, rawLine: "> Description.", content: "Description."}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{
			name:      "single blank lines pass",
			lines:     []parsedLine{title, blank1, description},
			wantCodes: nil,
		},
		{
			name:      "two consecutive blank lines fail",
			lines:     []parsedLine{title, blank1, blank2, description},
			wantCodes: []string{"TLDR011"},
		},
		{
			name:      "three consecutive blank lines fail once",
			lines:     []parsedLine{title, blank1, blank2, blank3, description},
			wantCodes: []string{"TLDR011"},
		},
		{
			name:      "page ending in blank run passes",
			lines:     []parsedLine{title, blank1, blank2, blank3},
			wantCodes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkConsecutiveBlankLines(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckNoTabs(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantCode string
	}{
		{name: "no tabs passes", raw: "# App\n", wantCode: ""},
		{name: "tab fails", raw: "# App\t\n", wantCode: "TLDR012"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkNoTabs(&parsedPage{rawContent: tt.raw, lines: parseLines(tt.raw)}, r)
			require.Equal(t, tt.wantCode, errorCode(r))
		})
	}
}

func TestCheckTrailingWhitespace(t *testing.T) {
	clean := parsedLine{kind: kindTitle, lineNumber: 0, rawLine: "# App", content: "App"}
	trailingSpace := parsedLine{kind: kindDescription, lineNumber: 1, rawLine: "> Description.  ", content: "Description."}
	trailingTab := parsedLine{kind: kindCommand, lineNumber: 2, rawLine: "`ls`\t", content: "ls", hasClosingBacktick: true}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{name: "clean lines pass", lines: []parsedLine{clean}, wantCodes: nil},
		{name: "trailing space fails", lines: []parsedLine{trailingSpace}, wantCodes: []string{"TLDR014"}},
		{name: "trailing tab fails", lines: []parsedLine{trailingTab}, wantCodes: []string{"TLDR014"}},
		{
			name:      "multiple trailing whitespace lines fail",
			lines:     []parsedLine{trailingSpace, trailingTab},
			wantCodes: []string{"TLDR014", "TLDR014"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkTrailingWhitespace(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}
