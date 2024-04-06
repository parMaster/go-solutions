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
