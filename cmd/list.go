package cmd

import (
	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list handled files",
	Short: "Will list all files in the configured repopath",
	Run: func(cmd *cobra.Command, args []string) {
		folderPath := cfg.FolderPath

		manager.ListFiles(folderPath)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
