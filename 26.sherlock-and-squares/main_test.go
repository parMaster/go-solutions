package main

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Complete the 'squares' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts following parameters:
 *  1. INTEGER a
 *  2. INTEGER b
 */

func squares(a int32, b int32) int32 {

	ra := int32(math.Floor(math.Sqrt(float64(a))))
	if float64(ra) == math.Sqrt(float64(a)) {
		ra--
	}
	rb := int32(math.Floor(math.Sqrt(float64(b))))
	fmt.Println("sqrt a, ra", math.Sqrt(float64(a)), ra)
	fmt.Println("sqrt b, rb", math.Sqrt(float64(b)), rb)

	return rb - ra
}

type test struct {
	a int32
	b int32
}

func TestSquares(t *testing.T) {

	tests := []struct {
		test test
		exp  int32
	}{
		{test{3, 9}, 2},
		{test{17, 24}, 0},
		{test{24, 49}, 3},
		{test{35, 70}, 3},
		{test{100, 1000}, 22},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.exp, squares(tc.test.a, tc.test.b))
	}

}
