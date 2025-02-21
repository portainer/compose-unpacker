package exec

import (
	"context"
)

type CommandExecutionContext struct {
	Context context.Context
}

func NewCommandExecutionContext(ctx context.Context) *CommandExecutionContext {
	return &CommandExecutionContext{
		Context: ctx,
	}
}
