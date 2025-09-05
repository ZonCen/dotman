package manager

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
)

func SyncRepo(repoPath string) error {
	// 1. Stage files
	if _, err := internal.Run("git", "-C", repoPath, "add", "-A"); err != nil {
		return fmt.Errorf("could not stage repo folder: %w", err)
	}

	// 2. Check if anything staged
	code, _ := internal.Run("git", "-C", repoPath, "diff", "--cached", "--quiet")
	if code == 1 {
		if _, err := internal.Run("git", "-C", repoPath, "commit", "-m", "dotman sync"); err != nil {
			return fmt.Errorf("could not commit changes: %w", err)
		}
	} else if code != 0 {
		return fmt.Errorf("git diff failed with exit code %d", code)
	}

	// 3. Push local commits
	if _, err := internal.Run("git", "-C", repoPath, "push"); err != nil {
		return fmt.Errorf("could not push changes: %w", err)
	}

	// 4. Optionally pull if ff-only
	if _, err := internal.Run("git", "-C", repoPath, "pull", "--ff-only"); err != nil {
		return fmt.Errorf("could not pull changes: %w", err)
	}

	return nil
}
