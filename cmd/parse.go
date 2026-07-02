package cmd

import (
	"flag"
	"os"
)

func Parse() (*CLI, error) {
	return parse(os.Args[1:])
}

func parse(args []string) (*CLI, error) {
	cli := &CLI{}

	fs := flag.NewFlagSet("tlgc", flag.ContinueOnError)

	// options
	fs.BoolVar(&cli.Quiet, "q", false, "")
	fs.BoolVar(&cli.Quiet, "quiet", false, "")

	fs.Var(
		&countValue{
			count: &cli.Verbose,
		},
		"verbose",
		"",
	)

	// operations
	fs.BoolVar(
		&cli.ConfigPath,
		"config-path",
		false,
		"print the configuration path",
	)

	fs.BoolVar(
		&cli.GenConfig,
		"gen-config",
		false,
		"print the default configuration",
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return cli, nil
}
