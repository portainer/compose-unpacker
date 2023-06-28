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
		Strs("composePath", cmd.ComposeRelativeFilePaths).
		Msg("Undeploying Compose stack from Git repository")

	i := strings.LastIndex(cmd.GitRepository, "/")
	if i == -1 {
		log.Error().
			Str("repository", cmd.GitRepository).
			Msg("Invalid Git repository URL")

		return errDeployComposeFailure
	}

	mountPath := makeWorkingDir(cmd.Destination, cmd.ProjectName)

	deployer, err := compose.NewComposeDeployer(BIN_PATH, PORTAINER_DOCKER_CONFIG_PATH)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Compose deployer")

		return errDeployComposeFailure
	}

	log.Debug().
		Str("projectName", cmd.ProjectName).
		Msg("Undeploying Compose stack")

	err = deployer.Remove(cmdCtx.context, cmd.ProjectName, nil, libstack.Options{})
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to remove Compose stack")
		return errDeployComposeFailure
	}

	log.Info().Msg("Compose stack remove complete")

	if !cmd.Keep { //stack stop request
		err = os.RemoveAll(mountPath)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to remove Compose stack project folder")
		}
	}

	return nil
}

func (cmd *SwarmUndeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	log.Info().
		Str("stack name", cmd.ProjectName).
		Str("destination", cmd.Destination).
		Msg("Undeploying Swarm stack from Git repository")

	command := path.Join(BIN_PATH, "docker")
	if runtime.GOOS == "windows" {
		command = path.Join(BIN_PATH, "docker.exe")
	}

	args := make([]string, 0)
	args = append(args, "stack", "rm", cmd.ProjectName)
	err := runCommandAndCaptureStdErr(command, args, nil, "")
	if err != nil {
		return err
	}

	mountPath := makeWorkingDir(cmd.Destination, cmd.ProjectName)
	if !cmd.Keep { //stack stop request
		err = os.RemoveAll(mountPath)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to remove Compose stack project folder")
		}
	}

	return nil
}
