package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// brute force naive solution
func countInversionsNaive(arr []int32) int64 {
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

func midway(from, to int) int {
	return int((to-from)/2) + from
}

// classic merge sort, counting number of inversions(insertions)
func mergeCount(arr []int32, from, mid, to int) int64 {
	// merging left arr[from:mid] and right arr[mid+1:to]
	var inversionsCount int64

	// i - left subarray iterator
	// j - right subarray iterator

	tmp := make([]int32, to-from+1)
	var tmpi int

	i := from
	j := mid + 1

	// combine subarrays into new temp array
	for i <= mid && j <= to {

		if arr[i] <= arr[j] {
			tmp[tmpi] = arr[i]
			i++
		} else {
			//swap
			tmp[tmpi] = arr[j]
			j++
			inversionsCount += int64((mid + 1) - i)
		}
		tmpi++
	}

	// put everything remained on the left side
	for i <= mid {
		tmp[tmpi] = arr[i]
		i++
		tmpi++
	}

	// put everything remained on the right size
	for j <= to {
		tmp[tmpi] = arr[j]
		j++
		tmpi++
	}

	// put everything back into original array from the temp array
	for tmpi := 0; tmpi <= to-from; tmpi++ {
		arr[tmpi+from] = tmp[tmpi]
	}

	return inversionsCount
}

// subCount - recursively divides the array and counts inversions in subarrays
func subCount(arr []int32, from, to int) int64 {
	var inversionsCount int64

	if from >= to {
		return 0
	}

	inversionsCount = subCount(arr, from, midway(from, to))
	inversionsCount += subCount(arr, midway(from, to)+1, to)
	inversionsCount += mergeCount(arr, from, midway(from, to), to)

	return inversionsCount
}

func countInversions(arr []int32) int64 {
	return subCount(arr, 0, len(arr)-1)
}

func Test_countInversions(t *testing.T) {
	assert.Equal(t, int64(0), countInversionsNaive([]int32{1, 1, 1, 2, 2}))
	assert.Equal(t, int64(4), countInversionsNaive([]int32{2, 1, 3, 1, 2}))
	assert.Equal(t, int64(1), countInversionsNaive([]int32{1, 5, 3, 7}))
	assert.Equal(t, int64(6), countInversionsNaive([]int32{7, 5, 3, 1}))
	assert.Equal(t, int64(0), countInversionsNaive([]int32{1, 3, 5, 7}))
	assert.Equal(t, int64(3), countInversionsNaive([]int32{3, 2, 1}))

	// Learn to properly determine the middle of given array indices
	assert.Equal(t, 2, midway(0, 4)) // id=2 is the middle of 5 elements array
	assert.Equal(t, 0, midway(0, 1)) // 2 elements array
	assert.Equal(t, 2, midway(0, 5)) // 6 elements

	// learn to merge arrays and count inversions
	assert.Equal(t, int64(1), mergeCount([]int32{1, 2, 1, 2, 3}, 0, 1, 4))
	assert.Equal(t, int64(2), mergeCount([]int32{1, 3, 1, 2, 3}, 0, 1, 4))

	// subCount for small arrays (edge cases)
	assert.Equal(t, int64(0), countInversions([]int32{1, 2}))
	assert.Equal(t, int64(1), countInversions([]int32{2, 1}))

	// All together - divide, merge and count
	assert.Equal(t, int64(0), countInversions([]int32{1, 1, 1, 2, 2}))
	assert.Equal(t, int64(4), countInversions([]int32{2, 1, 3, 1, 2}))
	assert.Equal(t, int64(1), countInversions([]int32{1, 5, 3, 7}))
	assert.Equal(t, int64(6), countInversions([]int32{7, 5, 3, 1}))
	assert.Equal(t, int64(0), countInversions([]int32{1, 3, 5, 7}))
	assert.Equal(t, int64(3), countInversions([]int32{3, 2, 1}))
}
