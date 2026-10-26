package render

import (
	"io"
	"strings"
)

// renderExample writes one example,
// a bullet line for the description (prefixed with ExamplePrefix when ShowHyphens is set)
// followed by the styled command text on the next line.
// In compact mode the blank line between bullet and command is omitted.
func (r *Renderer) renderExample(example Example) error {
	indent := strings.Repeat(" ", r.indent.Bullet)
	desc := example.Description

	if r.output.ShowHyphens {
		desc = r.output.ExamplePrefix + desc
	}

	if err := r.renderBulletLine(desc, indent); err != nil {
		return err
	}

	if !r.output.Compact {
		_, err := io.WriteString(r.w, "\n")
		if err != nil {
			return err
		}
	}

	_, err := io.WriteString(r.w, "\n")
	if err != nil {
		return err
	}

	if example.Command != "" {
		segments := ParseCommand(example.Command)
		if err := r.renderCommand(segments); err != nil {
			return err
		}
	}

	return nil
}

// renderExamples renders a sequence of examples to the configured [io.Writer].
//
// Each example is rendered using renderExample.
// In non-compact mode, a blank line is inserted
// between consecutive examples
func (r *Renderer) renderExamples(examples []Example) error {
	for i, example := range examples {
		if i > 0 && !r.output.Compact {
			_, err := io.WriteString(r.w, "\n")
			if err != nil {
				return err
			}
		}

		if err := r.renderExample(example); err != nil {
			return err
		}
	}

	return nil
}
