package internal

import (
	"os"
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
