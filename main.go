package main

import (
	"context"
	"fmt"
	"os"

	"github.com/portainer/compose-unpacker/commands"
	"github.com/portainer/compose-unpacker/exec"
	"github.com/portainer/compose-unpacker/log"
	"github.com/portainer/portainer/pkg/fips"

	"github.com/alecthomas/kong"
)

const UNPACKER_EXIT_ERROR = 255

var cli struct {
	LogLevel      log.Level                     `kong:"help='Set the logging level',default='INFO',enum='DEBUG,INFO,WARN,ERROR',env='LOG_LEVEL'"`
	PrettyLog     bool                          `kong:"help='Whether to enable or disable colored logs output',default='false',env='PRETTY_LOG'"`
	Deploy        commands.DeployCommand        `cmd:"" help:"Deploy a stack from a Git repository."`
	Undeploy      commands.UndeployCommand      `cmd:"" help:"Remove a stack from a Git repository."`
	SwarmDeploy   commands.SwarmDeployCommand   `cmd:"" help:"Deploy a Swarm stack from a Git repository."`
	SwarmUndeploy commands.SwarmUndeployCommand `cmd:"" help:"Remove a Swarm stack from a Git repository."`
	RemoveDir     commands.RemoveDirCommand     `cmd:"" help:"Remove a directory."`
	FipsMode      bool                          `kong:"help='Start in FIPS mode',name='fips-mode',env='FIPS_MODE',default:'false'"`
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

	fips.InitFIPS(cli.FipsMode)

	log.ConfigureLogger(cli.PrettyLog)
	log.SetLoggingLevel(cli.LogLevel)

	cmdCtx := exec.NewCommandExecutionContext(ctx)
	if err := cliCtx.Run(cmdCtx); err != nil {
		fmt.Println(err)
		os.Exit(UNPACKER_EXIT_ERROR)
	}
}
