package exec

import (
	"errors"

	"github.com/portainer/portainer/api/filesystem"
)

const BIN_PATH = "/app"

var PORTAINER_DOCKER_CONFIG_PATH = filesystem.JoinPaths(BIN_PATH, "portainer_docker_config")
var ErrDeployComposeFailure = errors.New("stack deployment failure")
