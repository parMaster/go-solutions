package main

import "fmt"

func checkMagazine(magazine []string, note []string) {

	m := make(map[string]int, 0)

	for _, mWord := range magazine {
		m[mWord]++
	}

	for _, nWord := range note {

		if 0 == m[nWord] {
			fmt.Println("No")
			return
		}
		m[nWord]--

	}
	fmt.Println("Yes")
	return

}

func main() {

	checkMagazine([]string{"give", "me", "one", "grand", "today", "night"}, []string{"give", "one", "grand", "today"})
}
