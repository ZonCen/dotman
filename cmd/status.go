package cmd

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show if the symlink file still exists",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := internal.ResolvePath("~/dotfiles/info.json")
		if err != nil {
			fmt.Println("could not resolve path:", err)
			return
		}
		err = manager.CheckStatus(filePath)
		if err != nil {
			fmt.Println("Could not run checkStatus:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
