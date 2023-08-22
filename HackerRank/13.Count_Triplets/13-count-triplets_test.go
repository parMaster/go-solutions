package main

import (
	"testing"
)

type testpair struct {
	arr      []int64
	r        int64
	expected int64
}

var tests = []testpair{
	{[]int64{1, 2, 2, 4}, 2, 2},
	{[]int64{1, 3, 9, 9, 27, 81}, 3, 6},
	{[]int64{1, 5, 5, 25, 125}, 5, 4},
	{[]int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 1, 161700},
}

func TestCountTriplets(t *testing.T) {
	for _, pair := range tests {
		v := countTriplets(pair.arr, pair.r)
		if v != pair.expected {
			t.Error(
				"For", pair.arr, pair.r,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}
