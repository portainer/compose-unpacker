package main

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/portainer/portainer/pkg/libstack"
	"github.com/portainer/portainer/pkg/libstack/compose"
	"github.com/rs/zerolog/log"
)

func (cmd *UndeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	log.Info().
		Str("repository", cmd.GitRepository).
		Strs("compose_path", cmd.ComposeRelativeFilePaths).
		Msg("Undeploying Compose stack from Git repository")

	if strings.LastIndex(cmd.GitRepository, "/") == -1 {
		log.Error().
			Str("repository", cmd.GitRepository).
			Msg("Invalid Git repository URL")

		return errDeployComposeFailure
	}

	mountPath := makeWorkingDir(cmd.Destination, cmd.ProjectName)

	deployer := compose.NewComposeDeployer()

	log.Debug().
		Str("project_name", cmd.ProjectName).
		Bool("remove_volumes", cmd.RemoveVolumes).
		Msg("Undeploying Compose stack")

	if err := deployer.Remove(cmdCtx.context, cmd.ProjectName, nil, libstack.RemoveOptions{Volumes: cmd.RemoveVolumes}); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to remove Compose stack")
		return errDeployComposeFailure
	}

	log.Info().Msg("Compose stack remove complete")

	if !cmd.Keep { //stack stop request
		if err := os.RemoveAll(mountPath); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to remove Compose stack project folder")
		}
	}

	return nil
}

func (cmd *SwarmUndeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	log.Info().
		Str("stack_name", cmd.ProjectName).
		Str("destination", cmd.Destination).
		Msg("Undeploying Swarm stack from Git repository")

	command := path.Join(BIN_PATH, "docker")
	if runtime.GOOS == "windows" {
		command = path.Join(BIN_PATH, "docker.exe")
	}

	args := make([]string, 0)
	args = append(args, "stack", "rm", "--detach=false", cmd.ProjectName)
	if err := runCommandAndCaptureStdErr(command, args, nil, ""); err != nil {
		return err
	}

	mountPath := makeWorkingDir(cmd.Destination, cmd.ProjectName)
	if !cmd.Keep { //stack stop request
		if err := os.RemoveAll(mountPath); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to remove Compose stack project folder")
		}
	}

	return nil
}
