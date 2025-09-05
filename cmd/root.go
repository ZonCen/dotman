package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

func initConfig() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".dotconfig")

	if !internal.FileExist(configPath) {
		fmt.Println("No config found, creating default one...")
		cfg = &config.Config{RepoPath: filepath.Join(home, "dotfiles")}
		_ = config.SaveConf(configPath, cfg)
	}

	var err error
	cfg, err = config.LoadConf(configPath)
	if err != nil {
		fmt.Println("Failed to load config")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

var rootCmd = &cobra.Command{
	Use:   "dotman",
	Short: "Dotman is a simple dotfiles manager",
	Long:  "Manage your dotfiles with ease: add, remove, list, and sync dotfiles across machines.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
