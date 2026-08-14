package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	raw := "# App\n\n> Brief description.\n> More information: https://example.com\n\n- Copy files\n\n`cp file file.bak`\n\n- Create a backup\n\n`tar czf backup.tar.gz file`"

	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 1,
		rawLine:    "# App",
		content:    "App",
	}
	blank_1 := parsedLine{
		kind:       kindBlank,
		lineNumber: 2,
		rawLine:    "",
	}
	desc := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> Brief description.",
		content:    "Brief description.",
	}
	link := parsedLine{
		kind:       kindDescription,
		lineNumber: 4,
		rawLine:    "> More information: https://example.com",
		content:    "More information: https://example.com",
	}
	blank_4 := parsedLine{
		kind:       kindBlank,
		lineNumber: 5,
		rawLine:    "",
	}
	copyDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 6,
		rawLine:    "- Copy files",
		content:    "Copy files",
	}
	blank_6 := parsedLine{
		kind:       kindBlank,
		lineNumber: 7,
		rawLine:    "",
	}
	copyCommand := parsedLine{
		kind:       kindCommand,
		lineNumber: 8,
		rawLine:    "`cp file file.bak`",
		content:    "cp file file.bak", hasClosingBacktick: true,
	}
	blank_8 := parsedLine{
		kind:       kindBlank,
		lineNumber: 9,
		rawLine:    "",
	}
	backupDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 10,
		rawLine:    "- Create a backup",
		content:    "Create a backup",
	}
	blank_10 := parsedLine{
		kind:       kindBlank,
		lineNumber: 11,
		rawLine:    "",
	}
	backupCommand := parsedLine{
		kind:               kindCommand,
		lineNumber:         12,
		rawLine:            "`tar czf backup.tar.gz file`",
		content:            "tar czf backup.tar.gz file",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name string
		raw  string
		want *parsedPage
	}{
		{
			name: "empty input returns nil",
			raw:  "",
			want: nil,
		},
		{
			name: "whitespace only input returns nil",
			raw:  " \n\t ",
			want: nil,
		},
		{
			name: "page without a title",
			raw:  "some text\n",
			want: &parsedPage{
				rawContent: "some text\n",
				lines: []parsedLine{
					{kind: kindText, lineNumber: 1, rawLine: "some text", content: "some text"},
					{kind: kindBlank, lineNumber: 2, rawLine: ""}, // trailing newline splits into a blank line
				},
			},
		},
		{
			name: "full page end to end",
			raw:  raw,
			want: &parsedPage{
				rawContent:      raw,
				lines:           []parsedLine{title, blank_1, desc, link, blank_4, copyDescription, blank_6, copyCommand, blank_8, backupDescription, blank_10, backupCommand},
				title:           "App",
				titleLineNumber: 1,
				descriptions:    []parsedLine{desc, link},
				infoLinks:       []parsedLine{link},
				exampleSections: []commandSection{
					{description: "Copy files", descriptionLineNumber: 6, commands: []parsedLine{copyCommand}},
					{description: "Create a backup", descriptionLineNumber: 10, commands: []parsedLine{backupCommand}},
				},
			},
		},
		{
			name: "crlf page classifies correctly",
			raw:  "# T\r\n\r\n> D\r\n",
			want: &parsedPage{
				rawContent: "# T\r\n\r\n> D\r\n",
				lines: []parsedLine{
					{kind: kindTitle, lineNumber: 1, rawLine: "# T\r", content: "T"},
					{kind: kindBlank, lineNumber: 2, rawLine: "\r"},
					{kind: kindDescription, lineNumber: 3, rawLine: "> D\r", content: "D"},
					{kind: kindBlank, lineNumber: 4, rawLine: ""},
				},
				title:           "T",
				titleLineNumber: 1,
				descriptions:    []parsedLine{{kind: kindDescription, lineNumber: 3, rawLine: "> D\r", content: "D"}},
				infoLinks:       []parsedLine{},
				exampleSections: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(tt.raw)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestBuildPage(t *testing.T) {
	// lines are hand-built, never via parseLines;
	// parseLines regression must not cascade into buildPage failures.
	title := parsedLine{
		kind:       kindTitle,
		lineNumber: 0,
		rawLine:    "# T",
		content:    "T",
	}
	titleMid := parsedLine{
		kind:       kindTitle,
		lineNumber: 2,
		rawLine:    "# T",
		content:    "T",
	}
	text := parsedLine{
		kind:       kindText,
		lineNumber: 1,
		rawLine:    "stray",
		content:    "stray",
	}
	descriptionBefore := parsedLine{
		kind:       kindDescription,
		lineNumber: 0,
		rawLine:    "> Before",
		content:    "Before",
	}
	descriptionAfter := parsedLine{
		kind:       kindDescription,
		lineNumber: 3,
		rawLine:    "> After",
		content:    "After",
	}
	descriptionOne := parsedLine{
		kind:       kindDescription,
		lineNumber: 1,
		rawLine:    "> One",
		content:    "One",
	}
	descriptionTwo := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> Two",
		content:    "Two",
	}
	descriptionOnly := parsedLine{
		kind:       kindDescription,
		lineNumber: 2,
		rawLine:    "> D",
		content:    "D",
	}
	blank := parsedLine{
		kind:       kindBlank,
		lineNumber: 3,
		rawLine:    "",
	}
	blankGap := parsedLine{
		kind:       kindBlank,
		lineNumber: 1,
		rawLine:    "",
	}
	exampleDescription := parsedLine{
		kind:       kindExampleDesc,
		lineNumber: 4,
		rawLine:    "- Do",
		content:    "Do",
	}
	command := parsedLine{
		kind:               kindCommand,
		lineNumber:         5,
		rawLine:            "`x`",
		content:            "x",
		hasClosingBacktick: true,
	}

	tests := []struct {
		name  string
		raw   string
		lines []parsedLine
		want  *parsedPage
	}{
		{
			name:  "no title leaves grouping empty",
			lines: []parsedLine{text, descriptionAfter},
			want: &parsedPage{
				rawContent: "",
				lines:      []parsedLine{text, descriptionAfter},
				// title, descriptions, infoLinks, exampleSections all nil:
				// never touched when there is no title.
			},
		},
		{
			name:  "title with descriptions and one example",
			lines: []parsedLine{title, descriptionOne, descriptionTwo, blank, exampleDescription, command},
			want: &parsedPage{
				rawContent:      "",
				lines:           []parsedLine{title, descriptionOne, descriptionTwo, blank, exampleDescription, command},
				title:           "T",
				titleLineNumber: 0,
				descriptions:    []parsedLine{descriptionOne, descriptionTwo},
				infoLinks:       []parsedLine{}, // make-backed: non-nil even when empty
				exampleSections: []commandSection{
					{description: "Do", descriptionLineNumber: 4, commands: []parsedLine{command}},
				},
			},
		},
		{
			name:  "blank between title and descriptions is tolerated",
			lines: []parsedLine{title, blankGap, descriptionOnly},
			want: &parsedPage{
				rawContent:      "",
				lines:           []parsedLine{title, blankGap, descriptionOnly},
				title:           "T",
				titleLineNumber: 0,
				descriptions:    []parsedLine{descriptionOnly},
				infoLinks:       []parsedLine{},
				exampleSections: nil, // no examples: bare var stays nil
			},
		},
		{
			name:  "title in the middle ignores earlier lines",
			lines: []parsedLine{descriptionBefore, text, titleMid, descriptionAfter},
			want: &parsedPage{
				rawContent:      "",
				lines:           []parsedLine{descriptionBefore, text, titleMid, descriptionAfter},
				title:           "T",
				titleLineNumber: 2,
				descriptions:    []parsedLine{descriptionAfter}, // grouping anchors at the first title
				infoLinks:       []parsedLine{},
				exampleSections: nil,
			},
		},
		{
			name:  "text right after title yields empty descriptions and nil sections",
			lines: []parsedLine{title, text},
			want: &parsedPage{
				rawContent:      "",
				lines:           []parsedLine{title, text},
				title:           "T",
				titleLineNumber: 0,
				descriptions:    []parsedLine{}, // collectDescriptions is make-backed
				infoLinks:       []parsedLine{},
				exampleSections: nil, // collectExampleSections is a bare var
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPage(tt.raw, tt.lines)
			require.Equal(t, tt.want, got)
		})
	}
}
