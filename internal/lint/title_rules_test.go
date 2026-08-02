package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckTitleDescriptionSeparator(t *testing.T) {
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 0,
		rawLine:    "# App",
		content:    "App",
	}
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 1,
		rawLine:    "",
	}
	text := parsedLine{
		kind:       kindText,
		lineNumber: 1,
		rawLine:    "stray",
		content:    "stray",
	}
	description_1 := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> D",
		content:    "D",
	}
	description_2 := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> D",
		content:    "D",
	}
	description_3 := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> D",
		content:    "D",
	}

	tests := []struct {
		name            string
		lines           []parsedLine
		descriptions    []parsedLine
		titleLineNumber int
		wantCode        string
	}{
		{
			name:     "no descriptions passes",
			lines:    []parsedLine{title},
			wantCode: "",
		},
		{
			name:            "description immediately after title fails",
			lines:           []parsedLine{title, description_1},
			descriptions:    []parsedLine{description_1},
			titleLineNumber: 0,
			wantCode:        "TLDR006",
		},
		{
			name:            "one blank line between title and description passes",
			lines:           []parsedLine{title, blank, description_2},
			descriptions:    []parsedLine{description_2},
			titleLineNumber: 0,
			wantCode:        "",
		},
		{
			name:            "multiple blank lines pass",
			lines:           []parsedLine{title, blank, blank, description_3},
			descriptions:    []parsedLine{description_3},
			titleLineNumber: 0,
			wantCode:        "",
		},
		{
			name:            "non-blank line between title and description fails",
			lines:           []parsedLine{title, text, description_2},
			descriptions:    []parsedLine{description_2},
			titleLineNumber: 0,
			wantCode:        "TLDR006",
		},
		{
			name:            "description line not present in lines passes",
			lines:           []parsedLine{title, blank, description_2},
			descriptions:    []parsedLine{description_3}, // lineNumber 3 does not exist
			titleLineNumber: 0,
			wantCode:        "",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				r := &Result{}
				checkTitleDescriptionSeparator(
					&parsedPage{
						lines:           tt.lines,
						descriptions:    tt.descriptions,
						titleLineNumber: tt.titleLineNumber,
					},
					r,
				)
				require.Equal(t, tt.wantCode, errorCode(r))
			},
		)
	}
}

func TestCheckTitleCharacters(t *testing.T) {
	tests := []struct {
		name            string
		title           string
		titleLineNumber int
		wantCode        string
	}{
		{name: "empty title passes", title: "", titleLineNumber: 0, wantCode: ""},
		{name: "plain title passes", title: "App", titleLineNumber: 0, wantCode: ""},
		{name: "letters, digits and punctuation pass", title: "Go 2.0 + Plugins!", titleLineNumber: 1, wantCode: ""},
		{name: "single period passes", title: ".", titleLineNumber: 0, wantCode: ""},
		{name: "space before period passes", title: " .", titleLineNumber: 0, wantCode: ""},
		{name: "invalid character fails", title: "App#", titleLineNumber: 0, wantCode: "TLDR013"},
		{name: "at sign fails", title: "App@", titleLineNumber: 0, wantCode: "TLDR013"},
		{name: "trailing period fails", title: "App.", titleLineNumber: 2, wantCode: "TLDR013"},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				r := &Result{}
				checkTitleCharacters(
					&parsedPage{
						title:           tt.title,
						titleLineNumber: tt.titleLineNumber,
					},
					r,
				)
				require.Equal(t, tt.wantCode, errorCode(r))
			},
		)
	}
}

func TestCheckTitleHash(t *testing.T) {
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 2,
		rawLine:    "# App",
		content:    "App",
	}
	description := parsedLine{
		kind:       kindDescription,
		lineNumber: 0,
		rawLine:    "> D",
		content:    "D",
	}
	ccommand := parsedLine{
		kind:               kindCommand,
		lineNumber:         1,
		rawLine:            "`c`",
		content:            "c",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name     string
		lines    []parsedLine
		wantCode string
	}{
		{
			name:     "title present passes",
			lines:    []parsedLine{title},
			wantCode: "",
		},
		{
			name:     "title in the middle passes",
			lines:    []parsedLine{description, ccommand, title},
			wantCode: "",
		},
		{
			name:     "no title fails",
			lines:    []parsedLine{description, ccommand},
			wantCode: "TLDR106",
		},
		{
			name:     "empty page fails",
			lines:    nil,
			wantCode: "TLDR106",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{}
			checkTitleHash(&parsedPage{lines: tt.lines}, r)
			require.Equal(t, tt.wantCode, errorCode(r))
		})
	}
}

func TestIsValidTitleRune(t *testing.T) {
	tests := []struct {
		name string
		in   rune
		want bool
	}{
		{name: "lowercase letter", in: 'a', want: true},
		{name: "uppercase letter", in: 'Z', want: true},
		{name: "digit", in: '7', want: true},
		{name: "underscore", in: '_', want: true},
		{name: "space", in: ' ', want: true},
		{name: "dash", in: '-', want: true},
		{name: "dot", in: '.', want: true},
		{name: "plus", in: '+', want: true},
		{name: "question mark", in: '?', want: true},
		{name: "hash", in: '#', want: false},
		{name: "at sign", in: '@', want: false},
		{name: "ampersand", in: '&', want: false},
		{name: "apostrophe", in: '\'', want: false},
		{name: "slash", in: '/', want: false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				got := isValidTitleRune(tt.in)
				require.Equal(t, tt.want, got)
			},
		)
	}
}
