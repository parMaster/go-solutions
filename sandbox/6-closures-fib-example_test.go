package sandbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	var f0, f1 int = 0, 1

	return func() int {
		res := f0
		t := f1
		f1 = f0 + f1
		f0 = t
		return res
	}
}

func TestFibExample(t *testing.T) {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
	assert.Equal(t, f(), 55)

}
