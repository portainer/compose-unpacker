package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/portainer/docker-compose-wrapper/compose"
	"github.com/portainer/portainer/api/filesystem"
)

var errDeployComposeFailure = errors.New("stack deployment failure")

func (cmd *DeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	cmdCtx.logger.Infow("Deploying Compose stack from Git repository",
		"repository", cmd.GitRepository,
		"composePath", cmd.ComposeRelativeFilePaths,
		"destination", cmd.Destination,
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

	cmdCtx.logger.Infow("Checking the file system...",
		"directory", cmd.Destination,
	)
	mountPath := fmt.Sprintf("%s/%s", cmd.Destination, "mount")
	if _, err := os.Stat(mountPath); err != nil {
		if os.IsNotExist(err) {
			cmdCtx.logger.Infow("Creating folder in the file system...",
				"directory", mountPath,
			)
			err := os.MkdirAll(mountPath, 0755)
			if err != nil {
				cmdCtx.logger.Errorw("Failed to create destination directory",
					"error", err,
				)
				return errDeployComposeFailure
			}
		} else {
			return err
		}
	} else {
		cmdCtx.logger.Infow("Backing up folder in the file system...",
			"directory", mountPath,
		)
		backupProjectPath := fmt.Sprintf("%s-old", mountPath)
		err = filesystem.MoveDirectory(mountPath, backupProjectPath)
		if err != nil {
			return err
		}
		defer func() {
			err = os.RemoveAll(backupProjectPath)
			if err != nil {
				log.Printf("[WARN] [http,stacks,git] [error: %s] [message: unable to remove git repository directory]", err)
			}
		}()
	}
	cmdCtx.logger.Infow("Creating target destination directory on disk",
		"directory", mountPath,
	)
	gitOptions := git.CloneOptions{
		URL:   cmd.GitRepository,
		Auth:  getAuth(cmd.User, cmd.Password),
		Depth: 1,
	}

	clonePath := path.Join(mountPath, repositoryName)

	cmdCtx.logger.Infow("Cloning git repository",
		"path", clonePath,
		"cloneOptions", gitOptions,
	)

	_, err := git.PlainCloneContext(cmdCtx.context, clonePath, false, &gitOptions)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to clone Git repository",
			"error", err,
		)

		return errDeployComposeFailure
	}

	cmdCtx.logger.Infow("Creating Compose deployer",
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

	cmdCtx.logger.Infow("Deploying Compose stack",
		"composeFilePaths", composeFilePaths,
		"workingDirectory", clonePath,
		"projectName", cmd.ProjectName,
	)

	err = deployer.Deploy(cmdCtx.context, clonePath, "", cmd.ProjectName, composeFilePaths, "", false)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to deploy Compose stack",
			"error", err,
		)

		return errDeployComposeFailure
	}

	cmdCtx.logger.Info("Compose stack deployment complete")

	return nil
}

func (cmd *SwarmDeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	cmdCtx.logger.Infow("Deploying Swarm stack from Git repository",
		"repository", cmd.GitRepository,
		"composePath", cmd.ComposeRelativeFilePaths,
		"destination", cmd.Destination,
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

	cmdCtx.logger.Infow("Checking the file system...",
		"directory", cmd.Destination,
	)
	mountPath := fmt.Sprintf("%s/%s", cmd.Destination, "mount")
	if _, err := os.Stat(mountPath); err != nil {
		if os.IsNotExist(err) {
			cmdCtx.logger.Infow("Creating folder in the file system...",
				"directory", mountPath,
			)
			err := os.MkdirAll(mountPath, 0755)
			if err != nil {
				cmdCtx.logger.Errorw("Failed to create destination directory",
					"error", err,
				)
				return errDeployComposeFailure
			}
		} else {
			return err
		}
	} else {
		cmdCtx.logger.Infow("Backing up folder in the file system...",
			"directory", mountPath,
		)
		backupProjectPath := fmt.Sprintf("%s-old", mountPath)
		err = filesystem.MoveDirectory(mountPath, backupProjectPath)
		if err != nil {
			return err
		}
		defer func() {
			err = os.RemoveAll(backupProjectPath)
			if err != nil {
				log.Printf("[WARN] [http,stacks,git] [error: %s] [message: unable to remove git repository directory]", err)
			}
		}()
	}
	cmdCtx.logger.Infow("Creating target destination directory on disk",
		"directory", mountPath,
	)
	gitOptions := git.CloneOptions{
		URL:   cmd.GitRepository,
		Auth:  getAuth(cmd.User, cmd.Password),
		Depth: 100,
	}

	clonePath := path.Join(mountPath, repositoryName)

	cmdCtx.logger.Infow("Cloning git repository",
		"path", clonePath,
		"cloneOptions", gitOptions,
	)

	_, err := git.PlainCloneContext(cmdCtx.context, clonePath, false, &gitOptions)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to clone Git repository",
			"error", err,
		)

		return errDeployComposeFailure
	}

	command := path.Join(BIN_PATH, "docker")

	if runtime.GOOS == "windows" {
		command = path.Join(BIN_PATH, "docker.exe")
	}
	args := make([]string, 0)

	if cmd.Prune {
		args = append(args, "stack", "deploy", "--prune", "--with-registry-auth")
	} else {
		args = append(args, "stack", "deploy", "--with-registry-auth")
	}
	if !cmd.Pull {
		args = append(args, "--resolve-image=never")
	}

	for _, cfile := range cmd.ComposeRelativeFilePaths {
		args = append(args, "--compose-file", path.Join(clonePath, cfile))
	}
	cmdCtx.logger.Infow("Deploying Swarm stack",
		"composeFilePaths", cmd.ComposeRelativeFilePaths,
		"workingDirectory", clonePath,
		"projectName", cmd.ProjectName,
	)
	args = append(args, cmd.ProjectName)

	env := make([]string, 0)
	/*
		for _, envvar := range stack.Env {
			env = append(env, envvar.Name+"="+envvar.Value)
		}

	*/
	err = runCommandAndCaptureStdErr(command, args, env, cmd.ProjectName)
	if err != nil {
		cmdCtx.logger.Errorw("Failed to swarm deplot Git repository",
			"error", err,
		)

		return errDeployComposeFailure
	}
	cmdCtx.logger.Info("Swarm stack deployment complete")

	return nil
}

func runCommandAndCaptureStdErr(command string, args []string, env []string, workingDir string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Dir = workingDir

	if env != nil {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, env...)
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return errors.New(stderr.String())
	}

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

func prepareDockerCommandAndArgs(binaryPath, configPath string) (string, []string, error) {
	// Assume Linux as a default
	command := path.Join(binaryPath, "docker")

	if runtime.GOOS == "windows" {
		command = path.Join(binaryPath, "docker.exe")
	}

	args := make([]string, 0)
	args = append(args, "--config", configPath)

	return command, args, nil
}
