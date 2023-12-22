package main

import (
	"fmt"
	"math"
)

func jumpingOnClouds(c []int32) int32 {

	var res int32 = 0
	var zeroes int32 = 0

	for i, v := range c {

		if 0 == v {
			zeroes++
		}

		if (1 == v) || (i == (len(c) - 1)) {

			res += int32(math.Floor(float64(zeroes) / 2.0))

			if 1 == v {
				res++
			}

			zeroes = 0

		}

	}

	return res
}

func main() {

	fmt.Print("Expected ", 1, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0})) // 1
	// os.Exit(1)

	fmt.Print("Expected ", 1, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0})) // 1
	fmt.Print("Expected ", 1, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0})) // 1

	fmt.Println()
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 0})) // 2
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 1, 0})) // 2
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 0})) // 2

	fmt.Println()
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 0, 0})) // 2
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 1, 0, 0})) // 3
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 1, 0})) // 2
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 0, 0})) // 2
	fmt.Print("Expected ", 2, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 1, 0})) // 2

	fmt.Println()
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 0, 0, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 0, 1, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 1, 0, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 1, 0, 1, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 1, 0, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 0, 1, 0})) // 3

	fmt.Println()
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 0, 0, 0, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 1, 0, 1, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 1, 0, 1, 0})) // 3

	fmt.Println()
	fmt.Print("Expected ", 4, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 1, 0, 0, 1, 0})) // 4
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 0, 0, 1, 0})) // 3
	fmt.Print("Expected ", 3, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 0, 0, 1, 0})) // 3
	fmt.Print("Expected ", 4, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 0, 1, 0, 0, 0, 1, 0})) // 4 == round(1/2) + round(3/2) + 1
	fmt.Print("Expected ", 5, " ")
	fmt.Println(jumpingOnClouds([]int32{0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0})) // 5 ==  0 + 4 + 1

}
