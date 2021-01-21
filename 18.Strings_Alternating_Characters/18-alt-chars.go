package main

import (
	"fmt"
	s "strings"
)

func alternatingCharacters(st string) int32 {

	replaced := 1
	var replacedTotal int32

	for 1 == replaced {

		replaced = 0
		if s.Contains(st, "AA") {
			st = s.Replace(st, "AA", "A", 1)
			replaced = 1
			replacedTotal++
		}
		if s.Contains(st, "BB") {
			st = s.Replace(st, "BB", "B", 1)
			replaced = 1
			replacedTotal++
		}
	}

	return replacedTotal

}

func main() {

	fmt.Println(alternatingCharacters("AAAA"))
	fmt.Println(alternatingCharacters("BBBBB"))
	fmt.Println(alternatingCharacters("ABABABAB"))
	fmt.Println(alternatingCharacters("BABABA"))
	fmt.Println(alternatingCharacters("AAABBB"))

}
