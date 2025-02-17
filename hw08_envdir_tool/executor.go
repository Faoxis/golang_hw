package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	if len(cmd) == 0 {
		return 1
	}
	//nolint: gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	// Copy current environment
	command.Env = os.Environ()

	// Apply our environment changes
	for key, value := range env {
		if value.NeedRemove {
			command.Env = removeFromEnv(command.Env, key)
			continue
		}
		command.Env = append(command.Env, key+"="+value.Value)
	}

	// Connect standard streams
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

func removeFromEnv(env []string, key string) []string {
	prefix := key + "="
	for i := 0; i < len(env); i++ {
		if len(env[i]) >= len(prefix) && env[i][:len(prefix)] == prefix {
			return append(env[:i], env[i+1:]...)
		}
	}
	return env
}
