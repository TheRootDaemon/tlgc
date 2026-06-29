package render

import (
	"regexp"
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
)

var (
	placeholderPattern = regexp.MustCompile(`\{\{(.*?)\}\}`)
	optionPattern      = regexp.MustCompile(`^\[(.+)\|(.+)\]$`)
)

type Kind int

const (
	Text Kind = iota
	Placeholder
	Option
)

type Segment struct {
	Kind  Kind
	Text  string
	Short string
	Long  string
}

type Example struct {
	Description string
	Command     string
}

type Page struct {
	Title       string
	URL         string
	Path        string
	Description []string
	Examples    []Example
}

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
