package exec

import (
	"runtime"

	"github.com/portainer/portainer/api/filesystem"
)

func GetDockerBinaryPath() string {
	command := filesystem.JoinPaths(BIN_PATH, "docker")
	if runtime.GOOS == "windows" {
		command = filesystem.JoinPaths(BIN_PATH, "docker.exe")
	}

	return command
}

func MakeWorkingDir(target, stackName string) string {
	return filesystem.JoinPaths(target, "stacks", stackName)
}
