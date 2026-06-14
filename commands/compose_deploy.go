package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/portainer/compose-unpacker/auth"
	"github.com/portainer/compose-unpacker/exec"
	"github.com/portainer/portainer/api/filesystem"
	portainergit "github.com/portainer/portainer/api/git"
	"github.com/portainer/portainer/pkg/fips"
	"github.com/portainer/portainer/pkg/libstack"
	"github.com/portainer/portainer/pkg/libstack/compose"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	gogitfs "github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/rs/zerolog/log"
)

type DeployCommand struct {
	User                     string   `help:"Username for Git authentication." short:"u"`
	Password                 string   `help:"Password or PAT for Git authentication" short:"p"`
	Prune                    bool     `help:"Prune services during deployment" short:"r"`
	Keep                     bool     `help:"Keep stack folder" short:"k"`
	SkipTLSVerify            bool     `help:"Skip TLS verification for git" name:"skip-tls-verify"`
	ForceRecreateStack       bool     `help:"Force to recreate the target stack regardless whether the image hash changes" name:"force-recreate"`
	Env                      []string `help:"OS ENV for stack" example:"key=value"`
	Registry                 []string `help:"Registry credentials" name:"registry"`
	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	Reference                string   `arg:"" help:"Reference of Git repository to deploy from." name:"git-ref"`
	ProjectName              string   `arg:"" help:"Name of the Compose stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file."  name:"compose-file-paths"`
}

func (cmd *DeployCommand) Run(cmdCtx *exec.CommandExecutionContext) error {
	log.Info().
		Str("repository", cmd.GitRepository).
		Strs("composePath", cmd.ComposeRelativeFilePaths).
		Str("destination", cmd.Destination).
		Strs("env", cmd.Env).
		Bool("skipTLSVerify", cmd.SkipTLSVerify).
		Msg("Deploying Compose stack from Git repository")

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
	clonePath := filesystem.JoinPaths(mountPath, repositoryName)
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
			InsecureSkipTLS: cmd.SkipTLSVerify && fips.CanTLSSkipVerify(),
			Tags:            git.NoTags,
		}

		log.Info().
			Str("repository", cmd.GitRepository).
			Str("path", clonePath).
			Str("url", gitOptions.URL).
			Int("depth", gitOptions.Depth).
			Msg("Cloning git repository")

		wt := portainergit.NewNoSymlinkFS(osfs.New(clonePath))
		dot := osfs.New(filesystem.JoinPaths(clonePath, ".git"))
		storer := gogitfs.NewStorage(dot, cache.NewObjectLRU(0))

		if _, err := git.CloneContext(cmdCtx.Context, storer, wt, &gitOptions); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to clone Git repository")
			return exec.ErrDeployComposeFailure
		}
	}

	deployer := compose.NewComposeDeployer()

	composeFilePaths := make([]string, len(cmd.ComposeRelativeFilePaths))
	for i := range len(cmd.ComposeRelativeFilePaths) {
		composeFilePaths[i] = filesystem.JoinPaths(clonePath, cmd.ComposeRelativeFilePaths[i])
	}

	log.Info().
		Strs("composeFilePaths", composeFilePaths).
		Str("workingDirectory", clonePath).
		Str("projectName", cmd.ProjectName).
		Msg("Deploying Compose stack")

	registries := exec.ParseRegistryCredentials(cmd.Registry)

	if err := deployer.Deploy(cmdCtx.Context, composeFilePaths, libstack.DeployOptions{
		Options: libstack.Options{
			WorkingDir:  clonePath,
			ProjectName: cmd.ProjectName,
			Env:         cmd.Env,
			Registries:  registries,
		},
		ForceRecreate: cmd.ForceRecreateStack,
		RemoveOrphans: cmd.Prune,
	}); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to deploy Compose stack")
		return fmt.Errorf("%w: %w", exec.ErrDeployComposeFailure, err)
	}

	log.Info().Msg("Compose stack deployment complete")
	return nil
}
