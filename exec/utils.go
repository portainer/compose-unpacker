package exec

import (
	"path"
	"path/filepath"
	"runtime"
)

func GetDockerBinaryPath() string {
	command := path.Join(BIN_PATH, "docker")
	if runtime.GOOS == "windows" {
		command = path.Join(BIN_PATH, "docker.exe")
	}
	return command
}

func MakeWorkingDir(target, stackName string) string {
	return filepath.Join(target, "stacks", stackName)
}
