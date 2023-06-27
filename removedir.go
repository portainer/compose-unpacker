package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

func (cmd *RemoveDirCommand) Run(cmdCtx *CommandExecutionContext) error {
	log.Info().
		Str("path", cmd.Path).
		Msg("Remove directory")

	err := os.RemoveAll(cmd.Path)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", cmd.Path).
			Msg("Failed to remove directory")
	}

	return err
}
