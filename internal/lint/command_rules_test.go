package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckCommandWhitespace(t *testing.T) {
	clean := parsedLine{
		kind:               kindCommand,
		lineNumber:         1,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	leadingSpace := parsedLine{
		kind:               kindCommand,
		lineNumber:         2,
		rawLine:            "` ls`",
		content:            " ls",
		hasClosingBacktick: true,
	}
	trailingSpace := parsedLine{
		kind:               kindCommand,
		lineNumber:         3,
		rawLine:            "`ls `",
		content:            "ls ",
		hasClosingBacktick: true,
	}
	escapedLeadingSpace := parsedLine{
		kind:               kindCommand,
		lineNumber:         4,
		rawLine:            "`\\ ls`",
		content:            `\ ls`,
		hasClosingBacktick: true,
	}
	escapedTrailingSpace := parsedLine{
		kind:               kindCommand,
		lineNumber:         5,
		rawLine:            "`ls\\ `",
		content:            `ls\ `,
		hasClosingBacktick: true,
	}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{name: "clean command passes", lines: []parsedLine{clean}, wantCodes: nil},
		{name: "leading space fails", lines: []parsedLine{leadingSpace}, wantCodes: []string{"TLDR021"}},
		{name: "trailing space fails", lines: []parsedLine{trailingSpace}, wantCodes: []string{"TLDR021"}},
		{name: "escaped leading space passes", lines: []parsedLine{escapedLeadingSpace}, wantCodes: nil},
		{name: "escaped trailing space passes", lines: []parsedLine{escapedTrailingSpace}, wantCodes: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkCommandWhitespace(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckCommandDescriptionAnnotated(t *testing.T) {
	text := parsedLine{
		kind:       kindText,
		lineNumber: 1,
		rawLine:    "unannotated",
		content:    "unannotated",
	}
	exampleDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 2,
		rawLine:    "- Example:",
		content:    "Example:",
	}
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> Description.",
		content:    "Description.",
	}
	command := parsedLine{
		kind:               kindCommand,
		lineNumber:         4,
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
			name:      "no stray text passes",
			lines:     []parsedLine{description, exampleDescription, command},
			wantCodes: nil,
		},
		{
			name:      "text before example description fails",
			lines:     []parsedLine{text, exampleDescription, command},
			wantCodes: []string{"TLDR101"},
		},
		{
			name:      "text before description passes",
			lines:     []parsedLine{text, description},
			wantCodes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkCommandDescriptionAnnotated(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckExampleDescriptionAnnotated(t *testing.T) {
	text := parsedLine{
		kind:       kindText,
		lineNumber: 1,
		rawLine:    "unannotated",
		content:    "unannotated",
	}
	exampleDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 2,
		rawLine:    "- Example:",
		content:    "Example:",
	}
	command := parsedLine{
		kind:               kindCommand,
		lineNumber:         3,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 4,
		rawLine:    "> Description.",
		content:    "Description.",
	}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{
			name:      "no stray text passes",
			lines:     []parsedLine{exampleDescription, command},
			wantCodes: nil,
		},
		{
			name:      "text after example description fails",
			lines:     []parsedLine{exampleDescription, text, command},
			wantCodes: []string{"TLDR102"},
		},
		{
			name:      "text followed by command fails",
			lines:     []parsedLine{text, command},
			wantCodes: []string{"TLDR102"},
		},
		{
			name:      "text followed by description passes",
			lines:     []parsedLine{text, description},
			wantCodes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkExampleDescriptionAnnotated(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckCommandClosingBacktick(t *testing.T) {
	closed := parsedLine{
		kind:               kindCommand,
		lineNumber:         1,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	unclosed := parsedLine{
		kind:               kindCommand,
		lineNumber:         2,
		rawLine:            "`ls",
		content:            "ls",
		hasClosingBacktick: false,
	}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{
			name:      "closed command passes",
			lines:     []parsedLine{closed},
			wantCodes: nil,
		},
		{
			name:      "unclosed command fails",
			lines:     []parsedLine{unclosed},
			wantCodes: []string{"TLDR103"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkCommandClosingBacktick(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckCommandNotEmpty(t *testing.T) {
	nonEmpty := parsedLine{
		kind:               kindCommand,
		lineNumber:         1,
		rawLine:            "`ls`",
		content:            "ls",
		hasClosingBacktick: true,
	}
	empty := parsedLine{
		kind:               kindCommand,
		lineNumber:         2,
		rawLine:            "``",
		content:            "",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name      string
		lines     []parsedLine
		wantCodes []string
	}{
		{
			name:      "non-empty command passes",
			lines:     []parsedLine{nonEmpty},
			wantCodes: nil,
		},
		{
			name:      "empty command fails",
			lines:     []parsedLine{empty},
			wantCodes: []string{"TLDR110"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkCommandNotEmpty(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}
