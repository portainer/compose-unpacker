package commands

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/portainer/compose-unpacker/auth"
	"github.com/portainer/compose-unpacker/exec"
	"github.com/rs/zerolog/log"
)

type SwarmDeployCommand struct {
	User                     string   `help:"Username for Git authentication." short:"u"`
	Password                 string   `help:"Password or PAT for Git authentication" short:"p"`
	Pull                     bool     `help:"Pull Image" short:"f"`
	Prune                    bool     `help:"Prune services during deployment" short:"r"`
	Keep                     bool     `help:"Keep stack folder" short:"k"`
	SkipTLSVerify            bool     `help:"Skip TLS verification for git" name:"skip-tls-verify"`
	ForceRecreateStack       bool     `help:"Force to recreate the target stack regardless whether the image hash changes" name:"force-recreate"`
	Env                      []string `help:"OS ENV for stack."`
	Registry                 []string `help:"Registry credentials" name:"registry"`
	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	Reference                string   `arg:"" help:"Reference of Git repository to deploy from." name:"git-ref"`
	ProjectName              string   `arg:"" help:"Name of the Swarm stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file."  name:"compose-file-paths"`
}

func (cmd *SwarmDeployCommand) Run(cmdCtx *exec.CommandExecutionContext) error {
	log.Info().
		Str("repository", cmd.GitRepository).
		Strs("composePath", cmd.ComposeRelativeFilePaths).
		Str("destination", cmd.Destination).
		Msg("Deploying Swarm stack from a Git repository")

	if err := auth.DockerLogin(cmd.Registry); err != nil {
		return fmt.Errorf("an error occured in swarm docker login. Error: %w", err)
	}
	defer auth.DockerLogout(cmd.Registry)

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

		return exec.ErrDeployComposeFailure
	}
	repositoryName := strings.TrimSuffix(cmd.GitRepository[i+1:], ".git")

	log.Info().
		Str("directory", cmd.Destination).
		Msg("Checking the file system...")

	mountPath := exec.MakeWorkingDir(cmd.Destination, cmd.ProjectName)
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
	if cmd.ForceRecreateStack && len(runningServices) > 0 {
		// To determine whether the current service needs to force update, it
		// is more reliable to check if there is a created service with the
		// stack name rather than to check if there is an existing git repository.
		forceUpdate = true
		log.Info().Msg("Set to force update")
	}

	if !cmd.Keep { // Stack create request
		if _, err := os.Stat(mountPath); err == nil {
			if err := os.RemoveAll(mountPath); err != nil {
				log.Error().
					Err(err).
					Msg("Failed to remove previous directory")
				return exec.ErrDeployComposeFailure
			}
		}

		if err := os.MkdirAll(mountPath, 0755); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to create destination directory")
			return exec.ErrDeployComposeFailure
		}

		log.Info().
			Str("directory", mountPath).
			Msg("Creating target destination directory on disk")

		gitOptions := git.CloneOptions{
			URL:             cmd.GitRepository,
			ReferenceName:   plumbing.ReferenceName(cmd.Reference),
			Auth:            auth.GetAuth(cmd.User, cmd.Password),
			Depth:           1,
			InsecureSkipTLS: cmd.SkipTLSVerify,
			Tags:            git.NoTags,
		}

		log.Info().
			Str("repository", cmd.GitRepository).
			Str("path", clonePath).
			Str("url", gitOptions.URL).
			Int("depth", gitOptions.Depth).
			Msg("Cloning git repository")

		if _, err = git.PlainCloneContext(cmdCtx.Context, clonePath, false, &gitOptions); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to clone Git repository")

			return exec.ErrDeployComposeFailure
		}
	}

	if err := deploySwarmStack(*cmd, clonePath); err != nil {
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
			if _, ok := runningServices[updatedServiceID]; ok {
				_ = updateService(updatedServiceID, forceUpdate)
			}
		}
	}

	return nil
}

func deploySwarmStack(cmd SwarmDeployCommand, clonePath string) error {
	command := exec.GetDockerBinaryPath()
	args := []string{"--config", exec.PORTAINER_DOCKER_CONFIG_PATH}

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
	log.Info().
		Strs("composeFilePaths", cmd.ComposeRelativeFilePaths).
		Str("workingDirectory", clonePath).
		Str("projectName", cmd.ProjectName).
		Msg("Deploying Swarm stack")

	args = append(args, cmd.ProjectName)

	err := exec.RunCommandAndCaptureStdErr(command, args, cmd.Env, clonePath)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to swarm deploy Git repository")
		return exec.ErrDeployComposeFailure
	}
	log.Info().
		Msg("Swarm stack deployment complete")

	return err
}

func checkRunningService(projectName string) ([]string, error) {
	command := exec.GetDockerBinaryPath()
	args := []string{"--config", exec.PORTAINER_DOCKER_CONFIG_PATH, "stack", "services", "--format={{.ID}}", projectName}

	log.Info().
		Strs("args", args).
		Msg("Checking Swarm stack")

	output, err := exec.RunCommand(command, args)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to check running swarm services")
		return nil, err
	}

	serviceIDs := splitLines(output)
	log.Info().
		Strs("serviceIDs", serviceIDs).
		Msg("Checking stack services")
	return serviceIDs, nil
}

func updateService(serviceID string, forceRecreate bool) error {
	command := exec.GetDockerBinaryPath()
	args := []string{"--config", exec.PORTAINER_DOCKER_CONFIG_PATH, "service", "update", serviceID}
	if forceRecreate {
		args = append(args, "--force")
	}

	log.Info().
		Strs("args", args).
		Msg("Updating Swarm service")

	out, err := exec.RunCommand(command, args)
	if err != nil {
		log.Error().
			Err(err).
			Str("standard_output", out).
			Str("context", "SwarmDeployerUpdateService").
			Msg("Failed to update swarm services")
		return err
	}

	log.Info().
		Str("standard_output", out).
		Str("context", "SwarmDeployerUpdateService").
		Msg("Update stack service completed")
	return nil
}

func splitLines(s string) []string {
	var separator string
	if runtime.GOOS == "windows" {
		separator = "\r\n"
	} else {
		separator = "\n"
	}
	parts := strings.Split(s, separator)

	ret := []string{}
	for _, part := range parts {
		// remove empty string
		if part != "" {
			ret = append(ret, part)
		}
	}
	return ret
}
