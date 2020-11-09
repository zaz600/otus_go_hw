package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if err := setEnv(env); err != nil {
		return -1
	}

	cmdName := cmd[0]
	subProc := exec.Command(cmdName, cmd[1:]...)
	subProc.Stdout = os.Stdout
	subProc.Stderr = os.Stderr
	subProc.Stdin = os.Stdin

	if err := subProc.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return -1
	}

	return 0
}

func setEnv(env Environment) error {
	for envVar, val := range env {
		if val == "" {
			os.Unsetenv(envVar)
			continue
		}
		if err := os.Setenv(envVar, val); err != nil {
			return err
		}
	}
	return nil
}
