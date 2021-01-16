package main

import (
	"fmt"
	"strings"
)

func twoStrings(s1 string, s2 string) string {

	for i := 0; i < len(s1)-1; i++ {

		if -1 != strings.Index(s2, string(s1[i])) {
			return "YES"
		}
	}

	return "NO"
}

func main() {

	fmt.Println(twoStrings("and", "art"))
	fmt.Println(twoStrings("be", "cat"))

}
