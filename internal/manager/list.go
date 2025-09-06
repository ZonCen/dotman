package manager

import (
	"fmt"
	"os"

	"github.com/ZonCen/dotman/internal"
)

func ListFiles(folderpath string) {
	if !internal.FileExist(folderpath) {
		fmt.Printf("could not find repofolder")
		return
	}

	entries, err := os.ReadDir(folderpath)
	if err != nil {
		fmt.Printf("could not read repofolder: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Println(entry.Name())
		}
	}
}
