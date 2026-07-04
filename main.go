package main

import (
	"os"

	"github.com/TheRootDaemon/tlgc/cmd"
	"github.com/TheRootDaemon/tlgc/internal/app"
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

	os.Exit(
		app.New().Run(cli),
	)
}
