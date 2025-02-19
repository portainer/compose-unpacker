package exec

import (
	"errors"
	"path"
)

const BIN_PATH = "/app"

var PORTAINER_DOCKER_CONFIG_PATH = path.Join(BIN_PATH, "portainer_docker_config")
var ErrDeployComposeFailure = errors.New("stack deployment failure")
