package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func ResolvePath(input string) (string, error) {
	if strings.HasPrefix(input, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, input[1:]), nil
	}
	if filepath.IsAbs(input) {
		return input, nil
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, input), nil
}

func Run(name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), err
	}

	if err != nil {
		return -1, err
	}

	return 0, nil
}

func RunOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return string(out), exitErr
	}
	return string(out), err
}
