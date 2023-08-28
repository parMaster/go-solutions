package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func which(arg string) (res string) {
	arguments := os.Args
	if len(arguments) == 1 && arg == "" {
		fmt.Println("Please provide an argument!")
		return
	}
	file := arguments[1]
	if arg != "" {
		file = arg
	}

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
					fmt.Println(fullPath)
					return fullPath
				}
			}
		}
	}
	return
}

func Test_which(t *testing.T) {

	assert.Equal(t, which("go"), "/usr/local/go/bin/go")
}
