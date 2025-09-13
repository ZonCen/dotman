package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZonCen/dotman/internal/testutils"
)

func TestSaveStatus(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	infoPath := filepath.Join(testDir, "info.json")
	info := map[string]FileInfo{
		".zshrc": {
			Symlink: "~/.zshrc",
			Path:    "~/dotfiles/.zshrc",
			Status:  "ok",
			Errors:  nil,
		},
		".vimrc": {
			Symlink: "~/.vimrc",
			Path:    "~/dotfiles/.vimrc",
			Status:  "Nok",
			Errors:  []string{"symlink does not exist"},
		},
	}

	err := SaveStatus(infoPath, info)
	if err != nil {
		t.Errorf("SaveStatus() error = %v", err)
	}

	// Verify file was created
	testutils.AssertFileExists(t, infoPath)
}

func TestReadFile(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	infoPath := filepath.Join(testDir, "info.json")

	// Test reading non-existent file
	_, err := ReadFile(infoPath)
	if err == nil {
		t.Error("Expected error when reading non-existent file")
	}

	// Test reading valid file
	validJSON := `{
  ".zshrc": {
    "symlink": "~/.zshrc",
    "path": "~/dotfiles/.zshrc",
    "status": "ok",
    "errors": null
  }
}`
	testutils.CreateTestFile(t, infoPath, validJSON)

	info, err := ReadFile(infoPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}

	if len(info) != 1 {
		t.Errorf("Expected 1 file info, got %d", len(info))
	}

	zshrcInfo, exists := info[".zshrc"]
	if !exists {
		t.Error("Expected .zshrc file info to exist")
	}

	// Note: ReadFile resolves paths, so ~ gets expanded to actual home directory
	homeDir := os.Getenv("HOME")
	expectedSymlink := filepath.Join(homeDir, ".zshrc")
	expectedPath := filepath.Join(homeDir, "dotfiles", ".zshrc")

	if zshrcInfo.Symlink != expectedSymlink {
		t.Errorf("Expected symlink %s, got %s", expectedSymlink, zshrcInfo.Symlink)
	}
	if zshrcInfo.Path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, zshrcInfo.Path)
	}
	if zshrcInfo.Status != "ok" {
		t.Errorf("Expected status ok, got %s", zshrcInfo.Status)
	}
}

func TestAddFiles(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	infoPath := filepath.Join(testDir, "info.json")

	// Create initial file
	initialInfo := map[string]FileInfo{
		".zshrc": {
			Symlink: "~/.zshrc",
			Path:    "~/dotfiles/.zshrc",
			Status:  "ok",
			Errors:  nil,
		},
	}

	err := SaveStatus(infoPath, initialInfo)
	if err != nil {
		t.Errorf("SaveStatus() error = %v", err)
	}

	// Add new file
	newInfo := map[string]FileInfo{
		".vimrc": {
			Symlink: "~/.vimrc",
			Path:    "~/dotfiles/.vimrc",
			Status:  "ok",
			Errors:  nil,
		},
	}

	err = AddFiles(infoPath, newInfo)
	if err != nil {
		t.Errorf("AddFiles() error = %v", err)
	}

	// Verify both files exist
	info, err := ReadFile(infoPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}

	if len(info) != 2 {
		t.Errorf("Expected 2 files, got %d", len(info))
	}

	if _, exists := info[".zshrc"]; !exists {
		t.Error("Expected .zshrc to still exist")
	}
	if _, exists := info[".vimrc"]; !exists {
		t.Error("Expected .vimrc to be added")
	}
}

func TestRemoveFiles(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	infoPath := filepath.Join(testDir, "info.json")

	// Create initial file with multiple entries
	initialInfo := map[string]FileInfo{
		".zshrc": {
			Symlink: "~/.zshrc",
			Path:    "~/dotfiles/.zshrc",
			Status:  "ok",
			Errors:  nil,
		},
		".vimrc": {
			Symlink: "~/.vimrc",
			Path:    "~/dotfiles/.vimrc",
			Status:  "ok",
			Errors:  nil,
		},
	}

	err := SaveStatus(infoPath, initialInfo)
	if err != nil {
		t.Errorf("SaveStatus() error = %v", err)
	}

	// Remove one file
	err = RemoveFiles(infoPath, ".zshrc")
	if err != nil {
		t.Errorf("RemoveFiles() error = %v", err)
	}

	// Verify only one file remains
	info, err := ReadFile(infoPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}

	if len(info) != 1 {
		t.Errorf("Expected 1 file, got %d", len(info))
	}

	if _, exists := info[".zshrc"]; exists {
		t.Error("Expected .zshrc to be removed")
	}
	if _, exists := info[".vimrc"]; !exists {
		t.Error("Expected .vimrc to still exist")
	}
}

func TestFileInfoRoundTrip(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	infoPath := filepath.Join(testDir, "info.json")
	originalInfo := map[string]FileInfo{
		".zshrc": {
			Symlink: "~/.zshrc",
			Path:    "~/dotfiles/.zshrc",
			Status:  "ok",
			Errors:  nil,
		},
		".vimrc": {
			Symlink: "~/.vimrc",
			Path:    "~/dotfiles/.vimrc",
			Status:  "Nok",
			Errors:  []string{"symlink does not exist", "file does not exist"},
		},
	}

	// Save
	err := SaveStatus(infoPath, originalInfo)
	if err != nil {
		t.Errorf("SaveStatus() error = %v", err)
	}

	// Load
	loadedInfo, err := ReadFile(infoPath)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}

	// Compare
	if len(loadedInfo) != len(originalInfo) {
		t.Errorf("File count mismatch: got %d, want %d", len(loadedInfo), len(originalInfo))
	}

	homeDir := os.Getenv("HOME")

	for filename, original := range originalInfo {
		loaded, exists := loadedInfo[filename]
		if !exists {
			t.Errorf("File %s not found in loaded info", filename)
			continue
		}

		// Note: ReadFile resolves paths, so we need to compare with resolved paths
		expectedSymlink := filepath.Join(homeDir, ".zshrc")
		expectedPath := filepath.Join(homeDir, "dotfiles", ".zshrc")

		if filename == ".zshrc" {
			if loaded.Symlink != expectedSymlink {
				t.Errorf("Symlink mismatch for %s: got %s, want %s", filename, loaded.Symlink, expectedSymlink)
			}
			if loaded.Path != expectedPath {
				t.Errorf("Path mismatch for %s: got %s, want %s", filename, loaded.Path, expectedPath)
			}
		} else if filename == ".vimrc" {
			expectedSymlink = filepath.Join(homeDir, ".vimrc")
			expectedPath = filepath.Join(homeDir, "dotfiles", ".vimrc")
			if loaded.Symlink != expectedSymlink {
				t.Errorf("Symlink mismatch for %s: got %s, want %s", filename, loaded.Symlink, expectedSymlink)
			}
			if loaded.Path != expectedPath {
				t.Errorf("Path mismatch for %s: got %s, want %s", filename, loaded.Path, expectedPath)
			}
		}

		if loaded.Status != original.Status {
			t.Errorf("Status mismatch for %s: got %s, want %s", filename, loaded.Status, original.Status)
		}
		if len(loaded.Errors) != len(original.Errors) {
			t.Errorf("Errors count mismatch for %s: got %d, want %d", filename, len(loaded.Errors), len(original.Errors))
		}
	}
}
