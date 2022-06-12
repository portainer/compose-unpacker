package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	fmt.Println("Unpacker begin to work")
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
	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}
	os.Exit(99)
}
