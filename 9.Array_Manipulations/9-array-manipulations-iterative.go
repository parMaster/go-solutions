package main

import (
	// "fmt"
	"fmt"
	"math"
)

func arrayManipulationIterative(n int32, queries [][]int32) int64 {

	// 24 Points, 7 test cases failed due to timeout

	var max int64 = 0
	res := make([]int64, n)

	for _, q := range queries {

		a, b, k := q[0], q[1], q[2]

		for j := a - 1; j < b; j++ {
			res[j] += int64(k)
		}

	}

	for _, v := range res {
		max = int64(math.Max(float64(max), float64(v)))
	}

	return max
}

func arrayManipulation(n int32, queries [][]int32) int64 {
	// Try #2: Most Wanted Number
	// Timeouts + 3 wrong answers

	var max, result int64 = 0, 0
	res := make([]int64, n+1)

	var i int32 = 1
	var imax int32 = 0

	for i = 1; i <= n+1; i++ {

		for _, q := range queries {

			if i >= q[0] && i <= q[1] {
				res[i]++
			}

		}

	}

	for i, v := range res {
		if max < v {
			max = v
			imax = int32(i)
		}
	}

	for _, q := range queries {

		if int32(imax) >= q[0] && int32(imax) <= q[1] {
			result += int64(q[2])
		}

	}

	return result
}

func main() {

	// fmt.Println(arrayManipulation(10, [][]int32{{1, 5, 3}, {4, 8, 7}, {6, 9, 1}}))
	fmt.Println(arrayManipulation(10, [][]int32{{2, 6, 8}, {3, 5, 7}, {1, 8, 1}, {5, 9, 15}}))

}
