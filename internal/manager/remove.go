package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZonCen/dotman/internal"
)

func RemoveFile(file string) error {
	internal.LogVerbose("Confirming if path %v is absolute", file)
	absPath, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("could not resolve path %s: %w", file, err)
	}
	internal.LogVerbose("Absolute path %v is correct", absPath)

	internal.LogVerbose("Checking if file %v is symlinked", absPath)
	isSym, err := internal.IsSymlink(absPath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if !isSym {
		return fmt.Errorf("file is not a symlink")
	}
	internal.LogVerbose("File is symlinked")

	internal.LogVerbose("Following symlink file %v to destination folder", absPath)
	folderpath, err := os.Readlink(absPath)
	if err != nil {
		return fmt.Errorf("failed to read symlink: %w", err)
	}
	internal.LogVerbose("Found destination at %v", folderpath)

	internal.LogVerbose("Checking if %v exists", folderpath)
	if !internal.FileExist(folderpath) {
		return fmt.Errorf("symlink target %s does not exist", folderpath)
	}
	internal.LogVerbose("File exists")

	internal.LogVerbose("Removing file %v", absPath)
	err = os.Remove(absPath)
	if err != nil {
		return fmt.Errorf("could not remove the file: %w", err)
	}
	internal.LogVerbose("File has been removed")

	internal.LogVerbose("Moving back %v to original path %v", absPath, file)
	err = os.Rename(folderpath, absPath)
	if err != nil {
		return fmt.Errorf("could not move the file: %w", err)
	}
	internal.LogVerbose("Moved")

	return nil
}
