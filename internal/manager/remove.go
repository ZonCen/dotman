package manager

import (
	"fmt"
	"os"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/files"
)

func RemoveFile(fileName, infoPath string, force bool) error {
	var (
		symPath  string
		filePath string
	)

	fileInfo, err := files.ReadFile(infoPath)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	symPath = fileInfo[fileName].Symlink
	filePath = fileInfo[fileName].Path

	_, err = checkSymlink(symPath)
	if err != nil && !force {
		return fmt.Errorf("file is not symlinked: %w", err)
	}

	_, err = checkPath(filePath)
	if err != nil && !force {
		return fmt.Errorf("could not process filepath: %w", err)
	}

	_, err = checkSamePath(symPath, filePath)
	if err != nil && !force {
		return fmt.Errorf("issues comparing paths: %w", err)
	}

	internal.LogVerbose("Removing file %v", symPath)
	err = os.Remove(symPath)
	if err != nil && !force {
		return fmt.Errorf("could not remove the file: %w", err)
	}

	internal.LogVerbose("Moving back %v to original path %v", filePath, symPath)
	err = os.Rename(filePath, symPath)
	if err != nil && !force {
		return fmt.Errorf("could not move the file: %w", err)
	}

	err = removeFromFile(infoPath, fileName)
	if err != nil {
		return fmt.Errorf("could not remove from file: %w", err)
	}

	return nil
}

func removeFromFile(infoPath, filename string) error {
	err := files.RemoveFiles(infoPath, filename)
	if err != nil {
		return fmt.Errorf("failed to remove from file: %w", err)
	}

	return nil
}
