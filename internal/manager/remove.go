package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZonCen/dotman/internal"
)

func RemoveFile(file string) error {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("could not resolve path %s: %w", file, err)
	}

	info, err := os.Lstat(absPath)
	if err != nil {
		return fmt.Errorf("file does not exist: %w", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("file is not a symlink")
	}

	repoPath, err := os.Readlink(absPath)
	if err != nil {
		return fmt.Errorf("failed to read symlink: %w", err)
	}

	if !internal.FileExist(repoPath) {
		return fmt.Errorf("symlink target %s does not exist", repoPath)
	}

	err = os.Remove(absPath)
	if err != nil {
		return fmt.Errorf("could not remove the file: %w", err)
	}

	err = os.Rename(repoPath, absPath)
	if err != nil {
		return fmt.Errorf("could not move the file: %w", err)
	}

	return nil
}
