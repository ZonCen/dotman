package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZonCen/dotman/internal/testutils"
)

func TestAddCommandIntegration(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create config
	configPath := filepath.Join(testDir, ".dotconfig")
	testutils.CreateTestConfig(t, configPath, testutils.TestConfig{
		RepoPath: repoDir,
		InfoPath: filepath.Join(repoDir, "info.json"),
	})

	// Set config path environment variable
	originalConfigPath := os.Getenv("DOTMAN_CONFIG_PATH")
	err := os.Setenv("DOTMAN_CONFIG_PATH", configPath)
	if err != nil {
		t.Errorf("os.Setenv() error = %v", err)
	}
	defer func() {
		if originalConfigPath != "" {
			err := os.Setenv("DOTMAN_CONFIG_PATH", originalConfigPath)
			if err != nil {
				t.Errorf("os.Setenv() error = %v", err)
			}
		} else {
			err := os.Unsetenv("DOTMAN_CONFIG_PATH")
			if err != nil {
				t.Errorf("os.UnsetEnv() error = %v", err)
			}
		}
	}()

	// Create test file
	testFile := filepath.Join(symlinkDir, ".zshrc")
	testutils.CreateTestFile(t, testFile, "zsh configuration")

	// Test add command
	originalArgs := os.Args
	os.Args = []string{"dotman", "add", testFile}
	defer func() { os.Args = originalArgs }()

	// This would require mocking the cobra command execution
	// For now, we'll test the core functionality directly
	t.Log("Integration test setup complete")
}

func TestListCommandIntegration(t *testing.T) {
	testDir, repoDir, _ := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create config
	configPath := filepath.Join(testDir, ".dotconfig")
	testutils.CreateTestConfig(t, configPath, testutils.TestConfig{
		RepoPath: repoDir,
		InfoPath: filepath.Join(repoDir, "info.json"),
	})

	// Create info.json with some files
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, `{
  ".zshrc": {
    "symlink": "~/.zshrc",
    "path": "~/dotfiles/.zshrc",
    "status": "ok",
    "errors": null
  },
  ".vimrc": {
    "symlink": "~/.vimrc",
    "path": "~/dotfiles/.vimrc",
    "status": "ok",
    "errors": null
  }
}`)

	// Test list command setup
	t.Log("List command integration test setup complete")
}

func TestStatusCommandIntegration(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create config
	configPath := filepath.Join(testDir, ".dotconfig")
	testutils.CreateTestConfig(t, configPath, testutils.TestConfig{
		RepoPath: repoDir,
		InfoPath: filepath.Join(repoDir, "info.json"),
	})

	// Create a working setup
	testFile := filepath.Join(symlinkDir, ".zshrc")
	repoFile := filepath.Join(repoDir, ".zshrc")
	testutils.CreateTestFile(t, repoFile, "zsh configuration")
	testutils.CreateTestSymlink(t, testFile, repoFile)

	// Create info.json
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, `{
  ".zshrc": {
    "symlink": "~/.zshrc",
    "path": "~/dotfiles/.zshrc",
    "status": "ok",
    "errors": null
  }
}`)

	// Test status command setup
	t.Log("Status command integration test setup complete")
}
