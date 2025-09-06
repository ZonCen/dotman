package git

import (
	"strings"

	"github.com/ZonCen/dotman/internal"
)

func CheckIfRepo(repoPath string) (int, error) {
	return internal.Run("git", "-C", repoPath, "rev-parse", "--is-inside-work-tree")
}

func Diff(repoPath string) (int, error) {
	return internal.Run("git", "-C", repoPath, "diff", "--cached", "--quiet")
}

func Status(repoPath string) (string, error) {
	return internal.RunOutput("git", "-C", repoPath, "status", "--porcelain")
}

func ListChanges(input string) []string {
	lines := strings.Split(input, "\n")
	sliceLines := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		sliceLines = append(sliceLines, line)
	}
	return sliceLines
}

func Add(repoPath string) (int, error) {
	return internal.Run("git", "-C", repoPath, "add", "-A")
}

func Commit(repoPath, message string) (int, error) {
	return internal.Run("git", "-C", repoPath, "commit", "-m", message)
}

func Push(repoPath string) (int, error) {
	return internal.Run("git", "-C", repoPath, "push")
}

func Pull(repoPath string) (int, error) {
	return internal.Run("git", "-C", repoPath, "pull", "--ff-only")
}
