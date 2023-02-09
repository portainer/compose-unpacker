package main

import (
	"os"
	"path"
	"runtime"
	"strings"

	libstack "github.com/portainer/docker-compose-wrapper"
	"github.com/portainer/docker-compose-wrapper/compose"
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
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")
	clonePath := path.Join(mountPath, repositoryName)
	log.Debug().
		Str("path", clonePath).
		Str("binPath", BIN_PATH).
		Msg("Creating Compose deployer")

	deployer, err := compose.NewComposeDeployer(BIN_PATH, "")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Compose deployer")

		return errDeployComposeFailure
	}

	composeFilePaths := make([]string, len(cmd.ComposeRelativeFilePaths))
	for i := 0; i < len(cmd.ComposeRelativeFilePaths); i++ {
		composeFilePaths[i] = path.Join(clonePath, cmd.ComposeRelativeFilePaths[i])
	}

	log.Debug().
		Strs("composeFilePaths", composeFilePaths).
		Str("workingDirectory", clonePath).
		Str("projectName", cmd.ProjectName).
		Msg("Undeploying Compose stack")

	err = deployer.Remove(cmdCtx.context, cmd.ProjectName, nil, libstack.Options{
		WorkingDir: clonePath,
	})
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
