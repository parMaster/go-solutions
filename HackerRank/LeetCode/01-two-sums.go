package main

import (
	"fmt"
)

func twoSum(nums []int, target int) []int {

	for i, v1 := range nums {
		for j, v2 := range nums {
			if (i != j) && ((v1 + v2) == target) {

				return []int{i, j}
			}
		}

	}

	return []int{0, 0}

}

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
	fmt.Println(twoSum([]int{3, 2, 4}, 6))
	fmt.Println(twoSum([]int{3, 3}, 6))
}
