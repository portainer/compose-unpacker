package exec

import (
	"strings"

	"github.com/portainer/portainer/api/filesystem"

	"github.com/docker/cli/cli/config/types"
	"github.com/rs/zerolog/log"
)

func MakeWorkingDir(target, stackName string) string {
	return filesystem.JoinPaths(target, "stacks", stackName)
}

func ParseRegistryCredentials(raw []string) []types.AuthConfig {
	var registries []types.AuthConfig
	for _, r := range raw {
		credentials := strings.SplitN(r, ":", 3)
		if len(credentials) != 3 {
			log.Warn().
				Str("registry", r).
				Msg("Registry is malformed, skipping login")
			continue
		}

		registries = append(registries, types.AuthConfig{
			Username:      credentials[0],
			Password:      credentials[1],
			ServerAddress: credentials[2],
		})
	}
	return registries
}
