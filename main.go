package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/portainer/compose-unpacker/log"
)

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

	log.ConfigureLogger(cli.PrettyLog)
	log.SetLoggingLevel(log.Level(cli.LogLevel))

	cmdCtx := NewCommandExecutionContext(ctx)
	err := cliCtx.Run(cmdCtx)
	if err != nil {
		fmt.Println(err)
		os.Exit(UNPACKER_EXIT_ERROR)
	}
}
