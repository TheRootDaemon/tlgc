package lint

// lineKind classifies a single line of a tldr page.
type lineKind int

const (
	kindBlank       lineKind = iota
	kindTitle                // starts with #
	kindDescription          // starts with >
	kindExampleDesc          // starts with -
	kindCommand              // starts with `
	kindText                 // everything else
)

// parsedLine holds classification results for a single source line.
type parsedLine struct {
	// kindBlank, kindTitle, kindDescription, kindExampleDesc, kindCommand, or kindText
	kind lineKind

	// 0-indexed line number within the raw content
	lineNumber int

	// original line text (without trailing \n)
	rawLine string

	// extracted content (text after the marker, or between backticks)
	content string

	// for kindCommand: whether a closing backtick was found
	hasClosingBacktick bool
}

// commandSection groups an example description with its command(s).
type commandSection struct {
	// example description text (after "- ")
	description string

	// 0-indexed line number of the description
	descriptionLineNumber int

	// command lines that belong to the example
	commands []parsedLine
}

// parsedPage is the structured representation of a tldr page.
type parsedPage struct {
	// original full content
	rawContent string

	// every parsed line, in source order
	lines []parsedLine

	// first title text (without the leading '#')
	title string

	// 0-indexed line number of the title
	titleLineNumber int

	// consecutive description lines following the title
	descriptions []parsedLine

	// description lines that are "More information: ..." links
	infoLinks []parsedLine

	// example descriptions paired with their commands
	exampleSections []commandSection
}

// parse is the top-level parse entry point.
//
// It returns nil for empty or whitespace-only input.
func parse(raw string) *parsedPage {
	lines := parseLines(raw)
	if lines == nil {
		return nil
	}

	return buildPage(raw, lines)
}

// buildPage groups parsed lines into a structured parsedPage.
func buildPage(raw string, lines []parsedLine) *parsedPage {
	p := &parsedPage{
		rawContent: raw,
		lines:      lines,
	}

	titleIndex := indexOfTitle(lines)
	if titleIndex < 0 {
		// no title means nothing else to group;
		// rules can inspect p.lines.
		return p
	}
	p.title = lines[titleIndex].content
	p.titleLineNumber = lines[titleIndex].lineNumber

	// descriptions (and any info links among them) follow the title,
	// optionally separated by blank lines.
	start := nextContentIndex(lines, titleIndex+1)
	p.descriptions, p.infoLinks, start = collectDescriptions(lines, start)

	// everything after the descriptions forms example sections.
	p.exampleSections = collectExampleSections(lines, start)

	return p
}
