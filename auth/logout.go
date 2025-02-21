package auth

import (
	"strings"

	"github.com/portainer/compose-unpacker/exec"
	"github.com/rs/zerolog/log"
)

func DockerLogout(registries []string) error {
	command := exec.GetDockerBinaryPath()

	for _, registry := range registries {
		credentials := strings.Split(registry, ":")
		if len(credentials) != 3 {
			log.Warn().
				Str("registry", registry).
				Msg("Registry is malformed, skipping logout")

			continue
		}

		args := make([]string, 0)
		args = append(args, "--config", exec.PORTAINER_DOCKER_CONFIG_PATH, "logout", credentials[2])

		if err := exec.RunCommandAndCaptureStdErr(command, args, nil, ""); err != nil {
			log.Warn().
				Err(err).
				Msgf("Docker logout %s failed, skipping logout", credentials[2])

			continue
		}

		log.Info().Msgf("Docker logout %s succedeed", credentials[2])
	}

	return nil
}
