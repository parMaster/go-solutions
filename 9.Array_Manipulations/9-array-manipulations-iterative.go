package main

import (
	// "fmt"
	"fmt"
	"math"
)

func arrayManipulation(n int32, queries [][]int32) int64 {

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

func main() {

	// fmt.Println(arrayManipulation(10, [][]int32{{1, 5, 3}, {4, 8, 7}, {6, 9, 1}}))
	fmt.Println(arrayManipulation(10, [][]int32{{2, 6, 8}, {3, 5, 7}, {1, 8, 1}, {5, 9, 15}}))

}
