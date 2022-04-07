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
