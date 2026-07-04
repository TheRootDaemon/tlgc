package app

import (
	"os"
	"strings"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/internal/render"
	"github.com/TheRootDaemon/tlgc/logger"
)

// lookupAndRenderPage finds a page by name and renders it to the terminal.
// Returns 0 on success, 1 on error.
func (a *App) lookupAndRenderPage(cli *cmd.CLI) int {
	p := a.resolvePlatform(cli.Platform)
	langs := a.resolveLanguages(cli.Languages)
	c := cache.New()

	query := strings.Join(cli.Page, "-")
	result, err := c.Find(query, p, langs)
	if err != nil {
		logger.Error("failed to find page: %v", err)
		return 1
	}

	if len(result.Matches) == 0 {
		logger.Error("page not found, try running tldr --update")
		return 1
	}

	data, err := os.ReadFile(result.Matches[0])
	if err != nil {
		logger.Error("failed to read page: %v", err)
		return 1
	}

	page := render.Parse(string(data))
	page.Path = result.Matches[0]

	renderer := render.New(a.Stdout, a.renderOptions(cli)...)
	if err := renderer.Render(p, page); err != nil {
		logger.Error("failed to render page: %v", err)
		return 1
	}
	return 0
}

// renderOptions builds the render options from CLI flags and config.
func (a *App) renderOptions(cli *cmd.CLI) []render.RenderOption {
	opts := []render.RenderOption{render.WithWriter(a.Stdout)}
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
