package git

import (
	"fmt"
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

func ChangeRemote(folderPath, desiredURL string) (int, error) {
	return internal.Run("git", "-C", folderPath, "remote", "set-url", "origin", desiredURL)
}

func AddRemote(folderPath, url string) (int, error) {
	return internal.Run("git", "-C", folderPath, "remote", "add", "origin", url)
}

func FetchOrigin(folderPath string) (int, error) {
	return internal.Run("git", "-C", folderPath, "fetch", "origin")
}

func FirstCheckout(folderpath, branch string) (int, error) {
	return internal.Run("git", "-C", folderpath, "checkout", "-b", branch, "--track", "origin/"+branch)
}

func GetRemoteURL(repoPath string) (string, error) {
	out, err := internal.RunOutput("git", "-C", repoPath, "remote", "get-url", "origin")
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL %w", err)
	}

	return strings.TrimSpace(out), nil
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

func Init(repoPath string) (int, error) {
	return internal.Run("git", "-C", repoPath, "init")
}
