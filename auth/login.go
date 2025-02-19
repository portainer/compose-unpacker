package auth

import (
	"strings"

	"github.com/portainer/compose-unpacker/exec"
	"github.com/rs/zerolog/log"
)

func DockerLogin(registries []string) error {
	command := exec.GetDockerBinaryPath()

	for _, registry := range registries {
		credentials := strings.Split(registry, ":")
		if len(credentials) != 3 {
			log.Warn().
				Str("registry", registry).
				Msg("registry is malformed. Skip login it.")

			continue
		}

		args := make([]string, 0)
		args = append(args, "--config", exec.PORTAINER_DOCKER_CONFIG_PATH, "login", "--username", credentials[0], "--password", credentials[1], credentials[2])

		if err := exec.RunCommandAndCaptureStdErr(command, args, nil, ""); err != nil {
			log.Warn().
				Err(err).
				Msgf("Docker login %s failed, skipping login", credentials[2])

			continue
		}

		log.Info().Msgf("Docker login %s succedeed", credentials[2])
	}

	return nil
}
