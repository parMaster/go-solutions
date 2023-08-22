package main

import (
	"fmt"
	s "strings"
)

func makeAnagram(a string, b string) int32 {

	// res := int32(len(a) + len(b))

	for i := 0; i < len(a); i++ {

		chr := string(a[i])
		if s.Contains(b, chr) {

			a = s.Replace(a, chr, "", 1)
			b = s.Replace(b, chr, "", 1)
			i--
		}

	}

	return int32(len(a) + len(b))

}

func main() {

	fmt.Println(makeAnagram("cde", "dcf"))
	fmt.Println(makeAnagram("cde", "abc"))
}
