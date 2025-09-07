package files

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"

	"github.com/ZonCen/dotman/internal"
)

type FileInfo struct {
	Symlink string   `json:"symlink"`
	Path    string   `json:"path"`
	Status  string   `json:"status"`
	Errors  []string `json:"errors"`
}

func SaveStatus(path string, info map[string]FileInfo) error {
	internal.LogVerbose("Marshal information and adding indentations")
	jsonBytes, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal information: %w", err)
	}

	internal.LogVerbose("Writing data to %v", path)
	if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write to disk: %w", err)
	}

	return nil
}

func ReadFile(path string) (map[string]FileInfo, error) {
	internal.LogVerbose("Reading bytes from %v", path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read the file: %w", err)
	}
	internal.LogVerbose("Unmarshal bytes")
	data := make(map[string]FileInfo)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal data: %w", err)
	}

	for fileName, info := range data {
		filePath, err := internal.ResolvePath(info.Path)
		if err != nil {
			return nil, fmt.Errorf("could not resolve path for filepath: %w", err)
		}
		symlinkPath, err := internal.ResolvePath(info.Symlink)
		if err != nil {
			return nil, fmt.Errorf("could not resolve path for symlink: %w", err)
		}
		info.Path = filePath
		info.Symlink = symlinkPath

		data[fileName] = info
	}

	return data, nil
}

func AddFiles(path string, info map[string]FileInfo) error {
	internal.LogVerbose("Collecting data from disk")
	localFileInfo, err := ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	internal.LogVerbose("Saving new data")
	maps.Copy(localFileInfo, info)

	internal.LogVerbose("Saving data to disk")
	err = SaveStatus(path, localFileInfo)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func RemoveFiles(path, filename string) error {
	internal.LogVerbose("Collecting data from disk")
	localFileInfo, err := ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	internal.LogVerbose("Removing data")
	delete(localFileInfo, filename)

	internal.LogVerbose("Saving data to disk")
	err = SaveStatus(path, localFileInfo)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
