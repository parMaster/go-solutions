package main

import (
	"fmt"
)

// Super is a struct with superstring
type Super struct {
	n        int32
	str      string
	subCount int64
}

// SearchAround finds a number of anagram around `pos` character
func (s *Super) SearchAround(pos int) {
	i, found, specialLetter := 0, int64(1), ""
	for found == 1 {
		found = 0
		i++
		if (0 <= pos-i) &&
			pos+i < len(s.str) &&
			s.str[pos-i] == s.str[pos+i] &&
			(specialLetter == "" || specialLetter == string(s.str[pos-i])) {

			specialLetter = string(s.str[pos-i])
			s.subCount++
			found = 1
		}
	}
}

// SearchForward is trying to construct the longest string of the same chars starting from pos
// Gauss(?) formula to determine the number of possible substring in a string of lenght N:
// SubstrCount = (n*(n+1))/2
// returns new position to continue from (don't use SearchAround for these kind of strings)
func (s *Super) SearchForward(pos int) int {

	i, found := 0, 1
	for found == 1 {
		found = 0
		i++
		if pos+i < len(s.str) &&
			s.str[pos] == s.str[pos+i] {
			found = 1
		}
	}
	s.subCount += int64((i * (i + 1)) / 2)
	return pos + i - 1
}

func substrCount(n int32, s string) int64 {

	superString := Super{n, s, 0}
	for i := 0; i < len(s); i++ {
		i = superString.SearchForward(i)
		superString.SearchAround(i)
	}
	return superString.subCount
}

func main() {
	fmt.Println("")
}
