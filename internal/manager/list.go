package manager

import (
	"fmt"
	"os"

	"github.com/ZonCen/dotman/internal"
)

func ListFiles(folderpath string) {
	internal.LogVerbose("Checking if %v exists", folderpath)
	if !internal.FileExist(folderpath) {
		fmt.Printf("could not find repofolder")
		return
	}
	internal.LogVerbose("Found %v", folderpath)

	internal.LogVerbose("Collecting files in %v", folderpath)
	entries, err := os.ReadDir(folderpath)
	if err != nil {
		fmt.Printf("could not read repofolder: %v\n", err)
		return
	}
	internal.LogVerbose("Presenting files in %v", folderpath)
	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Println(entry.Name())
		}
	}
}
