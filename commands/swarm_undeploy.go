package commands

import (
	"os"

	"github.com/portainer/compose-unpacker/exec"
	"github.com/portainer/portainer/pkg/libstack/swarm"
	"github.com/rs/zerolog/log"
)

type SwarmUndeployCommand struct {
	Keep        bool   `help:"Keep stack folder" short:"k"`
	ProjectName string `arg:"" help:"Name of the Compose (Swarm) stack." name:"project-name"`
	Destination string `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
}

func (cmd *SwarmUndeployCommand) Run(cmdCtx *exec.CommandExecutionContext) error {
	log.Info().
		Str("stack_name", cmd.ProjectName).
		Str("destination", cmd.Destination).
		Msg("Undeploying Swarm stack from Git repository")

	deployer := swarm.NewSwarmDeployer()

	if err := deployer.Remove(cmdCtx.Context, cmd.ProjectName, swarm.RemoveOptions{}); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to remove Swarm stack")
		return err
	}

	mountPath := exec.MakeWorkingDir(cmd.Destination, cmd.ProjectName)
	if !cmd.Keep { //stack stop request
		if err := os.RemoveAll(mountPath); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to remove Compose stack project folder")
		}
	}

	return nil
}
