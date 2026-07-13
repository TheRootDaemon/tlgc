package render

import (
	"os"

	"github.com/TheRootDaemon/tlgc/logger"
)

// renderRaw reads the raw markdown file at p.Path and
// writes it to the Renderer's writer.
func (r *Renderer) renderRaw(p *Page) error {
	if p.RawContent != "" {
		logger.Trace("using RawContent (%d bytes)", len(p.RawContent))
		_, err := r.w.Write(
			[]byte(p.RawContent),
		)
		return err
	}

	logger.Trace("reading from path=%s", p.Path)
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return err
	}

	_, err = r.w.Write(data)
	return err
}
