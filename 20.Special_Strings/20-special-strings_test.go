package main

import (
	"testing"
)

type cases struct {
	n        int32
	s        string
	expected int64
}

var tests = []cases{
	{5, "asasd", 7},
	{7, "abcbaba", 10},
	{4, "aaaa", 10},
}

func TestSubstrCount(t *testing.T) {
	for _, pair := range tests {
		v := substrCount(pair.n, pair.s)
		if v != pair.expected {
			t.Error(
				"For", pair.s,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

func Equal(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
