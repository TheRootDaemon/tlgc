package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckDescriptionStartsWithCapital(t *testing.T) {
	upper := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> Upper case.",
		content:    "Upper case.",
	}
	lower := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> lower case.",
		content:    "lower case.",
	}
	npm := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> npm install.",
		content:    "npm install.",
	}
	pnpm := parsedLine{
		kind:       kindDescription,
		lineNumber: 4,
		rawLine:    "> pnpm add.",
		content:    "pnpm add.",
	}
	empty := parsedLine{
		kind:       kindDescription,
		lineNumber: 5,
		rawLine:    ">",
		content:    "",
	}

	tests := []struct {
		name         string
		descriptions []parsedLine
		wantCodes    []string
	}{
		{name: "uppercase start passes", descriptions: []parsedLine{upper}, wantCodes: nil},
		{name: "lowercase start fails", descriptions: []parsedLine{lower}, wantCodes: []string{"TLDR003"}},
		{name: "npm exception passes", descriptions: []parsedLine{npm}, wantCodes: nil},
		{name: "pnpm exception passes", descriptions: []parsedLine{pnpm}, wantCodes: nil},
		{name: "empty description passes", descriptions: []parsedLine{empty}, wantCodes: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkDescriptionStartsWithCapital(&parsedPage{descriptions: tt.descriptions}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckDescriptionEndsWithPeriod(t *testing.T) {
	withPeriod := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> Description.",
		content:    "Description.",
	}
	withoutPeriod := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> Description",
		content:    "Description",
	}
	infoLink := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> More information: <https://example.com>",
		content:    "More information: <https://example.com>",
	}
	empty := parsedLine{
		kind:       kindDescription,
		lineNumber: 4,
		rawLine:    ">",
		content:    "",
	}

	tests := []struct {
		name         string
		descriptions []parsedLine
		wantCodes    []string
	}{
		{name: "ends with period passes", descriptions: []parsedLine{withPeriod}, wantCodes: nil},
		{name: "missing period fails", descriptions: []parsedLine{withoutPeriod}, wantCodes: []string{"TLDR004"}},
		{name: "info link without period passes", descriptions: []parsedLine{infoLink}, wantCodes: nil},
		{name: "empty description passes", descriptions: []parsedLine{empty}, wantCodes: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkDescriptionEndsWithPeriod(&parsedPage{descriptions: tt.descriptions}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckInformationLinkLabel(t *testing.T) {
	exact := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> More information: <https://example.com>",
		content:    "More information: <https://example.com>",
	}
	shortened := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> More info: <https://example.com>",
		content:    "More info: <https://example.com>",
	}
	noColon := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> More information <https://example.com>",
		content:    "More information <https://example.com>",
	}
	lowercase := parsedLine{
		kind:       kindDescription,
		lineNumber: 4,
		rawLine:    "> more information: <https://example.com>",
		content:    "more information: <https://example.com>",
	}
	notALink := parsedLine{
		kind:       kindDescription,
		lineNumber: 5,
		rawLine:    "> Description.",
		content:    "Description.",
	}

	tests := []struct {
		name         string
		descriptions []parsedLine
		wantCodes    []string
	}{
		{name: "exact label passes", descriptions: []parsedLine{exact}, wantCodes: nil},
		{name: "not a link passes", descriptions: []parsedLine{notALink}, wantCodes: nil},
		{name: "shortened label fails", descriptions: []parsedLine{shortened}, wantCodes: []string{"TLDR016"}},
		{name: "missing colon fails", descriptions: []parsedLine{noColon}, wantCodes: []string{"TLDR016"}},
		{name: "lowercase label fails", descriptions: []parsedLine{lowercase}, wantCodes: []string{"TLDR016"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkInformationLinkLabel(&parsedPage{descriptions: tt.descriptions}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckInformationLinkBrackets(t *testing.T) {
	bracketed := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> More information: <https://example.com>",
		content:    "More information: <https://example.com>",
	}
	unbracketed := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> More information: https://example.com",
		content:    "More information: https://example.com",
	}
	missingOpen := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> More information: https://example.com>",
		content:    "More information: https://example.com>",
	}

	tests := []struct {
		name      string
		infoLinks []parsedLine
		wantCodes []string
	}{
		{name: "bracketed link passes", infoLinks: []parsedLine{bracketed}, wantCodes: nil},
		{name: "unbracketed link fails", infoLinks: []parsedLine{unbracketed}, wantCodes: []string{"TLDR017"}},
		{name: "missing opening bracket fails", infoLinks: []parsedLine{missingOpen}, wantCodes: []string{"TLDR017"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkInformationLinkBrackets(&parsedPage{infoLinks: tt.infoLinks}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckSingleInformationLink(t *testing.T) {
	first := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> More information: <https://example.com>",
		content:    "More information: <https://example.com>",
	}
	second := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> More information: <https://example.org>",
		content:    "More information: <https://example.org>",
	}
	third := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> More information: <https://example.net>",
		content:    "More information: <https://example.net>",
	}

	tests := []struct {
		name      string
		infoLinks []parsedLine
		wantCodes []string
	}{
		{name: "no link passes", infoLinks: nil, wantCodes: nil},
		{name: "single link passes", infoLinks: []parsedLine{first}, wantCodes: nil},
		{name: "two links fail on extra", infoLinks: []parsedLine{first, second}, wantCodes: []string{"TLDR018"}},
		{
			name:      "three links fail on extras",
			infoLinks: []parsedLine{first, second, third},
			wantCodes: []string{"TLDR018", "TLDR018"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkSingleInformationLink(&parsedPage{infoLinks: tt.infoLinks}, r)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckNoteLabelFormat(t *testing.T) {
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> Description.",
		content:    "Description.",
	}
	lowercaseNote := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> note: something",
		content:    "note: something",
	}
	uppercaseNote := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> NOTE: something",
		content:    "NOTE: something",
	}
	exampleWithNote := commandSection{
		description:           "note: something",
		descriptionLineNumber: 4,
	}

	tests := []struct {
		name            string
		descriptions    []parsedLine
		exampleSections []commandSection
		wantCodes       []string
	}{
		{name: "plain description passes", descriptions: []parsedLine{description}, wantCodes: nil},
		{
			name:         "lowercase note fails",
			descriptions: []parsedLine{lowercaseNote},
			wantCodes:    []string{"TLDR020"},
		},
		{
			name:         "uppercase note fails",
			descriptions: []parsedLine{uppercaseNote},
			wantCodes:    []string{"TLDR020"},
		},
		{
			name:            "note in example description fails",
			exampleSections: []commandSection{exampleWithNote},
			wantCodes:       []string{"TLDR020"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkNoteLabelFormat(
				&parsedPage{
					descriptions:    tt.descriptions,
					exampleSections: tt.exampleSections,
				},
				r,
			)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}

func TestCheckStandardTermsInBackticks(t *testing.T) {
	backticked := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> Use `stdout`.",
		content:    "Use `stdout`.",
	}
	unbackticked := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> Writes to stdout.",
		content:    "Writes to stdout.",
	}
	inURL := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> See <https://example.com/stdout>.",
		content:    "See <https://example.com/stdout>.",
	}
	partOfWord := parsedLine{
		kind:       kindDescription,
		lineNumber: 4,
		rawLine:    "> List stdoutstreams.",
		content:    "List stdoutstreams.",
	}
	betweenSpans := parsedLine{
		kind:       kindDescription,
		lineNumber: 5,
		rawLine:    "> Use `foo` and stdin and `bar`.",
		content:    "Use `foo` and stdin and `bar`.",
	}
	multipleSpans := parsedLine{
		kind: kindDescription, lineNumber: 6,
		rawLine: "> Use `stdin`, `stdout`, and `stderr`.",
		content: "Use `stdin`, `stdout`, and `stderr`.",
	}
	unclosedBacktick := parsedLine{
		kind:       kindDescription,
		lineNumber: 7,
		rawLine:    "> Use `stdin to read.",
		content:    "Use `stdin to read.",
	}
	exampleWithTerm := commandSection{
		description:           "Send output to stderr:",
		descriptionLineNumber: 8,
	}

	tests := []struct {
		name            string
		descriptions    []parsedLine
		exampleSections []commandSection
		wantCodes       []string
	}{
		{name: "backticked term passes", descriptions: []parsedLine{backticked}, wantCodes: nil},
		{name: "term in URL passes", descriptions: []parsedLine{inURL}, wantCodes: nil},
		{name: "term inside word passes", descriptions: []parsedLine{partOfWord}, wantCodes: nil},
		{
			name:         "multiple backticked terms pass",
			descriptions: []parsedLine{multipleSpans},
			wantCodes:    nil,
		},
		{
			name:         "unclosed backtick before term passes",
			descriptions: []parsedLine{unclosedBacktick},
			wantCodes:    nil,
		},
		{
			name:         "unbackticked term fails",
			descriptions: []parsedLine{unbackticked},
			wantCodes:    []string{"TLDR112"},
		},
		{
			name:         "term between code spans fails",
			descriptions: []parsedLine{betweenSpans},
			wantCodes:    []string{"TLDR112"},
		},
		{
			name:            "term in example description fails",
			exampleSections: []commandSection{exampleWithTerm},
			wantCodes:       []string{"TLDR112"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkStandardTermsInBackticks(
				&parsedPage{
					descriptions:    tt.descriptions,
					exampleSections: tt.exampleSections,
				},
				r,
			)
			require.Equal(t, tt.wantCodes, errorCodes(r))
		})
	}
}
