package main

import (
	"fmt"
	"math"
	"strings"
)

func repeatedString(s string, n int64) int64 {
	return int64(strings.Count(s, "a")*int(math.Round(float64(n/int64(len(s))))) + strings.Count(s[0:int(n%int64(len(s)))], "a"))
}

func main() {

	fmt.Print("Expected ", 7, " ")
	fmt.Println(repeatedString("aba", 10))
	// os.Exit(1)
}
