package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// back of the envelope solution, while waiting for MacOS to update
func makingAnagrams(s1 string, s2 string) int32 {

	var deletions int32

	// frequency array
	m := make(map[byte]int, 26)

	for i := 0; i < len(s1); i++ {
		_, ok := m[s1[i]]
		if !ok {
			m[s1[i]] = 0
		}
		m[s1[i]]++
	}

	for i := 0; i < len(s2); i++ {
		_, ok := m[s2[i]]
		if !ok {
			m[s2[i]] = 0
		}
		m[s2[i]]--
	}

	for _, f := range m {
		deletions += int32(math.Abs(float64(f)))
	}

	return deletions
}

func Test_makingAnagrams(t *testing.T) {

	assert.Equal(t, int32(6), makingAnagrams("abc", "amnop"))
	assert.Equal(t, int32(4), makingAnagrams("cde", "abc"))
}
