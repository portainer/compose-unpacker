package main

import (
	"context"

	"go.uber.org/zap"
)

const (
	// Path to Docker binary
	BIN_PATH = "/usr/local/bin"
)

type CommandExecutionContext struct {
	context context.Context
	logger  *zap.SugaredLogger
}
type DeployCommand struct {
	User     string `help:"Username for Git authentication." short:"u"`
	Password string `help:"Password or PAT for Git authentication" short:"p"`

	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	ProjectName              string   `arg:"" help:"Name of the Compose stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file."  name:"compose-file-paths"`
}
type UndeployCommand struct {
	User     string `help:"Username for Git authentication." short:"u"`
	Password string `help:"Password or PAT for Git authentication" short:"p"`

	GitRepository            string   `arg:"" help:"Git repository to deploy from." name:"git-repo"`
	ProjectName              string   `arg:"" help:"Name of the Compose stack." name:"project-name"`
	Destination              string   `arg:"" help:"Path on disk where the Git repository will be cloned." type:"path" name:"destination"`
	ComposeRelativeFilePaths []string `arg:"" help:"Relative path to the Compose file." name:"compose-file-path"`
}

var cli struct {
	// Generic options
	Debug bool `help:"Enable debug mode."`

	// Commands
	Deploy DeployCommand `cmd:"" help:"Deploy a stack from a Git repository."`
}

func NewCommandExecutionContext(ctx context.Context, logger *zap.SugaredLogger) *CommandExecutionContext {
	return &CommandExecutionContext{
		context: ctx,
		logger:  logger,
	}
}
