package exec

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func RunCommandAndCaptureStdErr(command string, args []string, env []string, workingDir string) error {
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Dir = workingDir

	if env != nil {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, env...)
	}

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func RunCommand(command string, args []string) (string, error) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	cmd := exec.Command(command, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
}
