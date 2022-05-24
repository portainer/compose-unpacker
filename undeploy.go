package main

import (
	"errors"
	"github.com/portainer/docker-compose-wrapper/compose"
	"os"
	"path"
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
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")
	clonePath := path.Join(cmd.Destination, repositoryName)
	defer func() {
		err := os.RemoveAll(cmd.Destination)
		if nil != err {
			cmdCtx.logger.Errorw("Failed to delete the stack folder",
				"targetFolder", clonePath,
			)
		}
	}()

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
	//Remove(ctx context.Context, workingDir, host, projectName string, filePaths []string) error
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
