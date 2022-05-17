package main

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/portainer/docker-compose-wrapper/compose"
)

type DeployCommand struct {
	User     string `help:"Username for Git authentication." short:"u"`
	Password string `help:"Password or PAT for Git authentication" short:"p"`

	GitRepository           string `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	ComposeRelativeFilePath string `arg:"" help:"Relative path to the Compose file." type:"path" name:"compose-file-path"`
	ProjectName             string `arg:"" help:"Name of the Compose stack." name:"project-name"`
	Destination             string `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
}

var errDeployComposeFailure = errors.New("compose stack deployment failure")

func (cmd *DeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	cmdCtx.logger.Infow("Deploying Compose stack from Git repository",
		"repository", cmd.GitRepository,
		"composePath", cmd.ComposeRelativeFilePath,
	)

	if cmd.User != "" && cmd.Password != "" {
		cmdCtx.logger.Infow("Using Git authentication",
			"user", cmd.User,
			"password", "<redacted>",
		)
	}

	i := strings.LastIndex(cmd.GitRepository, "/")
	if i == -1 {
		cmdCtx.logger.Errorw("Invalid Git repository URL",
			"repository", cmd.GitRepository,
		)

		return errDeployComposeFailure
	}
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")

	cmdCtx.logger.Debugw("Creating target destination directory on disk",
		"directory", cmd.Destination,
	)

	err := os.MkdirAll(cmd.Destination, 0755)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to create destination directory",
			"error", err,
		)

		return errDeployComposeFailure
	}

	gitOptions := git.CloneOptions{
		URL:   cmd.GitRepository,
		Auth:  getAuth(cmd.User, cmd.Password),
		Depth: 1,
	}

	clonePath := path.Join(cmd.Destination, repositoryName)

	cmdCtx.logger.Debugw("Cloning git repository",
		"path", clonePath,
		"cloneOptions", gitOptions,
	)

	_, err = git.PlainCloneContext(cmdCtx.context, clonePath, false, &gitOptions)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to clone Git repository",
			"error", err,
		)

		return errDeployComposeFailure
	}

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

	composeFilePath := path.Join(clonePath, cmd.ComposeRelativeFilePath)

	cmdCtx.logger.Debugw("Deploying Compose stack",
		"composeFilePath", composeFilePath,
		"workingDirectory", clonePath,
		"projectName", cmd.ProjectName,
	)

	err = deployer.Deploy(cmdCtx.context, clonePath, "", cmd.ProjectName, []string{composeFilePath}, "", false)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to deploy Compose stack",
			"error", err,
		)

		return errDeployComposeFailure
	}

	cmdCtx.logger.Info("Compose stack deployment complete")

	return nil
}

func getAuth(username, password string) *http.BasicAuth {
	if password != "" {
		if username == "" {
			username = "token"
		}

		return &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	return nil
}
