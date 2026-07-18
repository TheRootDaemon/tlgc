package render

import "os"

// renderRaw reads the raw markdown file at p.Path and
// writes it to the Renderer's writer.
func (r *Renderer) renderRaw(p *Page) error {
	if p.RawContent != "" {
		_, err := r.w.Write(
			[]byte(p.RawContent),
		)
		return err
	}

	data, err := os.ReadFile(p.Path)
	if err != nil {
		return err
	}

	_, err = r.w.Write(data)
	return err
}
