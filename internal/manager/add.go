package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZonCen/dotman/internal"
)

// AddFile moves a file into the repository and creates a symlink back
func AddFile(filePath string, folderPath string) error {
	fileName := filepath.Base(filePath)
	destPath := filepath.Join(folderPath, fileName)

	// Create repo folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if internal.ConfirmWithUser("Folder was not found, do you want to create one? (y/N)") {
			err := internal.CreateFolder(folderPath)
			if err != nil {
				return fmt.Errorf("could not create folder: %w", err)
			}
		}
	}

	if internal.FileExist(destPath) {
		return fmt.Errorf("file already exists")
	}

	// Move file into repo
	err := os.Rename(filePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	// Create symlink from original location to repo
	err = os.Symlink(destPath, filePath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}
