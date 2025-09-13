package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZonCen/dotman/internal/testutils"
)

func TestLoadConf(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	configPath := filepath.Join(testDir, "config.yaml")

	// Test loading non-existent config
	_, err := LoadConf(configPath)
	if err == nil {
		t.Error("Expected error when loading non-existent config")
	}

	// Test loading valid config
	validConfig := `repo_path: /home/user/dotfiles
info_path: /home/user/dotfiles/info.json
`
	testutils.CreateTestFile(t, configPath, validConfig)

	cfg, err := LoadConf(configPath)
	if err != nil {
		t.Errorf("LoadConf() error = %v", err)
	}
	if cfg.FolderPath != "/home/user/dotfiles" {
		t.Errorf("LoadConf() FolderPath = %v, want /home/user/dotfiles", cfg.FolderPath)
	}
	if cfg.InfoPath != "/home/user/dotfiles/info.json" {
		t.Errorf("LoadConf() InfoPath = %v, want /home/user/dotfiles/info.json", cfg.InfoPath)
	}

	// Test loading invalid YAML
	invalidConfig := `repo_path: /home/user/dotfiles
info_path: /home/user/dotfiles/info.json
invalid_yaml: [unclosed
`
	testutils.CreateTestFile(t, configPath, invalidConfig)

	_, err = LoadConf(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML config")
	}
}

func TestSaveConf(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	configPath := filepath.Join(testDir, "config.yaml")
	cfg := &Config{
		FolderPath: "/home/user/dotfiles",
		InfoPath:   "/home/user/dotfiles/info.json",
	}

	// Test saving config
	err := SaveConf(configPath, cfg)
	if err != nil {
		t.Errorf("SaveConf() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Errorf("Failed to read saved config: %v", err)
	}

	expectedContent := "repo_path: /home/user/dotfiles\ninfo_path: /home/user/dotfiles/info.json\n"
	if string(content) != expectedContent {
		t.Errorf("Saved config content mismatch.\nExpected: %s\nGot: %s", expectedContent, string(content))
	}
}

func TestConfigRoundTrip(t *testing.T) {
	testDir := testutils.TestDir(t)
	defer testutils.CleanupTestDir(t, testDir)

	configPath := filepath.Join(testDir, "config.yaml")
	originalCfg := &Config{
		FolderPath: "/home/user/dotfiles",
		InfoPath:   "/home/user/dotfiles/info.json",
	}

	// Save config
	err := SaveConf(configPath, originalCfg)
	if err != nil {
		t.Errorf("SaveConf() error = %v", err)
	}

	// Load config
	loadedCfg, err := LoadConf(configPath)
	if err != nil {
		t.Errorf("LoadConf() error = %v", err)
	}

	// Compare
	if loadedCfg.FolderPath != originalCfg.FolderPath {
		t.Errorf("FolderPath mismatch: got %v, want %v", loadedCfg.FolderPath, originalCfg.FolderPath)
	}
	if loadedCfg.InfoPath != originalCfg.InfoPath {
		t.Errorf("InfoPath mismatch: got %v, want %v", loadedCfg.InfoPath, originalCfg.InfoPath)
	}
}
