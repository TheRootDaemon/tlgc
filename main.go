package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/cache"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/internal/render"
	"github.com/TheRootDaemon/tlgc/internal/upstream"
	"github.com/TheRootDaemon/tlgc/locale"
	"github.com/TheRootDaemon/tlgc/logger"
	"github.com/TheRootDaemon/tlgc/platform"
)

func main() {
	cli, err := cmd.Parse()
	if err != nil {
		logger.Error("%w", err)
		os.Exit(1)
	}

	logger.SetDefault(
		logger.New(
			cli.Quiet,
			cli.Verbose,
		),
	)

	os.Exit(run(cli))
}

func resolveLanguages(flagLangs []string) []string {
	if len(flagLangs) > 0 {
		return flagLangs
	}

	if cfgLangs := config.Cache().Languages; len(cfgLangs) > 0 {
		return cfgLangs
	}

	var langs []string
	locale.GetLanguages(&langs)
	return langs
}

func resolvePlatform(flagPlatform string) string {
	if flagPlatform != "" {
		return platform.Resolve(flagPlatform)
	}
	return platform.Default()
}

func renderOptions(cli *cmd.CLI) []render.RenderOption {
	opts := []render.RenderOption{render.WithWriter(os.Stdout)}

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

func run(cli *cmd.CLI) int {
	switch {
	case cli.Update:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		c := cache.New()
		languages := resolveLanguages(cli.Languages)

		client := upstream.New()
		if err := c.Update(
			context.Background(),
			languages,
			client,
		); err != nil {
			logger.Error("failed to update cache: %v", err)
			return 1
		}
		return 0

	case cli.List:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		c := cache.New()
		p := resolvePlatform(cli.Platform)
		pages, err := c.ListFor(p)
		if err != nil {
			logger.Error("failed to list pages: %v", err)
			return 1
		}

		for _, page := range pages {
			fmt.Println(page)
		}
		return 0

	case cli.ListAll:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		c := cache.New()
		pages, err := c.ListAll()
		if err != nil {
			logger.Error("failed to list pages: %v", err)
			return 1
		}

		for _, page := range pages {
			fmt.Println(page)
		}
		return 0

	case cli.Search != "":
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		p := resolvePlatform(cli.Platform)
		languages := resolveLanguages(cli.Languages)

		c := cache.New()
		results, err := c.Search(cli.Search, p, languages)
		if err != nil {
			logger.Error("search failed: %v", err)
			return 1
		}

		for _, r := range results {
			fmt.Printf("%s/%s\n", r.Platform, r.Page)
		}
		return 0

	case cli.ListPlatforms:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		c := cache.New()
		platforms, err := c.ListPlatforms()
		if err != nil {
			logger.Error("failed to list platforms: %v", err)
			return 1
		}

		for _, p := range platforms {
			fmt.Println(p)
		}
		return 0

	case cli.ListLanguages:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		c := cache.New()
		languages, err := c.ListLanguages()
		if err != nil {
			logger.Error("failed to list languages: %v", err)
			return 1
		}

		for _, l := range languages {
			fmt.Println(l)
		}
		return 0

	case cli.Info:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		c := cache.New()
		info, err := c.Info()
		if err != nil {
			logger.Error("failed to get cache info: %v", err)
			return 1
		}

		fmt.Printf("Cache directory: %s\n", info.CacheDir)
		fmt.Printf("Cache age: %s\n", info.Age)
		fmt.Printf("Total pages: %d\n", info.TotalPages)
		fmt.Printf("Auto update: %v\n", info.AutoUpdate)
		fmt.Printf("Max age (hours): %d\n", info.MaxAge)
		for _, ls := range info.LanguageStats {
			fmt.Printf("  %s: %d pages\n", ls.Language, ls.Pages)
		}
		return 0

	case cli.GenConfig:
		cfg, err := config.DefaultConfig()
		if err != nil {
			logger.Error(
				"failed to generate config: %w",
				err,
			)
		}
		fmt.Print(cfg)
		return 0

	case cli.ConfigPath:
		path := config.ConfigPath()
		fmt.Println(path)
		return 0

	case len(cli.Page) > 0:
		if err := config.Initialize(); err != nil {
			logger.Error("failed to load config: %v", err)
			return 1
		}

		p := resolvePlatform(cli.Platform)
		langs := resolveLanguages(cli.Languages)
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

		renderer := render.New(os.Stdout, renderOptions(cli)...)
		if err := renderer.Render(p, page); err != nil {
			logger.Error("failed to render page: %v", err)
			return 1
		}

		return 0

	default:
		return 0
	}
}
