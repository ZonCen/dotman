package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZonCen/dotman/internal"
)

// AddFile moves a file into the repository and creates a symlink back
func AddFile(filePath, folderPath string, force bool) error {
	fileName := filepath.Base(filePath)
	destPath := filepath.Join(folderPath, fileName)

	internal.LogVerbose("Checking for existing folder at %v", folderPath)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if internal.ConfirmWithUser("Folder was not found, do you want to create one? (y/N)") {
			err := internal.CreateFolder(folderPath)
			if err != nil {
				return fmt.Errorf("could not create folder: %w", err)
			}
			internal.LogVerbose("Folder created")
		}
	} else {
		internal.LogVerbose("Folder found")
	}

	internal.LogVerbose("Checking if file %v already exists", fileName)
	if internal.FileExist(destPath) && force {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("could not resolve path %s: %w", filePath, err)
		}
		isSym, err := internal.IsSymlink(absPath)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if isSym {
			return fmt.Errorf("file you trying to move (%v) is already a symlink", absPath)
		}
		internal.LogVerbose("File %v already exists, but will be overwritten", destPath)
	} else if internal.FileExist(destPath) {
		return fmt.Errorf("file already exists")
	}

	err := moveAndLink(filePath, destPath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func moveAndLink(filePath, destPath string) error {
	internal.LogVerbose("Moving %v to %v", filePath, destPath)
	err := os.Rename(filePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	internal.LogVerbose("Creating symlink at %v", filePath)
	err = os.Symlink(destPath, filePath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	return nil
}
