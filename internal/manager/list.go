package manager

import (
	"fmt"

	"github.com/ZonCen/dotman/internal"
	"github.com/ZonCen/dotman/internal/files"
)

func ListFiles(folderpath string) {
	internal.LogVerbose("Checking if %v exists", folderpath)
	if !internal.FileExist(folderpath) {
		fmt.Printf("could not find repofolder")
		return
	}

	internal.LogVerbose("Collecting information from in %v", folderpath+"/info.json")
	entries, err := files.ReadFile(folderpath + "/info.json")
	if err != nil {
		fmt.Printf("could not read info.json: %v\n", err)
		return
	}
	internal.LogVerbose("Presenting files in %v", folderpath)
	for filename := range entries {
		fmt.Println(filename)
	}
}
