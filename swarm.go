package main

import (
	"path"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

func deploySwarmStack(logger *zap.SugaredLogger, cmd SwarmDeployCommand, clonePath string) error {
	command := getDockerBinaryPath()
	args := make([]string, 0)

	if cmd.Prune {
		args = append(args, "stack", "deploy", "--prune", "--with-registry-auth")
	} else {
		args = append(args, "stack", "deploy", "--with-registry-auth")
	}
	if !cmd.Pull {
		args = append(args, "--resolve-image=never")
	}

	for _, cfile := range cmd.ComposeRelativeFilePaths {
		args = append(args, "--compose-file", path.Join(clonePath, cfile))
	}
	logger.Infow("Deploying Swarm stack", "composeFilePaths", cmd.ComposeRelativeFilePaths,
		"workingDirectory", clonePath, "projectName", cmd.ProjectName)
	args = append(args, cmd.ProjectName)

	err := runCommandAndCaptureStdErr(command, args, cmd.Env, clonePath)
	if err != nil {
		logger.Errorw("Failed to swarm deploy Git repository", "error", err)
		return errDeployComposeFailure
	}
	logger.Info("Swarm stack deployment complete")

	return err
}

func checkRunningService(logger *zap.SugaredLogger, cmd SwarmDeployCommand) ([]string, error) {
	command := getDockerBinaryPath()
	args := []string{"stack", "services", "--format={{.ID}}", cmd.ProjectName}

	logger.Infow("Checking Swarm stack", "args", args)
	output, err := runCommand(command, args)
	if err != nil {
		logger.Errorw("Failed to check running swarm services", "error", err)
		return nil, err
	}

	serviceIDs := splitLines(string(output))
	logger.Infow("Checking stack services", "service IDs", serviceIDs)
	return serviceIDs, nil
}

func updateService(logger *zap.SugaredLogger, cmd SwarmDeployCommand, serviceID string) error {
	command := getDockerBinaryPath()
	args := []string{"service", "update", serviceID}

	logger.Infow("Updating Swarm service", "args", args)
	_, err := runCommand(command, args)
	if err != nil {
		logger.Errorw("Failed to update swarm services", "error", err)
		return err
	}

	logger.Info("Update stack service completed")
	return nil
}

func splitLines(s string) []string {
	var separator string
	if runtime.GOOS == "windows" {
		separator = "\r\n"
	} else {
		separator = "\n"
	}
	parts := strings.Split(s, separator)

	ret := []string{}
	for _, part := range parts {
		// remove empty string
		if part != "" {
			ret = append(ret, part)
		}
	}
	return ret
}
