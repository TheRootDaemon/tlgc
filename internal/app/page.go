package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/internal/render"
	"github.com/TheRootDaemon/tlgc/internal/upstream"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/pathutil"
	"github.com/TheRootDaemon/tlgc/termcolor"
)

// lookupAndRenderPage finds a page by name and renders it to the terminal.
// Returns 0 on success, 1 on error.
func (a *App) lookupAndRenderPage(cli *cmd.CLI) int {
	p := a.resolvePlatform(cli.Platform)
	langs := a.resolveLanguages(cli.Languages)
	c := cache.New()

	if !cli.Offline {
		cfg := config.Cache()
		if cfg.AutoUpdate && c.NeedsUpdate(cfg.MaxAge) {
			client := upstream.New()
			if err := c.Update(context.Background(), langs, client); err != nil {
				logger.Warn("auto-update failed: %v", err)
			}
		}
	}

	query := strings.Join(cli.Page, "-")
	results, err := c.Find(query, p, langs)
	if err != nil {
		logger.Error("failed to find page: %v", err)
		return 1
	}

	pagePath, renderPlatform, err := a.selectPage(
		results,
		query,
		p,
	)
	if err != nil {
		logger.Error("%v", err)
		return 1
	}

	page, err := loadPage(pagePath)
	if err != nil {
		logger.Error("%v", err)
		return 1
	}

	if err := a.renderPage(cli, renderPlatform, page); err != nil {
		logger.Error("failed to render: %v", err)
		return 1
	}
	return 0
}

// renderLocalFile reads a local tldr markdown file, validates it,
// and renders it to the terminal.
// Returns 0 on success, 1 on error.
func (a *App) renderLocalFile(cli *cmd.CLI) int {
	page, err := loadPage(cli.Render)
	if err != nil {
		logger.Error("%v", err)
		return 1
	}

	if err := a.renderPage(cli, "", page); err != nil {
		logger.Error("failed to render: %v", err)
		return 1
	}
	return 0
}

func (a *App) renderPage(
	cli *cmd.CLI,
	platform string,
	page *render.Page,
) error {
	renderer := render.New(a.Stdout, a.renderOptions(cli)...)
	return renderer.Render(platform, page)
}

// selectPage chooses the best matching page
// and falls back to pages from other platforms
// when no exact match exists.
func (a *App) selectPage(
	results *cache.FindResult,
	query,
	requestedPlatform string,
) (string, string, error) {
	if len(results.Fallbacks) > 0 {
		logger.Warn(
			"%d page(s) found for other platforms:",
			len(results.Fallbacks),
		)

		for i, f := range results.Fallbacks {
			platform := pathutil.PagePlatform(f)
			_, err := fmt.Fprintf(
				a.Stderr,
				"%d. %s (tldr --platform %s %s)\n",
				i+1,
				platform,
				platform,
				query,
			)
			if err != nil {
				return "", "", err
			}
		}
	}

	switch {
	case len(results.Matches) > 0:
		return results.Matches[0], requestedPlatform, nil
	case len(results.Fallbacks) > 0:
		page := results.Fallbacks[0]
		return page, pathutil.PagePlatform(page), nil

	default:
		const (
			pageNotFound = `page not found, try running tldr --update

If the page does not exist, you can create an issue here:
%s
or document it yourself and create a pull request here:
%s`
			tldrIssues = "https://github.com/tldr-pages/tldr/issues"
			tldrPulls  = "https://github.com/tldr-pages/tldr/pulls"
		)
		return "", "", fmt.Errorf(
			pageNotFound,
			termcolor.Sprint("bold", tldrIssues),
			termcolor.Sprint("bold", tldrPulls),
		)
	}
}

// renderOptions builds the render options from the CLI flags and config.
func (a *App) renderOptions(cli *cmd.CLI) []render.RenderOption {
	opts := []render.RenderOption{
		render.WithWriter(a.Stdout),
	}
	switch cli.Color {
	case "always":
		opts = append(opts, render.WithColor(true))
	case "never":
		opts = append(opts, render.WithColor(false))
	}

	output := config.Output()
	switch {
	case cli.NoCompact:
		output.Compact = false
	case cli.Compact:
		output.Compact = true
	}

	switch {
	case cli.NoRaw:
		output.RawMarkdown = false
	case cli.Raw:
		output.RawMarkdown = true
	}

	switch {
	case cli.ShortOptions && cli.LongOptions:
		output.OptionStyle = config.OptionStyleCombined
	case cli.ShortOptions:
		output.OptionStyle = config.OptionStyleShort
	case cli.LongOptions:
		output.OptionStyle = config.OptionStyleLong
	}

	if cli.Edit {
		output.EditLink = true
	}

	opts = append(opts, render.WithOutput(output))
	return opts
}

// loadPage reads and validates the TL;DR markdown page at path,
// parses it into a render.Page,
// and sets the page's Path and RawContent fields.
// It returns an error if the file cannot be read
// or is not a valid TLDR page.
func loadPage(path string) (*render.Page, error) {
	root, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = root.Close()
	}()

	data, err := root.ReadFile(filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if err := render.Validate(string(data)); err != nil {
		return nil, fmt.Errorf("invalid tldr page: %w", err)
	}

	page := render.Parse(string(data))
	page.Path = path
	page.RawContent = string(data)

	return page, nil
}
