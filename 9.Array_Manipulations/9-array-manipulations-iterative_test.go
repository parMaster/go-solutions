package main

import (
	"testing"
)

type testpair struct {
	n        int32
	queries  [][]int32
	expected int64
}

var tests = []testpair{
	{10, [][]int32{{1, 5, 3}, {4, 8, 7}, {6, 9, 1}}, 10},
	{5, [][]int32{{1, 2, 100}, {2, 5, 100}, {3, 4, 100}}, 200},
	{10, [][]int32{{2, 6, 8}, {3, 5, 7}, {1, 8, 1}, {5, 9, 15}}, 31},
}

// TestArrayManipulation is obviously a test for arrayManipulation exercice
func TestArrayManipulation(t *testing.T) {
	for _, pair := range tests {
		v := arrayManipulation(pair.n, pair.queries)
		if v != pair.expected {
			t.Error(
				"For", pair.n, pair.queries,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}
