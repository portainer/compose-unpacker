package commands

import (
	"os"

	"github.com/portainer/compose-unpacker/exec"
	"github.com/rs/zerolog/log"
)

type RemoveDirCommand struct {
	Path string `arg:"" help:"The path be removed." name:"path"`
}

func (cmd *RemoveDirCommand) Run(cmdCtx *exec.CommandExecutionContext) error {
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
