package main

import "fmt"

func countTriplets(arr []int64, r int64) int64 {

	r1, r2, r3 := r, r*r, r*r*r

	var nr0, nr1, nr2, nr3 int
	var res int64

	// Look for r3, r2, r1
	for _, v := range arr {

		if 0 == (v % r3) {
			nr3++
			continue
		}
		if 0 == (v % r2) {
			nr2++
			continue
		}
		if 0 == (v % r1) {
			nr1++
		}
	}

	res = int64(nr3 * nr2 * nr1)

	nr0, nr1, nr2 = 0, 0, 0

	// Look for r2, r1, r0
	for _, v := range arr {

		if (0 == (v % r2)) && (v < r3) {
			nr2++
			continue
		}
		if (0 == (v % r1)) && (v < r3) {
			nr1++
			continue
		}
		if 1 == v {
			nr0++
		}

	}

	res += int64(nr2 * nr1 * nr0)

	return res

}

func main() {

	// fmt.Println(countTriplets([]int64{1, 2, 2, 4}, 2))
	// fmt.Println(countTriplets([]int64{1, 3, 9, 9, 27, 81}, 3))
	// fmt.Println(countTriplets([]int64{1, 5, 5, 25, 125}, 5))
	fmt.Println(countTriplets([]int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 1))

}
