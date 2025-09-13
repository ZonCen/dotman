package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	Verbose bool
)

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || !os.IsNotExist(err)
}

func FolderExist(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create repo directory: %w", err)
	}
	return nil
}

func CreateSymlink(symPath, filePath string) error {
	err := os.Symlink(filePath, symPath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	return nil
}
func ConfirmWithUser(msg string) bool {
	var confirm string
	for {
		fmt.Print(msg)
		_, err := fmt.Scanln(&confirm)
		if err != nil {
			// If there's an error reading input, default to "no"
			fmt.Println("Error reading input, defaulting to 'no'")
			return false
		}
		confirm = strings.ToLower(strings.TrimSpace(confirm))
		if confirm == "y" || confirm == "n" {
			break
		}
		fmt.Println("Please enter y or n")
	}

	if confirm == "y" {
		return true
	} else {
		return false
	}
}

func NormaliseRepoURL(url string) string {
	if strings.HasPrefix(url, "git@") {
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			return "https://github.com/" + parts[1]
		}
	} else if strings.HasPrefix(url, "https:") {
		parts := strings.SplitN(url, "https://github.com/", 2)
		if len(parts) == 2 {
			return "git@github.com:" + parts[1]
		}
	}

	return url
}

func NormaliseRepoSuffix(url string) string {
	if strings.HasSuffix(url, ".git") {
		url = strings.TrimSuffix(url, ".git")
	} else {
		url = url + ".git"
	}

	return url
}

func SameSuffix(s1, s2 string) bool {
	return strings.HasSuffix(s1, ".git") && strings.HasSuffix(s2, ".git")
}

func SamePrefix(s1, s2 string) bool {
	i := 0
	for i < len(s1) && i < len(s2) && s1[i] == s2[i] {
		i++
	}

	return i > 0
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

func ShrinkPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}

	return path
}

func IsSymlink(absPath string) (bool, error) {
	info, err := os.Lstat(absPath)
	if err != nil {
		return false, fmt.Errorf("file does not exist: %w", err)
	}

	return info.Mode()&os.ModeSymlink != 0, nil
}

func FollowSymlink(symPath string) (string, error) {
	folderpath, err := os.Readlink(symPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink: %w", err)
	}

	return folderpath, nil
}

func LogVerbose(msg string, args ...interface{}) {
	if Verbose {
		fmt.Printf(msg+"\n", args...)
	}
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
