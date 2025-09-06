package manager

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/git"
)

func Init(folderPath, repository, branch string, force bool) error {

	if !internal.FolderExist(folderPath) {
		if internal.ConfirmWithUser("FolderPath does not exists, do you want to create one? ") {
			err := internal.CreateFolder(folderPath)
			if err != nil {
				return fmt.Errorf("error creating folder %w", err)
			}
		}
	}

	if repository != "" {
		ok, _ := git.CheckIfRepo(folderPath)
		if ok == 128 {
			if internal.ConfirmWithUser("The folder has not been initialized to git, do you want to initialize it? ") {
				_, err := git.Init(folderPath)
				if err != nil {
					return fmt.Errorf("could not initialize repository: %w", err)
				}
				fmt.Println("Folder has been initialized")
				if internal.ConfirmWithUser("Folder has no remote, do you want to add the remote repository? ") {
					_, err := git.AddRemote(folderPath, repository)
					if err != nil {
						return fmt.Errorf("could not add origin to the repository: %w", err)
					}
				}
				if internal.ConfirmWithUser("Do you want to run git fetch, checkout and pull? ") {
					_, err := git.FetchOrigin(folderPath)
					if err != nil {
						return fmt.Errorf("could not fetch Origin: %w", err)

					}
					_, err = git.FirstCheckout(folderPath, branch)
					if err != nil {
						return fmt.Errorf("could not checkout: %w", err)

					}
					_, err = git.Pull(folderPath)
					if err != nil {
						return fmt.Errorf("could not pull from repository: %w", err)
					}
					fmt.Println("Files has been downloaded")
				}
			}
		} else if ok == 0 {
			currentURL, err := git.GetRemoteURL(folderPath)
			if err != nil {
				return fmt.Errorf("error checking repository: %w", err)

			}
			if !internal.SamePrefix(currentURL, repository) {
				currentURL = internal.NormaliseRepoURL(currentURL)
				if !internal.SameSuffix(currentURL, repository) {
					currentURL = internal.NormaliseRepoSuffix(currentURL)
				}
				if currentURL == repository && !force {
					fmt.Println("Folder already initialized to correct git")
				} else if currentURL == repository && force {
					if internal.ConfirmWithUser("Folder initialized to another repository. Do you want to change it? (y/N)") {
						_, err := git.ChangeRemote(folderPath, repository)
						if err != nil {
							return fmt.Errorf("error changing remote: %w", err)

						}
					}
				}
			}
		}
	}

	return nil
}
