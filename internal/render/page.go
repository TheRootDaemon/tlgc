package render

import (
	"regexp"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
)

var (
	// placeholderPattern matches tldr placeholder tokens like {{archive.tar}}
	// and captures the inner content (the text between the braces).
	placeholderPattern = regexp.MustCompile(`\{\{(.*?)\}\}`)

	// optionPattern matches the option syntax embedded inside a placeholder,
	// e.g. [-s|--long]. It captures the short form and the long form.
	optionPattern = regexp.MustCompile(`^\[(.+)\|(.+)\]$`)
)

// Kind represents the type of a parsed command segment.
type Kind int

const (
	// Text is a plain literal segment that should be rendered as-is.
	Text Kind = iota

	// Placeholder is a user-supplied value wrapped in {{...}}.
	Placeholder

	// Option is a command-line flag embedded in {{[...|...]}} syntax.
	Option
)

// Segment holds a parsed piece of a command string.
// A segment is either literal text,
// a user-supplied placeholder value,
// or a command-line option with short and long forms.
type Segment struct {
	Kind  Kind
	Text  string
	Short string
	Long  string
}

// Example is a single tldr example
// consisting of a description
// and the associated command text.
type Example struct {
	Description string
	Command     string
}

// Page is a parsed tldr page containing a title,
// zero or more description lines,
// an optional "More information" URL,
// the filesystem path to the source file,
// and a list of examples.
type Page struct {
	Title       string
	URL         string
	Path        string
	Description []string
	Examples    []Example
}

// Parse parses a raw markdown tldr page string into a Page.
//
// It recognises the following markdown structure:
//   - # Title
//   - > Description lines (with optional "More information: <url>" extraction)
//   - - Example descriptions (optionally ending with a colon)
//   - `command` lines (associated with the preceding example)
//
// Nil and empty content produce a Page with zero-valued fields.
func Parse(content string) *Page {
	p := &Page{}
	lines := strings.SplitSeq(content, "\n")

	for line := range lines {
		switch {
		case strings.HasPrefix(line, "# ") && p.Title == "":
			p.Title = line[2:]
		case strings.HasPrefix(line, "> "):
			body := line[2:]
			if url := parseURL(body); url != "" {
				p.URL = url
			} else {
				p.Description = append(p.Description, body)
			}
		case strings.HasPrefix(line, "- "):
			desc := strings.TrimSuffix(line[2:], ":")
			p.Examples = append(
				p.Examples,
				Example{
					Description: desc,
				},
			)
		case strings.HasPrefix(line, "`") && strings.HasSuffix(line, "`"):
			cmd := strings.TrimPrefix(
				strings.TrimSuffix(line, "`"),
				"`",
			)

			if len(p.Examples) == 0 {
				p.Examples = append(
					p.Examples,
					Example{
						Command: cmd,
					},
				)
			} else {
				last := &p.Examples[len(p.Examples)-1]
				if last.Command == "" {
					last.Command = cmd
				}
			}
		}
	}

	return p
}

// ParseCommand splits a raw command string into a slice of Segments.
//
// Placeholder tokens wrapped in {{...}} are extracted
// as either Placeholder or Option segments depending on their inner content.
// Text outside of placeholders becomes Text segments.
// An empty string returns nil.
func ParseCommand(raw string) []Segment {
	if raw == "" {
		return nil
	}

	matches := placeholderPattern.FindAllStringSubmatchIndex(raw, -1)
	if len(matches) == 0 {
		return []Segment{
			{
				Kind: Text,
				Text: raw,
			},
		}
	}

	lastEnd := 0
	var segments []Segment

	for _, match := range matches {
		matchStart := match[0]
		matchEnd := match[1]
		groupStart := match[2]
		groupEnd := match[3]

		if lastEnd < matchStart {
			segments = append(
				segments,
				Segment{
					Kind: Text,
					Text: raw[lastEnd:matchStart],
				},
			)
		}

		inner := raw[groupStart:groupEnd]
		segments = append(
			segments,
			parseInnerPlaceholders(inner),
		)

		lastEnd = matchEnd
	}

	if lastEnd < len(raw) {
		segments = append(
			segments,
			Segment{
				Kind: Text,
				Text: raw[lastEnd:],
			},
		)
	}

	return segments
}

// DisplayText returns the text
// that should be displayed for the segment,
// taking the option display style into account.
//
// For Option segments the return value depends on style:
//   - OptionStyleShort returns the short form (e.g. "-s")
//   - OptionStyleLong returns the long form (e.g. "--long")
//   - OptionStyleCombined returns both (e.g. "[-s|--long]")
//
// All other kinds return the Text field unchanged.
func (s Segment) DisplayText(style config.OptionStyle) string {
	switch s.Kind {
	case Option:
		switch style {
		case config.OptionStyleShort:
			return s.Short
		case config.OptionStyleLong:
			return s.Long
		case config.OptionStyleCombined:
			return "[" + s.Short + "|" + s.Long + "]"
		default:
			return s.Long
		}
	default:
		return s.Text
	}
}

// parseURL extracts the URL from a "More information" description line.
// It expects the URL to be wrapped in angle brackets like <https://example.org>.
// Returns the empty string if no valid URL is found.
func parseURL(body string) string {
	if !strings.Contains(body, "More information: <") {
		return ""
	}

	start := strings.Index(body, "<")
	end := strings.Index(body, ">")
	if start != -1 && end != -1 && start < end {
		return body[start+1 : end]
	}

	return ""
}

// parseInnerPlaceholders parses the inner content of a {{...}} placeholder.
//
// If the inner content matches the option pattern [short|long],
// an Option segment is returned with the short and long forms correctly
// identified regardless of their order.
// Otherwise a plain Placeholder segment is returned.
func parseInnerPlaceholders(inner string) Segment {
	if options := optionPattern.FindStringSubmatch(inner); len(options) == 3 {
		left, right := options[1], options[2]
		short, long := left, right

		if strings.HasPrefix(left, "--") && !strings.HasPrefix(right, "--") {
			short, long = right, left
		}

		return Segment{
			Kind:  Option,
			Long:  long,
			Short: short,
		}
	}

	return Segment{
		Kind: Placeholder,
		Text: inner,
	}
}
