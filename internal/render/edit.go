package render

import (
	"fmt"
	"io"

	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/pathutil"
)

// renderEditLink writes the edit URL to w, with a trailing newline
// unless output is in compact mode.
func (r *Renderer) renderEditLink(w io.Writer, url string) error {
	logger.Info("edit this page on GitHub")
	logger.Trace("url=%q compact=%t", url, r.output.Compact)
	_, err := io.WriteString(w, url)
	if err != nil {
		return err
	}

	if !r.output.Compact {
		_, err := io.WriteString(w, "\n")
		return err
	}

	return nil
}

// renderPageEditLink renders the edit link
// for p when edit links are enabled.
func (r *Renderer) renderPageEditLink(p *Page) error {
	if !r.output.EditLink {
		return nil
	}

	url := buildEditURL(p.Path, p.URL)
	if url == "" {
		return nil
	}

	return r.renderEditLink(r.w, url)
}

// buildEditURL returns the GitHub edit URL for a tldr page.
// If url is non-empty it is returned as-is (the page has a custom source).
// Otherwise the URL is constructed from the page's file path.
func buildEditURL(path, url string) string {
	if url != "" {
		logger.Trace("using custom url=%s", url)
		return url
	}

	if path != "" {
		page := pathutil.PageName(path)
		platform := pathutil.PagePlatform(path)
		result := fmt.Sprintf(
			"https://github.com/tldr-pages/tldr/edit/main/pages/%s/%s.md",
			platform,
			page,
		)
		logger.Trace("path=%s -> %s", path, result)
		return result
	}

	logger.Trace("empty path and url")
	return ""
}
