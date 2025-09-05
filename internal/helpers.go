package internal

import (
	"net/http"
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

func RepoExists(url string) bool {
	cleanURL := strings.TrimSuffix(url, ".git")
	resp, err := http.Head(cleanURL)
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently
}

func Run(name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), err
	}

	if err != nil {
		return -1, err
	}

	return 0, nil
}
