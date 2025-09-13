package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDir creates a temporary directory for testing and returns its path
func TestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "dotman_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// CleanupTestDir removes a test directory
func CleanupTestDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup test dir %s: %v", dir, err)
	}
}

// CreateTestFile creates a test file with given content
func CreateTestFile(t *testing.T, path, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create dir %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file %s: %v", path, err)
	}
}

// CreateTestSymlink creates a test symlink
func CreateTestSymlink(t *testing.T, symlinkPath, targetPath string) {
	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		t.Fatalf("Failed to create symlink %s -> %s: %v", symlinkPath, targetPath, err)
	}
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", path)
	}
}

// AssertFileNotExists checks if a file doesn't exist
func AssertFileNotExists(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		t.Errorf("Expected file %s to not exist, but it does", path)
	}
}

// AssertSymlink checks if a path is a symlink pointing to target
func AssertSymlink(t *testing.T, symlinkPath, expectedTarget string) {
	info, err := os.Lstat(symlinkPath)
	if err != nil {
		t.Errorf("Failed to stat symlink %s: %v", symlinkPath, err)
		return
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("Expected %s to be a symlink, but it's not", symlinkPath)
		return
	}

	target, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Errorf("Failed to read symlink %s: %v", symlinkPath, err)
		return
	}

	if target != expectedTarget {
		t.Errorf("Symlink %s points to %s, expected %s", symlinkPath, target, expectedTarget)
	}
}

// AssertFileContent checks if a file contains expected content
func AssertFileContent(t *testing.T, path, expectedContent string) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read file %s: %v", path, err)
		return
	}

	if string(content) != expectedContent {
		t.Errorf("File %s content mismatch.\nExpected: %s\nGot: %s", path, expectedContent, string(content))
	}
}

// SetupTestEnvironment creates a complete test environment
func SetupTestEnvironment(t *testing.T) (string, string, string) {
	testDir := TestDir(t)
	repoDir := filepath.Join(testDir, "repo")
	symlinkDir := filepath.Join(testDir, "symlinks")

	// Create directories
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}
	if err := os.MkdirAll(symlinkDir, 0755); err != nil {
		t.Fatalf("Failed to create symlink dir: %v", err)
	}

	return testDir, repoDir, symlinkDir
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment(t *testing.T, testDir string) {
	CleanupTestDir(t, testDir)
}

// MockUserInput simulates user input for testing
type MockUserInput struct {
	responses []string
	index     int
}

func NewMockUserInput(responses ...string) *MockUserInput {
	return &MockUserInput{
		responses: responses,
		index:     0,
	}
}

func (m *MockUserInput) GetResponse() string {
	if m.index >= len(m.responses) {
		return "n" // Default to "no" if we run out of responses
	}
	response := m.responses[m.index]
	m.index++
	return response
}

// TestConfig represents a test configuration
type TestConfig struct {
	RepoPath string
	InfoPath string
}

// CreateTestConfig creates a test configuration file
func CreateTestConfig(t *testing.T, configPath string, cfg TestConfig) {
	configContent := fmt.Sprintf(`repo_path: %s
info_path: %s
`, cfg.RepoPath, cfg.InfoPath)

	CreateTestFile(t, configPath, configContent)
}

// ReadFileContent reads file content as string
func ReadFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Contains checks if a string contains a substring
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
