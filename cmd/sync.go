package cmd

import (
	"fmt"

	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

var (
	dryRun   bool
	download bool
	upload   bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with your github repo",
	Run: func(cmd *cobra.Command, args []string) {
		folderPath := cfg.FolderPath
		// Dry-run
		if dryRun {
			fmt.Println("Will run in dry-run mode")
			download = false
			upload = false
		}
		// Download
		if download && !upload {
			fmt.Println("Will only download files")
		}
		// Upload
		if upload && !download {
			fmt.Println("Will only upload files")
		}

		err := manager.SyncRepo(folderPath, dryRun, download, upload)
		if err != nil {
			fmt.Printf("Error syncing with github: %v\n", err)
			return
		}
		fmt.Println("Sync completed")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Run sync without making changes")
	syncCmd.Flags().BoolVar(&download, "download", true, "Download remote changes only")
	syncCmd.Flags().BoolVar(&upload, "upload", true, "Uploads local changes only")
}
