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
		credentials := strings.Split(r, ":")
		partsLen := len(credentials)
		if partsLen != 3 && partsLen != 4 {
			log.Warn().
				Str("registry", r).
				Msg("Registry is malformed, skipping login")
			continue
		}

		serverAddr := credentials[2]
		if partsLen == 4 {
			serverAddr += ":" + credentials[3]
		}

		registries = append(registries, types.AuthConfig{
			Username:      credentials[0],
			Password:      credentials[1],
			ServerAddress: serverAddr,
		})
	}
	return registries
}
