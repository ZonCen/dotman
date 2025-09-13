package manager

import (
	"fmt"
	"path/filepath"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/files"
	"github.com/ZonCen/dotman/internal/git"
)

func Init(folderPath, repository, branch string, force bool) error {
	internal.LogVerbose("Checking if %v exist", folderPath)
	if !internal.FolderExist(folderPath) {
		if internal.ConfirmWithUser("FolderPath does not exists, do you want to create one? ") {
			err := internal.CreateFolder(folderPath)
			if err != nil {
				return fmt.Errorf("error creating folder %w", err)
			}
		}
	}

	internal.LogVerbose("Checking if repository is empty or not: %v", repository)
	if repository != "" {
		internal.LogVerbose("Checking if %v is inside a working tree", folderPath)
		ok, _ := git.CheckIfRepo(folderPath)
		if ok == 128 {
			internal.LogVerbose("Running initialization of the repository at %v", folderPath)
			if internal.ConfirmWithUser("The folder has not been initialized to git, do you want to initialize it? ") {
				internal.LogVerbose("Running git init")
				_, err := git.Init(folderPath)
				if err != nil {
					return fmt.Errorf("could not initialize repository: %w", err)
				}
				fmt.Println("Folder has been initialized")
			}
			urls, err := checkRemote(folderPath, repository)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			if urls != "" && internal.ConfirmWithUser("Do you want to run git fetch, checkout and pull? ") {
				internal.LogVerbose("Running git fetch origin")
				_, err := git.FetchOrigin(folderPath)
				if err != nil {
					return fmt.Errorf("could not fetch Origin: %w", err)
				}
				internal.LogVerbose("Running git checkout -b %v --track origin/%v", branch, branch)
				_, err = git.FirstCheckout(folderPath, branch)
				if err != nil {
					return fmt.Errorf("could not checkout: %w", err)
				}
				internal.LogVerbose("Running git pull --ff-only")
				_, err = git.Pull(folderPath)
				if err != nil {
					return fmt.Errorf("could not pull from repository: %w", err)
				}
				fmt.Println("Files has been downloaded")
			}
		} else if ok == 0 {
			currentURL, err := checkRemote(folderPath, repository)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			internal.LogVerbose("Checking if %v and %v has the same prefix or are the same", currentURL, repository)
			if !internal.SamePrefix(currentURL, repository) {
				internal.LogVerbose("Changing %v to match %v", currentURL, repository)
				currentURL = internal.NormaliseRepoURL(currentURL)
				internal.LogVerbose("Checking if %v and %v has the same suffix", currentURL, repository)
				if !internal.SameSuffix(currentURL, repository) {
					internal.LogVerbose("Changing suffix of %v to match %v", currentURL, repository)
					currentURL = internal.NormaliseRepoSuffix(currentURL)
				}
				internal.LogVerbose("Checking if %v and %v is similar, and if force (%v) is active", currentURL, repository, force)
				if currentURL == repository && !force {
					fmt.Println("Folder already initialized to correct git")
				} else if currentURL == repository && force {
					if internal.ConfirmWithUser("Folder initialized to another repository. Do you want to change it? (y/N)") {
						internal.LogVerbose("Changing repository to %v", repository)
						_, err := git.ChangeRemote(folderPath, repository)
						if err != nil {
							return fmt.Errorf("error changing remote: %w", err)
						}
					}
				}
			} else if currentURL == repository {
				fmt.Println("URLs are already the same")
			} else {
				return fmt.Errorf("unknown error when checking Repository")
			}
		}
		if internal.ConfirmWithUser("Do you want to add the symlinks to the correct paths? ") {
			internal.LogVerbose("Adding symlinks to the correct paths")
			err := addSymlinks(folderPath)
			if err != nil {
				return fmt.Errorf("could not add symlinks: %w", err)
			}
			internal.LogVerbose("Symlinks has been added")
		}
	} else {
		return fmt.Errorf("repository is empty")
	}

	return nil
}

func checkRemote(folderPath, repository string) (string, error) {
	internal.LogVerbose("Checking if we can locate Remote URLs at %v", folderPath)
	urls, err := git.GetRemoteURL(folderPath)
	if err != nil {
		if internal.ConfirmWithUser("Folder has no remote, do you want to add the remote repository? ") {
			internal.LogVerbose("Running git remote add origin %v", repository)
			_, err := git.AddRemote(folderPath, repository)
			if err != nil {
				return "", fmt.Errorf("could not add origin to the repository: %w", err)
			}
			internal.LogVerbose("Checking so remotes has been added at %v", folderPath)
			urls, err = git.GetRemoteURL(folderPath)
			if err != nil {
				return "", fmt.Errorf("could not find remote urls: %w", err)
			}
		}
	}
	internal.LogVerbose("Remotes has been found")
	return urls, nil
}

func addSymlinks(folderPath string) error {
	errors := make(map[string]string)
	files, err := files.ReadFile(filepath.Join(folderPath, "info.json"))
	if err != nil {
		return fmt.Errorf("could not read files: %w", err)
	}
	for _, file := range files {
		if internal.FileExist(file.Path) {
			err := internal.CreateSymlink(file.Symlink, file.Path)
			if err != nil {
				errors[file.Symlink] = err.Error()
			}
		} else {
			return fmt.Errorf("file %v does not exist", file.Path)
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("could not create symlinks: %v", errors)
	}
	return nil
}
