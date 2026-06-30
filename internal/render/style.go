package render

import (
	"strings"

	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/termcolor"
)

func (r *Renderer) applyStyle(s config.OutputStyle, t string) string {
	if !r.useColor {
		return t
	}

	return termcolor.Sprint(styleString(s), t)
}

func (r *Renderer) styleForSegment(s *Segment) config.OutputStyle {
	switch s.Kind {
	case Text:
		return r.style.Example
	case Placeholder:
		return r.style.Placeholder
	case Option:
		return r.style.Placeholder
	default:
		return r.style.Example
	}
}

func styleString(s config.OutputStyle) string {
	var parts []string
	if s.Bold {
		parts = append(parts, "bold")
	}
	if s.Dim {
		parts = append(parts, "dim")
	}
	if s.Italic {
		parts = append(parts, "italic")
	}
	if s.Underline {
		parts = append(parts, "underline")
	}
	if s.Strikethrough {
		parts = append(parts, "strikethrough")
	}
	if s.Color.Kind == config.ColorKindNamed && s.Color.Named != config.ColorDefault && s.Color.Named != "" {
		parts = append(parts, string(s.Color.Named))
	}
	if s.Background.Kind == config.ColorKindNamed && s.Background.Named != config.ColorDefault && s.Background.Named != "" {
		parts = append(parts, "on_"+string(s.Background.Named))
	}

	return strings.Join(parts, " ")
}
