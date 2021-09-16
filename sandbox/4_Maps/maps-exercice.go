package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {

	var wc = make(map[string]int)

	for _, word := range strings.Fields(s) {
		wc[word]++
	}

	return wc
}

func main() {
	wc.Test(WordCount)
}
