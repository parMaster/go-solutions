package main

import (
	"fmt"
	"math"
)

// Complete the sockMerchant function below.
func sockMerchant(n int32, ar []int32) int32 {

	var sum int32

	var pc = make(map[int32]int)
	for _, v := range ar {
		pc[v]++
	}

	// fmt.Println(pc)

	for _, v := range pc {
		sum += int32(math.Floor(float64(v / 2)))
	}

	return sum

}

func main() {

	res := sockMerchant(9, []int32{10, 20, 20, 10, 10, 30, 50, 10, 20})
	fmt.Println(res)

	res = sockMerchant(10, []int32{1, 1, 3, 1, 2, 1, 3, 3, 3, 3})
	fmt.Println(res)

}
