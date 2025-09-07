package cmd

import (
	"fmt"

	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [filename]",
	Short: "Remove symlink and move file from repofolder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		infoPath := cfg.InfoPath

		err := manager.RemoveFile(fileName, infoPath, force)
		if err != nil {
			fmt.Printf("Error removing file: %v\n", err)
			return
		}

		fmt.Println("Successfully removed file from path")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVar(&force, "force", false, "Use to delete entry from info.json")
}
