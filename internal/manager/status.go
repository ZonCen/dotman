package manager

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/files"
)

/*
First we need to check if a file exist or not for our status file
If file exist we should populate a FileInfo with the data.
For each file we should check so the symlink file exists (example if ~/.dotconfig)
For each file we should check so the Path file exists (example if ~/dotfiles/.dotconfig)
For each file we should update Status to either "Ok" or "Nok" based on above.
If NOK, provide what the issue is.
*/

func CheckStatus(filePath string) error {
	internal.LogVerbose("Checking if %v exists", filePath)
	if !internal.FileExist(filePath) {
		return fmt.Errorf("could not find the file %v", filePath)
	}

	internal.LogVerbose("Collecting data from %v", filePath)
	fileInfo, err := files.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	errorFiles := make(map[string]files.FileInfo)

	internal.LogVerbose("Checking entries")
	for filename, info := range fileInfo {
		internal.LogVerbose("Resetting errors before continue")
		if len(info.Errors) > 0 {
			info.Errors = nil
			fileInfo[filename] = info
		}
		internal.LogVerbose("Checking %s with current information: symlink=%s, path=%s, status=%s",
			filename, info.Symlink, info.Path, info.Status)
		symOK, err := checkSymlink(info.Symlink)
		if err != nil {
			info.Status = "Nok"
			info.Errors = append(info.Errors, err.Error())
			fileInfo[filename] = info
			errorFiles[filename] = info
		}
		fileOK, err := checkPath(info.Path)
		if err != nil {
			info.Status = "Nok"
			info.Errors = append(info.Errors, err.Error())
			fileInfo[filename] = info
			errorFiles[filename] = info
		}
		if symOK && fileOK {
			_, err := checkSamePath(info.Symlink, info.Path)
			if err != nil {
				info.Status = "Nok"
				info.Errors = append(info.Errors, err.Error())
				fileInfo[filename] = info
				errorFiles[filename] = info
			}
		}
	}

	if len(errorFiles) > 0 {
		if internal.Verbose {
			internal.LogVerbose("Presenting files that is in a Nok state")
		} else {
			fmt.Println("Following files are not in a good state")
		}
		for filename, info := range errorFiles {
			fmt.Printf("File: %s -> symlink=%s, path=%s, status=%s\n",
				filename, info.Symlink, info.Path, info.Status)

			for _, err := range info.Errors {
				fmt.Printf("  Error: %s\n", err)
			}
		}
	}

	err = files.SaveStatus(filePath, fileInfo)
	if err != nil {
		return fmt.Errorf("could not save file %v due to error: %w", filePath, err)
	}

	return nil
}

func checkSymlink(symlink string) (bool, error) {
	path, err := internal.ResolvePath(symlink)
	if err != nil {
		return false, fmt.Errorf("could not resolve path %w", err)
	}
	internal.LogVerbose("Checking if %v exists and is a symlink", path)

	exists := internal.FileExist(path)
	if !exists {
		return false, fmt.Errorf("symlink file %v does not exist", path)
	}

	isSYm, err := internal.IsSymlink(path)
	if err != nil {
		return false, fmt.Errorf("file %v is not a symlink: %w", path, err)
	}

	return isSYm, nil
}

func checkPath(path string) (bool, error) {
	absPath, err := internal.ResolvePath(path)
	if err != nil {
		return false, fmt.Errorf("could not resolve path (%v): %w", path, err)
	}
	internal.LogVerbose("Checking if %v exists", absPath)

	if !internal.FileExist(absPath) {
		return false, fmt.Errorf("file %v does not exist", absPath)
	}

	return true, nil
}

func checkSamePath(symlink, path string) (bool, error) {
	internal.LogVerbose("Checking if symlink is pointing to correct filepath")
	symPath, err := internal.ResolvePath(symlink)
	if err != nil {
		return false, fmt.Errorf("could not resolve symlink path (%v): %w", symlink, err)
	}

	filePath, err := internal.ResolvePath(path)
	if err != nil {
		return false, fmt.Errorf("could not resolve filepath (%v): %w", path, err)
	}

	sp, err := internal.FollowSymlink(symPath)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	if sp != filePath {
		return false, fmt.Errorf("following symlink (%v) does not match filepath (%v)", sp, filePath)
	}

	return true, nil
}
