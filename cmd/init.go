package cmd

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/manager"
	"github.com/spf13/cobra"
)

var (
	folderPath string
	repository string
	branch     string
	force      bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise your dotman configuration and folder",
	Run: func(cmd *cobra.Command, args []string) {
		internal.LogVerbose("Trying to resolve path %v", folderPath)
		folderPath, err := internal.ResolvePath(folderPath)
		if err != nil {
			fmt.Printf("Could not read the folderPath")
			return
		}

		internal.LogVerbose("Starting the initialization")
		err = manager.Init(folderPath, repository, branch, force)
		if err != nil {
			fmt.Printf("Error initialize repository: %v\n", err)
			return
		}
		fmt.Println("init has been run successfully")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&folderPath, "folderpath", "~/dotfiles", "Folder where your dot files will be saved")
	initCmd.Flags().StringVar(&repository, "repository", "", "Which git repository you want to use. Used to run git init if needed")
	initCmd.Flags().StringVar(&branch, "branch", "main", "Which branch to use")
	initCmd.Flags().BoolVar(&force, "force", false, "If the init should change remote URL to new auth_type")
}
