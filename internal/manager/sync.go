package manager

import (
	"fmt"
	"strings"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/git"
)

func SyncRepo(folderPath string, dryrun, download, upload bool) error {
	fmt.Println("Checking for valid repository")
	code, err := git.CheckIfRepo(folderPath)
	if err != nil || code != 0 {
		return fmt.Errorf("not a git repository: %s", folderPath)
	}

	fmt.Println("Repository detected at:", folderPath)
	if dryrun {
		fmt.Println("[dry-run] Collecting local changes")
	} else {
		fmt.Println("Collecting local changes")
	}

	output, err := git.Status(folderPath)
	if err != nil {
		return fmt.Errorf("failed to collect status: %w", err)
	}
	if strings.TrimSpace(output) == "" {
		if dryrun {
			return fmt.Errorf("[dry-run] no changes detected")
		} else {
			return fmt.Errorf("no changes detected")
		}
	}

	if dryrun {
		fmt.Println("[dry-run] Changes detected, following files would be staged and commited with commit message 'dotman sync':")
		printChanges(output)

		return nil
	}

	if upload {
		fmt.Println("Following files will be commited and pushed:")
		printChanges(output)

		if _, err := git.Add(folderPath); err != nil {
			return fmt.Errorf("could not stage repo folder: %w", err)
		}

		code, _ = git.Diff(folderPath)
		if code == 1 {
			if _, err := git.Commit(folderPath, "dotman sync"); err != nil {
				return fmt.Errorf("could not commit changes: %w", err)
			}
		} else if code != 0 {
			return fmt.Errorf("git diff failed with exit code %d", code)
		}

		if _, err := git.Push(folderPath); err != nil {
			return fmt.Errorf("could not push changes: %w", err)
		}
	}

	if download {
		output, err := git.Status(folderPath)
		if err != nil {
			return fmt.Errorf("failed to collect status: %w", err)
		}

		if strings.TrimSpace(output) != "" {
			if !internal.ConfirmWithUser("[warning] Local changes detected, pull may fail or cause conflicts. Continue? (y/N)") {
				return fmt.Errorf("aborting downloading changes from git")
			}
		}

		if _, err := git.Pull(folderPath); err != nil {
			return fmt.Errorf("could not pull changes: %w", err)
		}
	}

	return nil
}

func printChanges(output string) {
	for _, change := range git.ListChanges(output) {
		fmt.Println(change)
	}
}
