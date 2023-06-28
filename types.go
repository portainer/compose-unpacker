package main

import (
	"context"
	"path"

	"github.com/portainer/compose-unpacker/log"
)

const (
	BIN_PATH            = "/app"
	UNPACKER_EXIT_ERROR = 255
)

var PORTAINER_DOCKER_CONFIG_PATH = path.Join(BIN_PATH, "portainer_docker_config")

type CommandExecutionContext struct {
	context context.Context
}

type DeployCommand struct {
	User                     string   `help:"Username for Git authentication." short:"u"`
	Password                 string   `help:"Password or PAT for Git authentication" short:"p"`
	Keep                     bool     `help:"Keep stack folder" short:"k"`
	SkipTLSVerify            bool     `help:"Skip TLS verification for git" name:"skip-tls-verify"`
	Env                      []string `help:"OS ENV for stack" example:"key=value"`
	Registry                 []string `help:"Registry credentials" name:"registry"`
	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	Reference                string   `arg:"" help:"Reference of Git repository to deploy from." name:"git-ref"`
	ProjectName              string   `arg:"" help:"Name of the Compose stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file."  name:"compose-file-paths"`
}

type SwarmDeployCommand struct {
	User                     string   `help:"Username for Git authentication." short:"u"`
	Password                 string   `help:"Password or PAT for Git authentication" short:"p"`
	Pull                     bool     `help:"Pull Image" short:"f"`
	Prune                    bool     `help:"Prune services during deployment" short:"r"`
	Keep                     bool     `help:"Keep stack folder" short:"k"`
	SkipTLSVerify            bool     `help:"Skip TLS verification for git" name:"skip-tls-verify"`
	Env                      []string `help:"OS ENV for stack."`
	Registry                 []string `help:"Registry credentials" name:"registry"`
	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	Reference                string   `arg:"" help:"Reference of Git repository to deploy from." name:"git-ref"`
	ProjectName              string   `arg:"" help:"Name of the Swarm stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file."  name:"compose-file-paths"`
}

type UndeployCommand struct {
	User     string `help:"Username for Git authentication." short:"u"`
	Password string `help:"Password or PAT for Git authentication" short:"p"`
	Keep     bool   `help:"Keep stack folder" short:"k"`

	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	ProjectName              string   `arg:"" help:"Name of the Compose stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file." name:"compose-file-path"`
}

type SwarmUndeployCommand struct {
	Keep        bool   `help:"Keep stack folder" short:"k"`
	ProjectName string `arg:"" help:"Name of the Compose (Swarm) stack." name:"project-name"`
	Destination string `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
}

type RemoveDirCommand struct {
	Path string `arg:"" help:"The path be removed." name:"path"`
}

var cli struct {
	LogLevel      log.Level            `kong:"help='Set the logging level',default='INFO',enum='DEBUG,INFO,WARN,ERROR',env='LOG_LEVEL'"`
	PrettyLog     bool                 `kong:"help='Whether to enable or disable colored logs output',default='false',env='PRETTY_LOG'"`
	Deploy        DeployCommand        `cmd:"" help:"Deploy a stack from a Git repository."`
	Undeploy      UndeployCommand      `cmd:"" help:"Remove a stack from a Git repository."`
	SwarmDeploy   SwarmDeployCommand   `cmd:"" help:"Deploy a Swarm stack from a Git repository."`
	SwarmUndeploy SwarmUndeployCommand `cmd:"" help:"Remove a Swarm stack from a Git repository."`
	RemoveDir     RemoveDirCommand     `cmd:"" help:"Remove a directory."`
}

func NewCommandExecutionContext(ctx context.Context) *CommandExecutionContext {
	return &CommandExecutionContext{
		context: ctx,
	}
}
