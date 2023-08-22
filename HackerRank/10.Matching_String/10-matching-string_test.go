package main

import (
	"testing"
)

type testpair struct {
	strings  []string
	queries  []string
	expected []int32
}

var tests = []testpair{
	{[]string{"aba", "baba", "aba", "xzxb"}, []string{"aba", "xzxb", "ab"}, []int32{2, 1, 0}},
	{[]string{"ab", "ab", "abc"}, []string{"ab", "abc", "bc"}, []int32{2, 1, 0}},
	{[]string{"def", "de", "fgh"}, []string{"de", "lmn", "fgh"}, []int32{1, 0, 1}},
	{[]string{"abcde", "sdaklfj", "asdjf", "na", "basdn", "sdaklfj", "asdjf", "na", "asdjf", "na", "basdn", "sdaklfj", "asdjf"}, []string{"abcde", "sdaklfj", "asdjf", "na", "basdn"}, []int32{1, 3, 4, 3, 2}},
}

// A test for matchingStrings func
func TestMatchingStrings(t *testing.T) {
	for _, pair := range tests {
		v := matchingStrings(pair.strings, pair.queries)
		if !Equal(v, pair.expected) {
			t.Error(
				"For", pair.strings, pair.queries,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
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
