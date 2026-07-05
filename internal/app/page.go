package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/internal/render"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/pathutil"
)

// lookupAndRenderPage finds a page by name and renders it to the terminal.
// Returns 0 on success, 1 on error.
func (a *App) lookupAndRenderPage(cli *cmd.CLI) int {
	p := a.resolvePlatform(cli.Platform)
	langs := a.resolveLanguages(cli.Languages)
	c := cache.New()

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

	data, err := os.ReadFile(pagePath)
	if err != nil {
		logger.Error("failed to read page: %v", err)
		return 1
	}

	page := render.Parse(string(data))
	page.Path = pagePath

	renderer := render.New(a.Stdout, a.renderOptions(cli)...)
	if err := renderer.Render(renderPlatform, page); err != nil {
		logger.Error("failed to render page: %v", err)
		return 1
	}
	return 0
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
			fmt.Fprintf(
				a.Stderr,
				"%d. %s (tldr --platform %s %s)\n",
				i+1,
				platform,
				platform,
				query,
			)
		}
	}

	switch {
	case len(results.Matches) > 0:
		return results.Matches[0], requestedPlatform, nil
	case len(results.Fallbacks) > 0:
		page := results.Fallbacks[0]
		return page, pathutil.PagePlatform(page), nil

	default:
		return "", "", fmt.Errorf("page not found, try running tldr --update")
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
	case cli.Compact:
		output.Compact = true
	case cli.Raw:
		output.RawMarkdown = true
	case cli.Edit:
		output.EditLink = true
	case cli.ShortOptions:
		output.OptionStyle = config.OptionStyleShort
	case cli.LongOptions:
		output.OptionStyle = config.OptionStyleLong
	}

	opts = append(opts, render.WithOutput(output))
	return opts
}
