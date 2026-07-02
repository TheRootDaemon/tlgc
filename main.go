package main

import (
	"fmt"
	"os"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/config"
	"github.com/TheRootDaemon/tlgc/logger"
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

func run(cli *cmd.CLI) int {
	switch {
	case cli.ConfigPath:
		path := config.ConfigPath()
		fmt.Println(path)
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

	default:
		return 0
	}
}
