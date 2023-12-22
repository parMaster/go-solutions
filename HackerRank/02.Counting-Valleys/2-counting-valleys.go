package main

import (
	"fmt"
)

func countingValleys(steps int32, path string) int32 {
	// Write your code here

	var depth int = 0
	var valleys int32 = 0

	for i := int32(0); i < steps; i++ {

		switch string(path[i]) {
		case "U":
			if -1 == depth {
				valleys++
			}
			depth++
		case "D":
			depth--
		}

	}

	return valleys
}

func main() {

	res := countingValleys(8, "UDDDUDUU")
	fmt.Println(res)

}
