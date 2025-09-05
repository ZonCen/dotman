package cmd

import (
	"fmt"

	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with your github repo",
	Run: func(cmd *cobra.Command, args []string) {
		// Get config
		repoPath := cfg.RepoPath
		err := manager.SyncRepo(repoPath)
		if err != nil {
			fmt.Printf("Error syncing with github: %v\n", err)
			return
		}
		fmt.Println("sync called")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
