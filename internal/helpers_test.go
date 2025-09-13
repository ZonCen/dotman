package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZonCen/dotman/internal/testutils"
)

func TestFileExist(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "test.txt")

	// Test non-existent file
	if FileExist(testFile) {
		t.Error("Expected FileExist to return false for non-existent file")
	}

	// Test existing file
	testutils.CreateTestFile(t, testFile, "test content")
	if !FileExist(testFile) {
		t.Error("Expected FileExist to return true for existing file")
	}
}

func TestFolderExist(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	testFolder := filepath.Join(testDir, "test_folder")

	// Test non-existent folder
	if FolderExist(testFolder) {
		t.Error("Expected FolderExist to return false for non-existent folder")
	}

	// Test existing folder
	if err := os.MkdirAll(testFolder, 0755); err != nil {
		t.Fatalf("Failed to create test folder: %v", err)
	}
	if !FolderExist(testFolder) {
		t.Error("Expected FolderExist to return true for existing folder")
	}
}

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: filepath.Join(os.Getenv("HOME"), "relative/path"),
		},
		{
			name:     "tilde path",
			input:    "~/test/path",
			expected: filepath.Join(os.Getenv("HOME"), "test/path"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolvePath(tt.input)
			if err != nil {
				t.Errorf("ResolvePath() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("ResolvePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShrinkPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path in home directory",
			input:    filepath.Join(home, "test", "file.txt"),
			expected: "~/test/file.txt",
		},
		{
			name:     "path outside home directory",
			input:    "/usr/local/bin/test",
			expected: "/usr/local/bin/test",
		},
		{
			name:     "home directory itself",
			input:    home,
			expected: "~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShrinkPath(tt.input)
			if result != tt.expected {
				t.Errorf("ShrinkPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsSymlink(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	regularFile := filepath.Join(testDir, "regular.txt")
	symlinkFile := filepath.Join(testDir, "symlink.txt")
	targetFile := filepath.Join(testDir, "target.txt")

	// Create regular file
	testutils.CreateTestFile(t, regularFile, "content")
	testutils.CreateTestFile(t, targetFile, "target content")

	// Create symlink
	testutils.CreateTestSymlink(t, symlinkFile, targetFile)

	// Test regular file
	isSym, err := IsSymlink(regularFile)
	if err != nil {
		t.Errorf("IsSymlink() error = %v", err)
	}
	if isSym {
		t.Error("Expected IsSymlink to return false for regular file")
	}

	// Test symlink
	isSym, err = IsSymlink(symlinkFile)
	if err != nil {
		t.Errorf("IsSymlink() error = %v", err)
	}
	if !isSym {
		t.Error("Expected IsSymlink to return true for symlink")
	}

	// Test non-existent file
	_, err = IsSymlink(filepath.Join(testDir, "nonexistent.txt"))
	if err == nil {
		t.Error("Expected IsSymlink to return error for non-existent file")
	}
}

func TestFollowSymlink(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	symlinkFile := filepath.Join(testDir, "symlink.txt")
	targetFile := filepath.Join(testDir, "target.txt")

	testutils.CreateTestFile(t, targetFile, "target content")
	testutils.CreateTestSymlink(t, symlinkFile, targetFile)

	result, err := FollowSymlink(symlinkFile)
	if err != nil {
		t.Errorf("FollowSymlink() error = %v", err)
	}
	if result != targetFile {
		t.Errorf("FollowSymlink() = %v, want %v", result, targetFile)
	}
}

func TestNormaliseRepoURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SSH to HTTPS",
			input:    "git@github.com:user/repo.git",
			expected: "https://github.com/user/repo.git",
		},
		{
			name:     "HTTPS to SSH",
			input:    "https://github.com/user/repo.git",
			expected: "git@github.com:user/repo.git",
		},
		{
			name:     "unchanged URL",
			input:    "https://gitlab.com/user/repo.git",
			expected: "https://gitlab.com/user/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormaliseRepoURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormaliseRepoURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNormaliseRepoSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "add .git suffix",
			input:    "https://github.com/user/repo",
			expected: "https://github.com/user/repo.git",
		},
		{
			name:     "remove .git suffix",
			input:    "https://github.com/user/repo.git",
			expected: "https://github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormaliseRepoSuffix(tt.input)
			if result != tt.expected {
				t.Errorf("NormaliseRepoSuffix() = %v, want %v", result, tt.expected)
			}
		})
	}
}
