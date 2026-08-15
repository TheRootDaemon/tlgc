package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckExampleDescriptionEndsWithColon(t *testing.T) {
	withColon := commandSection{
		description:           "List all files:",
		descriptionLineNumber: 1,
	}
	withoutColon := commandSection{
		description:           "List all files",
		descriptionLineNumber: 2,
	}
	empty := commandSection{
		description:           "",
		descriptionLineNumber: 3,
	}

	tests := []struct {
		name            string
		exampleSections []commandSection
		wantCodes       []string
	}{
		{name: "ends with colon passes", exampleSections: []commandSection{withColon}, wantCodes: nil},
		{name: "missing colon fails", exampleSections: []commandSection{withoutColon}, wantCodes: []string{"TLDR005"}},
		{name: "empty description passes", exampleSections: []commandSection{empty}, wantCodes: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkExampleDescriptionEndsWithColon(&parsedPage{exampleSections: tt.exampleSections}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckExampleDescriptionSurroundedByBlankLines(t *testing.T) {
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
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> Description.",
		content:    "Description.",
	}
	exampleDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 3,
		rawLine:    "- List all files:",
		content:    "List all files:",
	}
	exampleDescriptionNoBlankBefore := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 5,
		rawLine:    "- List all files:",
		content:    "List all files:",
	}
	command := parsedLine{
		kind:               kindCommand,
		lineNumber:         4,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	commandNoBlankBefore := parsedLine{
		kind:               kindCommand,
		lineNumber:         6,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	commandWithBlankBefore := parsedLine{
		kind:               kindCommand,
		lineNumber:         7,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name            string
		lines           []parsedLine
		exampleSections []commandSection
		wantCodes       []string
	}{
		{
			name:  "surrounded by blank lines passes",
			lines: []parsedLine{blank, exampleDescription, blank, command},
			exampleSections: []commandSection{
				{
					description:           "List all files:",
					descriptionLineNumber: 3,
					commands:              []parsedLine{command},
				},
			},
			wantCodes: nil,
		},
		{
			name: "no blank before description fails",
			lines: []parsedLine{
				title,
				description,
				exampleDescriptionNoBlankBefore,
				blank,
				commandWithBlankBefore,
			},
			exampleSections: []commandSection{
				{
					description:           "List all files:",
					descriptionLineNumber: 5,
					commands:              []parsedLine{commandWithBlankBefore},
				},
			},
			wantCodes: []string{"TLDR007"},
		},
		{
			name: "no blank before command fails",
			lines: []parsedLine{
				blank,
				exampleDescription,
				commandNoBlankBefore,
			},
			exampleSections: []commandSection{
				{
					description:           "List all files:",
					descriptionLineNumber: 3,
					commands:              []parsedLine{commandNoBlankBefore},
				},
			},
			wantCodes: []string{"TLDR007"},
		},
		{
			name:  "description without command passes",
			lines: []parsedLine{blank, exampleDescription},
			exampleSections: []commandSection{
				{
					description:           "List all files:",
					descriptionLineNumber: 3,
				},
			},
			wantCodes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkExampleDescriptionSurroundedByBlankLines(
				&parsedPage{
					lines:           tt.lines,
					exampleSections: tt.exampleSections,
				},
				r,
			)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckExampleDescriptionStartsWithCapital(t *testing.T) {
	capital := commandSection{
		description:           "List all files:",
		descriptionLineNumber: 1,
	}
	lowercase := commandSection{
		description:           "list all files:",
		descriptionLineNumber: 2,
	}
	placeholder := commandSection{
		description:           "[file] to copy:",
		descriptionLineNumber: 3,
	}
	empty := commandSection{
		description:           "",
		descriptionLineNumber: 4,
	}

	tests := []struct {
		name            string
		exampleSections []commandSection
		wantCodes       []string
	}{
		{name: "capital start passes", exampleSections: []commandSection{capital}, wantCodes: nil},
		{name: "lowercase start fails", exampleSections: []commandSection{lowercase}, wantCodes: []string{"TLDR015"}},
		{name: "placeholder start passes", exampleSections: []commandSection{placeholder}, wantCodes: nil},
		{name: "empty description passes", exampleSections: []commandSection{empty}, wantCodes: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkExampleDescriptionStartsWithCapital(&parsedPage{exampleSections: tt.exampleSections}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckMaximumExampleCount(t *testing.T) {
	section := commandSection{
		description:           "List all files:",
		descriptionLineNumber: 1,
	}

	tests := []struct {
		name            string
		exampleSections []commandSection
		wantCode        string
	}{
		{name: "eight examples pass", exampleSections: makeExampleSections(section, 8), wantCode: ""},
		{name: "nine examples fail", exampleSections: makeExampleSections(section, 9), wantCode: "TLDR019"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkMaximumExampleCount(&parsedPage{exampleSections: tt.exampleSections}, r)
			require.Equal(t, tt.wantCode, errorCode(r))
		})
	}
}

func TestCheckInfinitiveTense(t *testing.T) {
	infinitive := commandSection{
		description:           "List all files:",
		descriptionLineNumber: 1,
	}
	present := commandSection{
		description:           "Writes files:",
		descriptionLineNumber: 2,
	}
	gerund := commandSection{
		description:           "Writing files:",
		descriptionLineNumber: 3,
	}

	tests := []struct {
		name            string
		exampleSections []commandSection
		wantCodes       []string
	}{
		{name: "infinitive tense passes", exampleSections: []commandSection{infinitive}, wantCodes: nil},
		{name: "present tense fails", exampleSections: []commandSection{present}, wantCodes: []string{"TLDR104"}},
		{name: "gerund fails", exampleSections: []commandSection{gerund}, wantCodes: []string{"TLDR104"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkInfinitiveTense(&parsedPage{exampleSections: tt.exampleSections}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckSingleCommandPerExample(t *testing.T) {
	command_1 := parsedLine{
		kind:               kindCommand,
		lineNumber:         2,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	command_2 := parsedLine{
		kind:               kindCommand,
		lineNumber:         3,
		rawLine:            "`ls -la`",
		content:            "ls -la",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name            string
		exampleSections []commandSection
		wantCodes       []string
	}{
		{
			name: "single command passes",
			exampleSections: []commandSection{
				{
					description:           "List all files:",
					descriptionLineNumber: 1,
					commands:              []parsedLine{command_1},
				},
			},
		},
		{
			name: "two commands fail on the second",
			exampleSections: []commandSection{
				{
					description:           "List all files:",
					descriptionLineNumber: 1,
					commands: []parsedLine{
						command_1,
						command_2,
					},
				},
			},
			wantCodes: []string{"TLDR105"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkSingleCommandPerExample(&parsedPage{exampleSections: tt.exampleSections}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestLineIndex(t *testing.T) {
	tests := []struct {
		name       string
		lines      []parsedLine
		lineNumber int
		want       int
	}{
		{
			name:       "line not found returns minus one",
			lines:      []parsedLine{{lineNumber: 0}, {lineNumber: 1}},
			lineNumber: 5,
			want:       -1,
		},
		{
			name:       "empty lines returns minus one",
			lines:      nil,
			lineNumber: 0,
			want:       -1,
		},
		{
			name:       "line is first",
			lines:      []parsedLine{{lineNumber: 2}, {lineNumber: 4}},
			lineNumber: 2,
			want:       0,
		},
		{
			name:       "line is in the middle",
			lines:      []parsedLine{{lineNumber: 0}, {lineNumber: 3}, {lineNumber: 9}},
			lineNumber: 3,
			want:       1,
		},
		{
			name:       "duplicate line numbers returns first",
			lines:      []parsedLine{{lineNumber: 5}, {lineNumber: 7}, {lineNumber: 5}},
			lineNumber: 5,
			want:       0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lineIndex(tt.lines, tt.lineNumber)
			require.Equal(t, tt.want, got)
		})
	}
}

// makeExampleSections returns a slice containing n copies of section.
func makeExampleSections(section commandSection, n int) []commandSection {
	sections := make([]commandSection, n)
	for i := range sections {
		sections[i] = section
	}
	return sections
}
