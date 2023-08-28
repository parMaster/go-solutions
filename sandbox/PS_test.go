package sandbox

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ps() error {

	proc, err := os.ReadDir("/proc")
	if err != nil {
		return errors.Join(err, fmt.Errorf("error reading /proc dir"))
	}

	ps := make(map[int]fs.DirEntry, len(proc))
	for i, p := range proc {
		pstr := filepath.Clean(p.Name())

		pnum, err := strconv.Atoi(pstr)
		if err != nil {
			continue
		}

		fmt.Printf("i, p: %d, %v\n", i, pnum)
		ps[pnum] = p
	}

	for _, p := range ps {
		fmt.Printf("p: %d\n", p)

		// read status file

		statPath := filepath.Join("/proc", p.Name(), "status")

		_, err := os.Stat(statPath)
		if err != nil {
			return errors.Join(err, fmt.Errorf("error reading %s", statPath))
		}
		statFile, err := os.ReadFile(statPath)
		if err != nil {
			return errors.Join(err, fmt.Errorf("error reading %s", statPath))
		}
		statFileStr := string(statFile)

		fmt.Printf("statFileStr: %s\n", statFileStr)

	}

	return nil
}

func Test_ps(t *testing.T) {

	assert.NoError(t, ps())
}
