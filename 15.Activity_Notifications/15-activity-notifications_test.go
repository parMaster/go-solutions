package main

import (
	"testing"
)

type casesMedian struct {
	a        []int32
	expected float32
}

type casesQsort struct {
	arr      []int32
	expected []int32
}

type cases struct {
	e        []int32
	d        int32
	expected int32
}

var testsMedian = []casesMedian{
	{[]int32{1, 2, 3, 4}, float32(2.5)},
	{[]int32{1, 2, 3, 4, 5}, 3},
}

var testsQsort = []casesQsort{
	{[]int32{1, 2, 3, 4, 5}, []int32{1, 2, 3, 4, 5}},
	{[]int32{1, 5, 3, 2, 4}, []int32{1, 2, 3, 4, 5}},
}

var tests = []cases{
	{[]int32{10, 20, 30, 40, 50}, 3, 1},
}

func TestMedian(t *testing.T) {
	for _, pair := range testsMedian {
		v := median(pair.a)
		if v != pair.expected {
			t.Error(
				"For", pair.a,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

func TestQsort(t *testing.T) {
	for _, pair := range testsQsort {
		v := quicksort(pair.arr)
		if !Equal(v, pair.expected) {
			t.Error(
				"For", pair.arr,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

func TestActivityNotifications(t *testing.T) {
	for _, pair := range tests {
		v := activityNotifications(pair.e, pair.d)
		if v != pair.expected {
			t.Error(
				"For", pair.e, pair.d,
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
