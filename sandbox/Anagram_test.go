package sandbox

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isAnagram(s1, s2 string) bool {

	for _, v := range s1 {
		if string(v) != " " && strings.Count(s1, string(v)) != strings.Count(s2, string(v)) {
			return false
		}
	}
	for _, v := range s2 {
		if string(v) != " " && strings.Count(s1, string(v)) != strings.Count(s2, string(v)) {
			return false
		}
	}

	return true
}

func Test_Anagram(t *testing.T) {
	assert.True(t, isAnagram("schoolmaster", "the classroom"), "the classroom")
	assert.True(t, isAnagram("punishment", "nine thumps"), "nine thumps")
	assert.True(t, isAnagram("debit card", "bad credit"), "bad credit")
	assert.False(t, isAnagram("debit card", "good credit"), "good credit")
}

// how many characters to delete from both strings to get anagram
// insane person implementation with only one minimal hashmap
func diffCount(s1, s2 string) int {
	if len(s2) > len(s1) {
		s1, s2 = s2, s1
	}

	m := map[string]int{}
	for _, r := range s1 {
		s := string(r)
		if _, ok := m[s]; !ok {
			m[s] = strings.Count(s1, s) - strings.Count(s2, s)
		}
	}
	sum := 0
	for _, count := range m {
		sum += count
	}
	diff := len(s2) - (len(s1) - sum) + sum

	return diff
}

func Test_MinAnagram(t *testing.T) {
	assert.Equal(t, 6, diffCount("hello", "billion"))
	assert.Equal(t, 7, diffCount("hello", "billlion"))
	assert.Equal(t, 0, diffCount("llo", "lol"))
}
