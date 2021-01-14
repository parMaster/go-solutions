package main

import (
	"fmt"
)

func rotLeft(a []int32, d int32) []int32 {

	return append(a[d%int32(len(a)):], a[0:d%int32(len(a))]...)

}

func main() {

	fmt.Print("Expected ", "77 97 58 1 86 58 26 10 86 51 41 73 89 7 10 1 59 58 84 77", " ")
	fmt.Println(rotLeft([]int32{41, 73, 89, 7, 10, 1, 59, 58, 84, 77, 77, 97, 58, 1, 86, 58, 26, 10, 86, 51}, 10))

	fmt.Print("Expected ", "5 1 2 3 4", " ")
	fmt.Println(rotLeft([]int32{1, 2, 3, 4, 5}, 9))

	fmt.Print("Expected ", "5 1 2 3 4", " ")
	fmt.Println(rotLeft([]int32{1, 2, 3, 4, 5}, 4))

	fmt.Print("Expected ", "1 2 3 4 5", " ")
	fmt.Println(rotLeft([]int32{1, 2, 3, 4, 5}, 5))

	fmt.Print("Expected ", "1 2 3 4 5", " ")
	fmt.Println(rotLeft([]int32{1, 2, 3, 4, 5}, 0))

}
