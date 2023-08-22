package main

import (
	"fmt"
)

func countSwaps(a []int32) {

	var swaps int

	swap := func(a int32, b int32) (int32, int32) {
		return b, a
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a)-1; j++ {

			if a[j] > a[j+1] {
				a[j], a[j+1] = swap(a[j], a[j+1])
				swaps++
			}
		}

	}

	fmt.Println("Array is sorted in", swaps, "swaps.")
	fmt.Println("First Element:", a[0])
	fmt.Println("Last Element:", a[len(a)-1])
}

func main() {

}
