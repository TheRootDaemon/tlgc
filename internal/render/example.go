package render

import (
	"io"
	"strings"

	"github.com/TheRootDaemon/tlgc/logger"
)

// renderExample writes one example,
// a bullet line for the description (prefixed with ExamplePrefix when ShowHyphens is set)
// followed by the styled command text on the next line.
// In compact mode the blank line between bullet and command is omitted.
func (r *Renderer) renderExample(w io.Writer, ex Example) error {
	logger.Trace(
		"desc=%q hasCmd=%t compact=%t",
		ex.Description,
		ex.Command != "",
		r.output.Compact,
	)
	indent := strings.Repeat(" ", r.indent.Bullet)
	desc := ex.Description

	if r.output.ShowHyphens {
		desc = r.output.ExamplePrefix + desc
	}

	if err := r.renderBulletLine(w, desc, indent); err != nil {
		return err
	}

	if !r.output.Compact {
		_, err := io.WriteString(w, "\n")
		if err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "\n")
	if err != nil {
		return err
	}

	if ex.Command != "" {
		segments := ParseCommand(ex.Command)
		if err := r.renderCommand(w, segments); err != nil {
			return err
		}
	}

	return nil
}

// renderExamples renders a sequence of examples to w.
//
// Each example is rendered using renderExample.
// In non-compact mode, a blank line is inserted
// between consecutive examples
func (r *Renderer) renderExamples(w io.Writer, examples []Example) error {
	for i, ex := range examples {
		if i > 0 && !r.output.Compact {
			_, err := io.WriteString(w, "\n")
			if err != nil {
				return err
			}
		}

		if err := r.renderExample(r.w, ex); err != nil {
			return err
		}
	}

	return nil
}
