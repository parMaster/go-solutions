package sandbox

import (
	"strings"
	"testing"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {

	var wc = make(map[string]int)

	for _, word := range strings.Fields(s) {
		wc[word]++
	}

	return wc
}

func TestMapsExercice(t *testing.T) {
	wc.Test(WordCount)
}
