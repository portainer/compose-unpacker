package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	libstack "github.com/portainer/docker-compose-wrapper"
	"github.com/portainer/docker-compose-wrapper/compose"
	"github.com/rs/zerolog/log"
)

var errDeployComposeFailure = errors.New("stack deployment failure")
var errDockerLoginFailure = errors.New("docker login failure")
var errDockerLogoutFailure = errors.New("docker logout failure")

func (cmd *DeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	log.Info().
		Str("repository", cmd.GitRepository).
		Strs("composePath", cmd.ComposeRelativeFilePaths).
		Str("destination", cmd.Destination).
		Strs("env", cmd.Env).
		Bool("skipTLSVerify", cmd.SkipTLSVerify).
		Msg("Deploying Compose stack from Git repository")

	defer dockerLogout(cmd.Registry)
	err := dockerLogin(cmd.Registry)
	if err != nil {
		return err
	}

	if cmd.User != "" && cmd.Password != "" {
		log.Info().
			Str("user", cmd.User).
			Msg("Using Git authentication")
	}

	i := strings.LastIndex(cmd.GitRepository, "/")
	if i == -1 {

		log.Error().
			Str("repository", cmd.GitRepository).
			Msg("Invalid Git repository URL")
		return errDeployComposeFailure
	}
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")

	log.Info().
		Str("directory", cmd.Destination).
		Msg("Checking the file system...")

	mountPath := makeWorkingDir(cmd.Destination, cmd.ProjectName)
	clonePath := path.Join(mountPath, repositoryName)
	if !cmd.Keep { //stack create request
		_, err := os.Stat(mountPath)
		if err == nil {
			err = os.RemoveAll(mountPath)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed to remove previous directory")
				return errDeployComposeFailure
			}
		}

		err = os.MkdirAll(mountPath, 0755)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to create destination directory")
			return errDeployComposeFailure
		}

		log.Info().
			Str("directory", mountPath).
			Msg("Creating target destination directory on disk")

		gitOptions := git.CloneOptions{
			URL:             cmd.GitRepository,
			ReferenceName:   plumbing.ReferenceName(cmd.Reference),
			Auth:            getAuth(cmd.User, cmd.Password),
			Depth:           1,
			InsecureSkipTLS: cmd.SkipTLSVerify,
		}

		log.Info().
			Str("repository", cmd.GitRepository).
			Str("path", clonePath).
			Str("url", gitOptions.URL).
			Int("depth", gitOptions.Depth).
			Msg("Cloning git repository")

		_, err = git.PlainCloneContext(cmdCtx.context, clonePath, false, &gitOptions)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to clone Git repository")
			return errDeployComposeFailure
		}
	}

	deployer, err := compose.NewComposeDeployer(BIN_PATH, PORTAINER_DOCKER_CONFIG_PATH)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create Compose deployer")
		return errDeployComposeFailure
	}

	composeFilePaths := make([]string, len(cmd.ComposeRelativeFilePaths))
	for i := 0; i < len(cmd.ComposeRelativeFilePaths); i++ {
		composeFilePaths[i] = path.Join(clonePath, cmd.ComposeRelativeFilePaths[i])
	}

	log.Info().
		Strs("composeFilePaths", composeFilePaths).
		Str("workingDirectory", clonePath).
		Str("projectName", cmd.ProjectName).
		Msg("Deploying Compose stack")

	err = deployer.Deploy(cmdCtx.context, composeFilePaths, libstack.DeployOptions{
		Options: libstack.Options{
			WorkingDir:  clonePath,
			ProjectName: cmd.ProjectName,
			Env:         cmd.Env,
		},
		ForceRecreate: true,
	})

	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to deploy Compose stack")
		return errDeployComposeFailure
	}

	log.Info().Msg("Compose stack deployment complete")
	return nil
}

func (cmd *SwarmDeployCommand) Run(cmdCtx *CommandExecutionContext) error {
	log.Info().
		Str("repository", cmd.GitRepository).
		Strs("composePath", cmd.ComposeRelativeFilePaths).
		Str("destination", cmd.Destination).
		Msg("Deploying Swarm stack from a Git repository")

	defer dockerLogout(cmd.Registry)
	err := dockerLogin(cmd.Registry)
	if err != nil {
		return err
	}

	if cmd.User != "" && cmd.Password != "" {
		log.Info().
			Str("user", cmd.User).
			Msg("Using Git authentication")
	}

	i := strings.LastIndex(cmd.GitRepository, "/")
	if i == -1 {
		log.Error().
			Str("repository", cmd.GitRepository).
			Msg("Invalid Git repository URL")

		return errDeployComposeFailure
	}
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")

	log.Info().
		Str("directory", cmd.Destination).
		Msg("Checking the file system...")

	mountPath := makeWorkingDir(cmd.Destination, cmd.ProjectName)
	clonePath := path.Join(mountPath, repositoryName)

	// Record running services before deployment/redeployment
	serviceIDs, err := checkRunningService(cmd.ProjectName)
	if err != nil {
		return err
	}

	runningServices := make(map[string]struct{}, 0)
	for _, serviceID := range serviceIDs {
		runningServices[serviceID] = struct{}{}
	}

	forceUpdate := false
	if len(runningServices) > 0 {
		// To determine whether the current service needs to force update, it
		// is more reliable to check if there is a created service with the
		// stack name rather than to check if there is an existing git repository.
		forceUpdate = true
		log.Info().Msg("Set to force update")
	}

	if !cmd.Keep { //stack create request
		_, err := os.Stat(mountPath)
		if err == nil {
			err = os.RemoveAll(mountPath)
			if err != nil {
				log.Error().
					Err(err).
					Msg("Failed to remove previous directory")
				return errDeployComposeFailure
			}
		}
		err = os.MkdirAll(mountPath, 0755)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to create destination directory")
			return errDeployComposeFailure
		}

		log.Info().
			Str("directory", mountPath).
			Msg("Creating target destination directory on disk")

		gitOptions := git.CloneOptions{
			URL:             cmd.GitRepository,
			ReferenceName:   plumbing.ReferenceName(cmd.Reference),
			Auth:            getAuth(cmd.User, cmd.Password),
			Depth:           100,
			InsecureSkipTLS: cmd.SkipTLSVerify,
		}

		log.Info().
			Str("repository", cmd.GitRepository).
			Str("path", clonePath).
			Str("url", gitOptions.URL).
			Int("depth", gitOptions.Depth).
			Msg("Cloning git repository")

		_, err = git.PlainCloneContext(cmdCtx.context, clonePath, false, &gitOptions)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to clone Git repository")
			return errDeployComposeFailure
		}
	}

	err = deploySwarmStack(*cmd, clonePath)
	if err != nil {
		return err
	}

	if forceUpdate {
		// If the process executes redeployment, the running services need
		// to be recreated forcibly
		updatedServiceIDs, err := checkRunningService(cmd.ProjectName)
		if err != nil {
			return err
		}

		for _, updatedServiceID := range updatedServiceIDs {
			_, ok := runningServices[updatedServiceID]
			if ok {
				_ = updateService(updatedServiceID)
			}
		}
	}

	return nil
}

func dockerLogin(registries []string) error {
	command := getDockerBinaryPath()

	for _, registry := range registries {
		credentials := strings.Split(registry, ":")
		if len(credentials) != 3 {
			return errDockerLoginFailure
		}

		args := make([]string, 0)
		args = append(args, "--config", PORTAINER_DOCKER_CONFIG_PATH, "login", "--username", credentials[0], "--password", credentials[1], credentials[2])

		err := runCommandAndCaptureStdErr(command, args, nil, "")
		if err != nil {
			log.Error().
				Err(err).
				Msg(fmt.Sprintf("Docker login %s failed", credentials[2]))
			return errDockerLoginFailure
		}
		log.Info().
			Msg(fmt.Sprintf("Docker login %s successed", credentials[2]))
	}

	return nil
}

func dockerLogout(registries []string) error {
	command := getDockerBinaryPath()

	for _, registry := range registries {
		credentials := strings.Split(registry, ":")
		if len(credentials) != 3 {
			return errDockerLogoutFailure
		}

		args := make([]string, 0)
		args = append(args, "--config", PORTAINER_DOCKER_CONFIG_PATH, "logout", credentials[2])

		err := runCommandAndCaptureStdErr(command, args, nil, "")
		if err != nil {
			log.Error().
				Err(err).
				Msg(fmt.Sprintf("Docker logout %s failed", credentials[2]))
			return errDockerLogoutFailure
		}
		log.Info().
			Msg(fmt.Sprintf("Docker logout %s successed", credentials[2]))
	}

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
		return errors.New(stderr.String())
	}
	return nil
}

func runCommand(command string, args []string) (string, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
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

func makeWorkingDir(target, stackName string) string {
	return filepath.Join(target, "stacks", stackName)
}

func getDockerBinaryPath() string {
	command := path.Join(BIN_PATH, "docker")
	if runtime.GOOS == "windows" {
		command = path.Join(BIN_PATH, "docker.exe")
	}
	return command
}
