package main

import (
	"errors"
	"fmt"
	"github.com/portainer/docker-compose-wrapper/compose"
	"path"
	"runtime"
	"strings"
)

var errUndeployComposeFailure = errors.New("compose stack remove failure")

func (cmd *UndeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	cmdCtx.logger.Infow("Undeploying Compose stack from Git repository",
		"repository", cmd.GitRepository,
		"composePath", cmd.ComposeRelativeFilePaths,
	)

	i := strings.LastIndex(cmd.GitRepository, "/")
	if i == -1 {
		cmdCtx.logger.Errorw("Invalid Git repository URL",
			"repository", cmd.GitRepository,
		)

		return errDeployComposeFailure
	}
	mountPath := fmt.Sprintf("%s/%s", cmd.Destination, "mount")
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")
	clonePath := path.Join(mountPath, repositoryName)
	cmdCtx.logger.Debugw("Current git repository",
		"path", clonePath,
	)

	cmdCtx.logger.Debugw("Creating Compose deployer",
		"binPath", BIN_PATH,
	)
	deployer, err := compose.NewComposeDeployer(BIN_PATH, "")
	if err != nil {
		cmdCtx.logger.Errorw("Failed to create Compose deployer",
			"error", err,
		)

		return errDeployComposeFailure
	}
	composeFilePaths := make([]string, len(cmd.ComposeRelativeFilePaths))
	for i := 0; i < len(cmd.ComposeRelativeFilePaths); i++ {
		composeFilePaths[i] = path.Join(clonePath, cmd.ComposeRelativeFilePaths[i])
	}
	cmdCtx.logger.Debugw("Undeploying Compose stack",
		"composeFilePaths", composeFilePaths,
		"workingDirectory", clonePath,
		"projectName", cmd.ProjectName,
	)
	err = deployer.Remove(cmdCtx.context, clonePath, "", cmd.ProjectName, composeFilePaths)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to remove Compose stack",
			"error", err,
		)
		return errDeployComposeFailure
	}
	cmdCtx.logger.Info("Compose stack remove complete")
	return nil
}

func (cmd *SwarmUndeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	cmdCtx.logger.Infow("Undeploying Swarm stack from Git repository",
		"stack name", cmd.ProjectName,
		"destination", cmd.Destination,
	)
	command := path.Join(BIN_PATH, "docker")
	if runtime.GOOS == "windows" {
		command = path.Join(BIN_PATH, "docker.exe")
	}
	args := make([]string, 0)
	args = append(args, "stack", "rm", cmd.ProjectName)
	return runCommandAndCaptureStdErr(command, args, nil, "")
}
