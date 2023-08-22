package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func pangrams(s string) string {
	s = strings.ToLower(s)

	m := make(map[byte]int)
	for i := 0; i <= len(s)-1; i++ {

		if s[i] == 32 { // filtering out spaces
			continue
		}
		_, ok := m[s[i]]
		if ok {
			m[s[i]]++
		} else {
			m[s[i]] = 1
		}
	}

	if len(m) == 26 {
		return "pangram"
	} else {
		return "not pangram"
	}
}

func Test_pangrams(t *testing.T) {
	assert.Equal(t, "pangram", pangrams("We promptly judged antique ivory buckles for the next prize"))
}
