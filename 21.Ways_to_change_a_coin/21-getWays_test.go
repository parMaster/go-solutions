// https://www.hackerrank.com/challenges/coin-change/problem
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type memo struct {
	solutions map[int32]int64
}

// getWays from HackerRank
func getWays(n int32, c []int64) int64 {

	m := memo{}
	m.solutions = make(map[int32]int64, n+1)
	for i := int32(0); i <= n; i++ {
		m.solutions[i] = m.getSubWays(n, c)
	}

	return m.solutions[n]
}

func (m *memo) getSubWays(n int32, c []int64) int64 {

	elem, ok := m.solutions[n]
	if ok {
		return elem
	}

	if n == 0 {
		return 1
	}

	if n < 0 {
		return 0
	}

	var ways int64 = 0
	for i, v := range c {
		ways += m.getSubWays(n-int32(v), c[i:])
	}

	return ways
}

func Test_getWays(t *testing.T) {

	testPairs := []struct {
		expected  int64
		targetSum int32
		numbers   []int64
	}{
		{int64(3), 3, []int64{8, 3, 1, 2}},
		{int64(4), 4, []int64{1, 2, 3}},
		{int64(5), 10, []int64{2, 5, 3, 6}},
		// {int64(96190959), 166, []int64{5, 37, 8, 39, 33, 17, 22, 32, 13, 7, 10, 35, 40, 2, 43, 49, 46, 19, 41, 1, 12, 11, 28}},
	}

	for _, p := range testPairs {
		// memo := make(map[int][]int)
		assert.Equal(t, p.expected, getWays(p.targetSum, p.numbers))
	}

	// timeout on HR, 16 seconds on Xeon
	// {int64(96190959), 166, []int64{5, 37, 8, 39, 33, 17, 22, 32, 13, 7, 10, 35, 40, 2, 43, 49, 46, 19, 41, 1, 12, 11, 28}},

}

//i func getWays(n int, c []int) int {

// }
