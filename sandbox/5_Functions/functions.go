package main

import (
	"fmt"
)

func stub() bool {
	return false
}

//func flushICache(begin, end uintptr) // implemented externally, fuction body ommited

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	fmt.Println(stub())
	fmt.Println(min(1, 2))
}

// func IndexRune(s string, r rune) int {
// 	for i, c := range s {
// 		if c == r {
// 			return i
// 		}
// 	}
// 	// invalid: missing return statement
// 	return false
// }
