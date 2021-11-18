package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// brute force naive solution
func countInversions(arr []int32) int64 {
	var res int64
	for i, e := range arr {
		for j := i; j < len(arr); j++ {
			if e > arr[j] {
				res++
			}
		}
	}
	return res
}

// Placeholder for faster solution
func countInversions2(arr []int32) int64 {
	var res int64

	return res
}

func Test_countInversions(t *testing.T) {
	assert.Equal(t, int64(0), countInversions([]int32{1, 1, 1, 2, 2}))
	assert.Equal(t, int64(4), countInversions([]int32{2, 1, 3, 1, 2}))
	assert.Equal(t, int64(1), countInversions([]int32{1, 5, 3, 7}))
	assert.Equal(t, int64(6), countInversions([]int32{7, 5, 3, 1}))
	assert.Equal(t, int64(0), countInversions([]int32{1, 3, 5, 7}))
	assert.Equal(t, int64(3), countInversions([]int32{3, 2, 1}))

	assert.Equal(t, int64(0), countInversions2([]int32{1, 1, 1, 2, 2}))
	assert.Equal(t, int64(4), countInversions2([]int32{2, 1, 3, 1, 2}))
	assert.Equal(t, int64(1), countInversions2([]int32{1, 5, 3, 7}))
	assert.Equal(t, int64(6), countInversions2([]int32{7, 5, 3, 1}))
	assert.Equal(t, int64(0), countInversions2([]int32{1, 3, 5, 7}))
	assert.Equal(t, int64(3), countInversions2([]int32{3, 2, 1}))
}
