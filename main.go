package main

import (
	"context"
	"log"

	"github.com/alecthomas/kong"
	"go.uber.org/zap"
)

func initializeLogger(debug bool) (*zap.SugaredLogger, error) {
	if debug {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		return logger.Sugar(), nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func main() {
	ctx := context.Background()

	cliCtx := kong.Parse(&cli,
		kong.Name("unpacker"),
		kong.Description("A tool to deploy Docker stacks from Git repositories."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	logger, err := initializeLogger(cli.Debug)
	if err != nil {
		log.Fatalf("Unable to initialize logger: %s", err)
	}

	cmdCtx := NewCommandExecutionContext(ctx, logger)
	err = cliCtx.Run(cmdCtx)
	cliCtx.FatalIfErrorf(err)
}
