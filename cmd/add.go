package cmd

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a file to the dotfiles repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		filePath, _ := internal.ResolvePath(file)

		folderPath := cfg.FolderPath

		err := manager.AddFile(filePath, folderPath, force)
		if err != nil {
			fmt.Printf("Error adding file: %v\n", err)
			return
		}

		fmt.Printf("Successfully added %s to repository\n", file)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVar(&force, "force", false, "Use to override files that already exists in the folder.")
}
