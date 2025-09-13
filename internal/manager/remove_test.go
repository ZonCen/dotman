package manager

import (
	"path/filepath"
	"testing"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/testutils"
)

func TestRemoveFile(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Setup: Create a file in repo and symlink to it
	testFile := filepath.Join(symlinkDir, ".zshrc")
	repoFile := filepath.Join(repoDir, ".zshrc")
	testutils.CreateTestFile(t, repoFile, "zsh configuration")
	testutils.CreateTestSymlink(t, testFile, repoFile)

	// Create info.json
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestConfig(t, infoPath, testutils.TestConfig{
		RepoPath: repoDir,
		InfoPath: infoPath,
	})

	// Add entry to info.json with actual test paths
	testutils.CreateTestFile(t, infoPath, `{
  ".zshrc": {
    "symlink": "`+testFile+`",
    "path": "`+repoFile+`",
    "status": "ok",
    "errors": null
  }
}`)

	// Test removing file
	err := RemoveFile(".zshrc", infoPath, false)
	if err != nil {
		t.Errorf("RemoveFile() error = %v", err)
	}

	// Verify file was moved back to original location (no longer a symlink)
	testutils.AssertFileExists(t, testFile)
	testutils.AssertFileContent(t, testFile, "zsh configuration")

	// Verify it's no longer a symlink
	if isSym, _ := internal.IsSymlink(testFile); isSym {
		t.Error("Expected file to no longer be a symlink after removal")
	}

	// Verify repo file no longer exists
	testutils.AssertFileNotExists(t, repoFile)
}

func TestRemoveFileForce(t *testing.T) {
	testDir, repoDir, _ := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Setup: Create a file in repo but no symlink (broken state)
	repoFile := filepath.Join(repoDir, ".zshrc")
	testutils.CreateTestFile(t, repoFile, "zsh configuration")

	// Create info.json with actual test paths
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, `{
  ".zshrc": {
    "symlink": "`+filepath.Join(testDir, "symlinks", ".zshrc")+`",
    "path": "`+repoFile+`",
    "status": "Nok",
    "errors": ["symlink does not exist"]
  }
}`)

	// Test removing file with force (should work even with broken symlink)
	err := RemoveFile(".zshrc", infoPath, true)
	if err != nil {
		t.Errorf("RemoveFile() with force error = %v", err)
	}

	// Verify repo file no longer exists
	testutils.AssertFileNotExists(t, repoFile)
}

func TestRemoveFileNonExistent(t *testing.T) {
	testDir, repoDir, _ := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create info.json with non-existent file
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, `{
  ".nonexistent": {
    "symlink": "`+filepath.Join(testDir, "symlinks", ".nonexistent")+`",
    "path": "`+filepath.Join(repoDir, ".nonexistent")+`",
    "status": "Nok",
    "errors": ["file does not exist"]
  }
}`)

	// Test removing non-existent file without force - should fail
	err := RemoveFile(".nonexistent", infoPath, false)
	if err == nil {
		t.Error("Expected error when removing non-existent file without force")
	}

	// Test removing non-existent file with force - should work
	err = RemoveFile(".nonexistent", infoPath, true)
	if err != nil {
		t.Errorf("RemoveFile() with force error = %v", err)
	}
}

func TestRemoveFromFile(t *testing.T) {
	testDir, repoDir, _ := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create info.json with multiple entries
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

	// Remove one file
	err := removeFromFile(infoPath, ".zshrc")
	if err != nil {
		t.Errorf("removeFromFile() error = %v", err)
	}

	// Verify only .vimrc remains
	content, err := testutils.ReadFileContent(infoPath)
	if err != nil {
		t.Errorf("Failed to read info.json: %v", err)
	}

	if !testutils.Contains(content, ".vimrc") {
		t.Error("Expected .vimrc to remain in info.json")
	}
	if testutils.Contains(content, ".zshrc") {
		t.Error("Expected .zshrc to be removed from info.json")
	}
}
