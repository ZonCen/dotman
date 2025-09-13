package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/config"
)

var (
	cfg *config.Config
)

func initConfig() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".dotconfig")

	if !internal.FileExist(configPath) {
		if internal.ConfirmWithUser("No config found, do you want to create one? (y/N)") {
			cfg = &config.Config{FolderPath: filepath.Join(home, "dotfiles"),
				InfoPath: filepath.Join(home, "dotfiles", "info.json")}
			_ = config.SaveConf(configPath, cfg)
		}
	}

	var err error
	cfg, err = config.LoadConf(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
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

	rootCmd.PersistentFlags().BoolVarP(&internal.Verbose, "verbose", "v", false, "Show detailed output")
}
