package sandbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Fibonacci recursion
func fr(n1, n2, i, n int) int {
	if i == n {
		return n2
	}
	return fr(n2, n1+n2, i+1, n)
}

func Test_Fib_recursion(t *testing.T) {
	fmt.Println(fr(1, 1, 3, 4))
	assert.Equal(t, 12586269025, fr(1, 1, 2, 50))
	assert.Equal(t, 3, fr(1, 2, 3, 4))
}
