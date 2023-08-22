package main

import (
	"fmt"
)

func alternatingCharacters(st string) int32 {

	var replacedTotal int32 = 0
	last := st[0]

	for i := 1; i < len(st); i++ {

		if last == st[i] {
			replacedTotal++

		}
		last = st[i]

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
