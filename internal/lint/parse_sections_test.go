package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexOfTitle(t *testing.T) {
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 1,
		rawLine:    "# T",
		content:    "T",
	}
	secondTitle := parsedLine{
		kind:       kindTitle,
		lineNumber: 4,
		rawLine:    "# U",
		content:    "U",
	}
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 0,
		rawLine:    "> D",
		content:    "D",
	}
	command := parsedLine{
		kind:               kindCommand,
		lineNumber:         2,
		rawLine:            "`c`",
		content:            "c",
		hasClosingBacktick: true,
	}
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 3,
		rawLine:    "",
	}

	tests := []struct {
		name  string
		lines []parsedLine
		want  int
	}{
		{
			name:  "nil lines",
			lines: nil,
			want:  -1,
		},
		{
			name:  "no title",
			lines: []parsedLine{description, command},
			want:  -1,
		},
		{
			name:  "title at the start",
			lines: []parsedLine{title, description},
			want:  0,
		},
		{
			name:  "title in the middle",
			lines: []parsedLine{description, command, title},
			want:  2,
		},
		{
			name:  "two titles: the first wins",
			lines: []parsedLine{title, blank, secondTitle},
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := indexOfTitle(tt.lines)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNextContentIndex(t *testing.T) {
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 0,
		rawLine:    "# T",
		content:    "T",
	}
	text := parsedLine{
		kind:       kindText,
		lineNumber: 2,
		rawLine:    "x",
		content:    "x",
	}
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 1,
		rawLine:    "",
	}

	tests := []struct {
		name  string
		lines []parsedLine
		i     int
		want  int
	}{
		{
			name:  "nil lines",
			lines: nil,
			i:     0,
			want:  0,
		},
		{
			name:  "starts on content returns the same index",
			lines: []parsedLine{title, blank},
			i:     0,
			want:  0,
		},
		{
			name:  "skips leading blanks",
			lines: []parsedLine{blank, blank, title},
			i:     0,
			want:  2,
		},
		{
			name:  "i beyond len returns i unchanged",
			lines: []parsedLine{title},
			i:     5,
			want:  5,
		},
		{
			name:  "all-blank suffix returns len",
			lines: []parsedLine{title, blank, blank},
			i:     1,
			want:  3,
		},
		{
			name:  "stops at any non-blank, including text",
			lines: []parsedLine{blank, blank, text},
			i:     0,
			want:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nextContentIndex(tt.lines, tt.i)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCollectDescriptions(t *testing.T) {
	description_A := parsedLine{
		kind:       kindDescription,
		lineNumber: 0,
		rawLine:    "> A",
		content:    "A",
	}
	description_B := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> B",
		content:    "B",
	}
	link := parsedLine{
		kind:       kindDescription,
		lineNumber: 0,
		rawLine:    "> More information: https://x",
		content:    "More information: https://x",
	}
	noSpace := parsedLine{
		kind:       kindDescription,
		lineNumber: 0,
		rawLine:    "> More information:",
		content:    "More information:",
	}
	lower := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> more information: x",
		content:    "more information: x",
	}
	command := parsedLine{
		kind:       kindCommand,
		lineNumber: 2,
		rawLine:    "`c`",
		content:    "c", hasClosingBacktick: true,
	}
	text := parsedLine{
		kind:       kindText,
		lineNumber: 1,
		rawLine:    "x",
		content:    "x",
	}
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 1,
		rawLine:    "",
	}

	tests := []struct {
		name            string
		lines           []parsedLine
		i               int
		wantDescription []parsedLine
		wantInfo        []parsedLine
		wantNext        int
	}{
		{
			name:            "empty input yields non-nil empty results",
			lines:           nil,
			i:               0,
			wantDescription: []parsedLine{}, // make-backed: non-nil even when empty
			wantInfo:        []parsedLine{},
			wantNext:        0,
		},
		{
			name:            "consecutive run",
			lines:           []parsedLine{description_A, description_B, command},
			i:               0,
			wantDescription: []parsedLine{description_A, description_B},
			wantInfo:        []parsedLine{},
			wantNext:        2,
		},
		{
			name:            "stops at the first non-description",
			lines:           []parsedLine{description_A, text, description_B},
			i:               0,
			wantDescription: []parsedLine{description_A},
			wantInfo:        []parsedLine{},
			wantNext:        1,
		},
		{
			name:            "a blank line ends the run",
			lines:           []parsedLine{description_A, blank, description_B},
			i:               0,
			wantDescription: []parsedLine{description_A},
			wantInfo:        []parsedLine{},
			wantNext:        1,
		},
		{
			name:            "info link is collected and detected",
			lines:           []parsedLine{link, description_A},
			i:               0,
			wantDescription: []parsedLine{link, description_A},
			wantInfo:        []parsedLine{link},
			wantNext:        2,
		},
		{
			name:            "exact prefix is required",
			lines:           []parsedLine{noSpace, lower},
			i:               0,
			wantDescription: []parsedLine{noSpace, lower},
			wantInfo:        []parsedLine{},
			wantNext:        2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			description, info, next := collectDescriptions(tt.lines, tt.i)
			require.Equal(t, tt.wantDescription, description)
			require.Equal(t, tt.wantInfo, info)
			require.Equal(t, tt.wantNext, next)
		})
	}
}

func TestCollectExampleSections(t *testing.T) {
	exampleDescription_1 := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 0,
		rawLine:    "- A",
		content:    "A",
	}
	command_1 := parsedLine{
		kind:               kindCommand,
		lineNumber:         1,
		rawLine:            "`a`",
		content:            "a",
		hasClosingBacktick: true,
	}
	blank_A := parsedLine{
		kind:       kindBlank,
		lineNumber: 2,
		rawLine:    "",
	}
	blank_B := parsedLine{
		kind:       kindBlank,
		lineNumber: 3,
		rawLine:    "",
	}
	exampleDescription_2 := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 4,
		rawLine:    "- B",
		content:    "B",
	}
	command_2 := parsedLine{
		kind:               kindCommand,
		lineNumber:         5,
		rawLine:            "`b`",
		content:            "b",
		hasClosingBacktick: true,
	}
	text := parsedLine{
		kind:       kindText,
		lineNumber: 6,
		rawLine:    "stray",
		content:    "stray",
	}
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 7,
		rawLine:    "> D",
		content:    "D",
	}

	tests := []struct {
		name  string
		lines []parsedLine
		i     int
		want  []commandSection
	}{
		{
			name:  "no examples returns nil",
			lines: nil,
			i:     0,
			want:  nil, // bare var: nil when empty
		},
		{
			name:  "leading blank is skipped",
			lines: []parsedLine{blank_A, exampleDescription_1, command_1},
			i:     0,
			want: []commandSection{
				{description: "A", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
			},
		},
		{
			name:  "blank gap between sections",
			lines: []parsedLine{exampleDescription_1, command_1, blank_A, blank_B, exampleDescription_2, command_2},
			i:     0,
			want: []commandSection{
				{description: "A", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
				{description: "B", descriptionLineNumber: 4, commands: []parsedLine{command_2}},
			},
		},
		{
			name:  "stray text lines are skipped",
			lines: []parsedLine{exampleDescription_1, command_1, text, exampleDescription_2, command_2},
			i:     0,
			want: []commandSection{
				{description: "A", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
				{description: "B", descriptionLineNumber: 4, commands: []parsedLine{command_2}},
			},
		},
		{
			name:  "description after examples is stray",
			lines: []parsedLine{exampleDescription_1, command_1, description},
			i:     0,
			want: []commandSection{
				{description: "A", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
			},
		},
		{
			name:  "consecutive example descriptions",
			lines: []parsedLine{exampleDescription_1, exampleDescription_2, command_1},
			i:     0,
			want: []commandSection{
				{description: "A", descriptionLineNumber: 0}, // no commands
				{description: "B", descriptionLineNumber: 4, commands: []parsedLine{command_1}},
			},
		},
		{
			name:  "trailing blanks after the last command",
			lines: []parsedLine{exampleDescription_1, command_1, blank_A, blank_B},
			i:     0,
			want: []commandSection{
				{description: "A", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectExampleSections(tt.lines, tt.i)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBuildExampleSection(t *testing.T) {
	description := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 0,
		rawLine:    "- Copy files",
		content:    "Copy files",
	}
	command_1 := parsedLine{
		kind:               kindCommand,
		lineNumber:         1,
		rawLine:            "`cp a b`",
		content:            "cp a b",
		hasClosingBacktick: true,
	}
	command_2 := parsedLine{
		kind:               kindCommand,
		lineNumber:         3,
		rawLine:            "`scp a b`",
		content:            "scp a b",
		hasClosingBacktick: true,
	}
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 2,
		rawLine:    "",
	}
	text := parsedLine{
		kind:       kindText,
		lineNumber: 4,
		rawLine:    "stray",
		content:    "stray",
	}

	tests := []struct {
		name     string
		lines    []parsedLine
		i        int
		want     commandSection
		wantNext int
	}{
		{
			name:     "description with one command",
			lines:    []parsedLine{description, command_1},
			i:        0,
			want:     commandSection{description: "Copy files", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
			wantNext: 2,
		},
		{
			name:     "blank between description and command is tolerated",
			lines:    []parsedLine{description, blank, command_1},
			i:        0,
			want:     commandSection{description: "Copy files", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
			wantNext: 3,
		},
		{
			name:     "blank between commands is tolerated",
			lines:    []parsedLine{description, command_1, blank, command_2},
			i:        0,
			want:     commandSection{description: "Copy files", descriptionLineNumber: 0, commands: []parsedLine{command_1, command_2}},
			wantNext: 4,
		},
		{
			name:     "description with no commands leaves commands nil",
			lines:    []parsedLine{description, text},
			i:        0,
			want:     commandSection{description: "Copy files", descriptionLineNumber: 0}, // commands: nil — append never ran
			wantNext: 1,
		},
		{
			name:     "stops at the next non-command line without consuming it",
			lines:    []parsedLine{description, command_1, text},
			i:        0,
			want:     commandSection{description: "Copy files", descriptionLineNumber: 0, commands: []parsedLine{command_1}},
			wantNext: 2, // next points at text; the caller decides what happens to it
		},
		{
			name: "line numbers come from the parsed lines themselves",
			lines: []parsedLine{
				{kind: kindExampleDesc, lineNumber: 7, rawLine: "- X", content: "X"},
				{kind: kindCommand, lineNumber: 9, rawLine: "`y`", content: "y", hasClosingBacktick: true},
			},
			i: 0,
			want: commandSection{
				description:           "X",
				descriptionLineNumber: 7,
				commands:              []parsedLine{{kind: kindCommand, lineNumber: 9, rawLine: "`y`", content: "y", hasClosingBacktick: true}},
			},
			wantNext: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, next := buildExampleSection(tt.lines, tt.i)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantNext, next)
		})
	}
}

func TestIsInfoLink(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "exact prefix",
			content: "More information: https://example.com",
			want:    true,
		},
		{
			name:    "missing trailing space",
			content: "More information:",
			want:    false,
		},
		{
			name:    "lowercase prefix",
			content: "more information: https://example.com",
			want:    false,
		},
		{
			name:    "empty content",
			content: "",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isInfoLink(tt.content)
			require.Equal(t, tt.want, got)
		})
	}
}
