package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func which(args []string) (foundPaths []string) {
	arguments := os.Args
	if len(arguments) == 1 && len(args) == 0 {
		fmt.Println("Please provide an argument!")
		return
	}
	if len(args) != 0 {
		arguments = args
	}

	for _, file := range arguments {
		path := os.Getenv("PATH")
		pathSplit := filepath.SplitList(path)
		for _, directory := range pathSplit {
			fullPath := filepath.Join(directory, file)
			// Does it exist?
			fileInfo, err := os.Stat(fullPath)
			if err == nil {
				mode := fileInfo.Mode()
				// Is it a regular file?
				if mode.IsRegular() {
					// Is it executable?
					if mode&0111 != 0 {
						foundPaths = append(foundPaths, fullPath)
					}
				}
			}
		}
	}
	return slices.Compact(foundPaths)
}

func Test_which(t *testing.T) {

	assert.Equal(t, which([]string{"go"}), []string{"/usr/local/go/bin/go"})
}
