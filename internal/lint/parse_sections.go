package lint

import "strings"

// indexOfTitle returns the index of the first title line,
// or -1 if the page has no title.
func indexOfTitle(lines []parsedLine) int {
	for i, l := range lines {
		if l.kind == kindTitle {
			return i
		}
	}
	return -1
}

// nextContentIndex returns the index of the first non-blank line at or after i.
func nextContentIndex(lines []parsedLine, i int) int {
	for i < len(lines) && lines[i].kind == kindBlank {
		i++
	}
	return i
}

// collectDescriptions gathers the consecutive description lines starting at i.
//
// It returns the descriptions,
// the subset that are information links,
// and the index of the first line after the descriptions.
func collectDescriptions(lines []parsedLine, i int) ([]parsedLine, []parsedLine, int) {
	descriptions := make([]parsedLine, 0, len(lines))
	infoLinks := make([]parsedLine, 0, len(lines))

	for i < len(lines) && lines[i].kind == kindDescription {
		line := lines[i]
		descriptions = append(descriptions, line)

		if isInfoLink(line.content) {
			infoLinks = append(infoLinks, line)
		}

		i++
	}

	return descriptions, infoLinks, i
}

// collectExampleSections walks the lines starting at i,
// building one commandSection per example-description line (followed by its commands).
func collectExampleSections(lines []parsedLine, i int) []commandSection {
	var sections []commandSection

	for i < len(lines) {
		i = nextContentIndex(lines, i)
		if i >= len(lines) {
			break
		}
		if lines[i].kind != kindExampleDesc {
			// stray text or other line – skip it;
			// rules can inspect p.lines.
			i++
			continue
		}

		var section commandSection
		section, i = buildExampleSection(lines, i)
		sections = append(sections, section)
	}
	return sections
}

// buildExampleSection consumes a single example description
// and the command lines that follow it,
// starting at the description line.
//
// It returns the section and the index of the first line after it.
func buildExampleSection(lines []parsedLine, i int) (section commandSection, next int) {
	section = commandSection{
		description:           lines[i].content,
		descriptionLineNumber: lines[i].lineNumber,
	}

	i = nextContentIndex(lines, i+1)

	for i < len(lines) && lines[i].kind == kindCommand {
		section.commands = append(section.commands, lines[i])
		i = nextContentIndex(lines, i+1)
	}
	return section, i
}

// isInfoLink reports whether a description value looks like an information link.
func isInfoLink(content string) bool {
	return strings.HasPrefix(content, "More information: ")
}
