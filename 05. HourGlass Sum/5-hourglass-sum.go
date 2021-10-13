package main

import (
	"fmt"
)

func hourglassSum(arr [][]int32) int32 {

	var res int32 = -63

	hgs := func(arr *[][]int32, i, j int32) int32 {

		var sum int32 = 0

		sum =
			(*arr)[i-1][j-1] +
				(*arr)[i-1][j] +
				(*arr)[i-1][j+1] +
				(*arr)[i][j] +
				(*arr)[i+1][j-1] +
				(*arr)[i+1][j] +
				(*arr)[i+1][j+1]

		return sum
	}

	max := func(a, b int32) int32 {
		if a > b {
			return a
		}
		return b
	}

	var i, j int32 = 1, 1

	for i = 1; i < int32(len(arr)-1); i++ {
		for j = 1; j < int32(len(arr[i])-1); j++ {

			res = max(res, hgs(&arr, i, j))
		}
	}

	return res

}

func main() {

	fmt.Print("Expected ", 19, " ")
	fmt.Println(hourglassSum([][]int32{
		{1, 1, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 0},
		{1, 1, 1, 0, 0, 0},
		{0, 0, 2, 4, 4, 0},
		{0, 0, 0, 2, 0, 0},
		{0, 0, 1, 2, 4, 0},
	}))

	fmt.Print("Expected ", 0, " ")
	fmt.Println(hourglassSum([][]int32{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}))
	fmt.Print("Expected ", -6, " ")
	fmt.Println(hourglassSum([][]int32{
		{-1, -1, 0, -9, -2, -2},
		{-2, -1, -6, -8, -2, -5},
		{-1, -1, -1, -2, -3, -4},
		{-1, -9, -2, -4, -4, -5},
		{-7, -3, -3, -2, -9, -9},
		{-1, -3, -1, -2, -4, -5},
	}))

}
