package manager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZonCen/dotman/internal/testutils"
)

func TestAddFile(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create a test file to add
	testFile := filepath.Join(symlinkDir, ".zshrc")
	testutils.CreateTestFile(t, testFile, "zsh configuration")

	// Create empty info.json first
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, "{}")

	// Test adding file
	err := AddFile(testFile, repoDir, false)
	if err != nil {
		t.Errorf("AddFile() error = %v", err)
	}

	// Verify file was moved to repo
	repoFile := filepath.Join(repoDir, ".zshrc")
	testutils.AssertFileExists(t, repoFile)
	testutils.AssertFileContent(t, repoFile, "zsh configuration")

	// Verify symlink was created
	testutils.AssertSymlink(t, testFile, repoFile)

	// Verify info.json was updated
	testutils.AssertFileExists(t, infoPath)
}

func TestAddFileForce(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create empty info.json first
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, "{}")

	// Create a test file
	testFile := filepath.Join(symlinkDir, ".zshrc")
	testutils.CreateTestFile(t, testFile, "zsh configuration")

	// Add file first time
	err := AddFile(testFile, repoDir, false)
	if err != nil {
		t.Errorf("AddFile() error = %v", err)
	}

	// Remove the symlink to simulate a broken state
	err = os.Remove(testFile)
	if err != nil {
		t.Errorf("os.Remove() error = %v", err)
	}
	// Create another file with same name
	testFile2 := filepath.Join(symlinkDir, ".zshrc")
	testutils.CreateTestFile(t, testFile2, "new zsh configuration")

	// Try to add again with force
	err = AddFile(testFile2, repoDir, true)
	if err != nil {
		t.Errorf("AddFile() with force error = %v", err)
	}

	// Verify the file was overwritten
	repoFile := filepath.Join(repoDir, ".zshrc")
	testutils.AssertFileContent(t, repoFile, "new zsh configuration")
}

func TestAddFileAlreadyExists(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create empty info.json first
	infoPath := filepath.Join(repoDir, "info.json")
	testutils.CreateTestFile(t, infoPath, "{}")

	// Create a test file
	testFile := filepath.Join(symlinkDir, ".zshrc")
	testutils.CreateTestFile(t, testFile, "zsh configuration")

	// Add file first time
	err := AddFile(testFile, repoDir, false)
	if err != nil {
		t.Errorf("AddFile() error = %v", err)
	}

	// Create another file with same name
	testFile2 := filepath.Join(symlinkDir, ".zshrc")
	testutils.CreateTestFile(t, testFile2, "new zsh configuration")

	// Try to add again without force - should fail
	err = AddFile(testFile2, repoDir, false)
	if err == nil {
		t.Error("Expected error when adding file that already exists without force")
	}
}

func TestAddFileNonExistentSource(t *testing.T) {
	testDir, repoDir, _ := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Try to add non-existent file
	nonExistentFile := filepath.Join(testDir, "nonexistent.txt")

	err := AddFile(nonExistentFile, repoDir, false)
	if err == nil {
		t.Error("Expected error when adding non-existent file")
	}
}

func TestMoveAndLink(t *testing.T) {
	testDir, repoDir, symlinkDir := testutils.SetupTestEnvironment(t)
	defer testutils.CleanupTestEnvironment(t, testDir)

	// Create a test file
	testFile := filepath.Join(symlinkDir, ".vimrc")
	testutils.CreateTestFile(t, testFile, "vim configuration")

	// Move and link
	destPath := filepath.Join(repoDir, ".vimrc")
	err := moveAndLink(testFile, destPath)
	if err != nil {
		t.Errorf("moveAndLink() error = %v", err)
	}

	// Verify file was moved
	testutils.AssertFileExists(t, destPath)
	testutils.AssertFileContent(t, destPath, "vim configuration")

	// Verify symlink was created (original path should now be a symlink)
	testutils.AssertSymlink(t, testFile, destPath)
}
