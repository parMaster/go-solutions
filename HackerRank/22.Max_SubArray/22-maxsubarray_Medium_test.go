package main

// https://www.hackerrank.com/challenges/maxsubarray/problem?isFullScreen=false

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func maxSubarray(arr []int32) []int32 {

	var maxElem int32 = 1 - 1<<31
	for _, current := range arr {
		maxElem = max32(current, maxElem)
	}
	// init maxSubset with a max negative numbere if negatives-only array given
	maxSubset := maxElem - max32(0, maxElem)

	var prefixSum int32 = 0
	var maxSum int32 = 1 - 1<<31
	for _, current := range arr {

		prefixSum = max32(current, prefixSum+current)
		maxSum = max32(maxSum, prefixSum)

		maxSubset += max32(current, 0)
	}

	return []int32{maxSum, maxSubset}
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func Test_MaxSubarray(t *testing.T) {

	assert.Equal(t, []int32{10, 10}, maxSubarray([]int32{1, 2, 3, 4}))
	assert.Equal(t, []int32{10, 11}, maxSubarray([]int32{2, -1, 2, 3, 4, -5}))

	assert.Equal(t, []int32{-1, -1}, maxSubarray([]int32{-2, -3, -1, -4, -6}))

}
