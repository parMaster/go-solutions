package main

import (
	"fmt"
)

// Super is a struct with superstring
type Super struct {
	n        int32
	str      string
	subCount int64
	specials map[string]int
}

// SearchAround is method
func (s *Super) SearchAround(pos int) int64 {

	i, found, letterAround := 0, 1, ""
	s.subCount++
	for 1 == found {
		found = 0
		i++
		if (0 <= pos-i) &&
			pos+i < len(s.str) &&
			s.str[pos-i] == s.str[pos+i] &&
			(letterAround == "" || letterAround == string(s.str[pos-i])) {

			letterAround = string(s.str[pos-i])
			// builtString = string(s.str[pos-i]) + builtString + string(s.str[pos+i])
			// fmt.Println(builtString)
			// _, ok := s.specials[builtString]
			// if !ok {
			s.subCount++
			// s.specials[builtString] = 1
			// }
			found = 1

		}
	}
	return 0
}

// SearchForward is method
func (s *Super) SearchForward(pos int) int64 {

	i, found := 0, 1
	for 1 == found {
		found = 0
		i++
		if pos+i < len(s.str) &&
			s.str[pos] == s.str[pos+i] {
			s.subCount++
			found = 1
		}
	}
	return 0
}

func substrCount(n int32, s string) int64 {

	superString := Super{n, s, 0, map[string]int{}}
	for i := 0; i < len(s); i++ {
		superString.SearchAround(i)
		superString.SearchForward(i)
	}

	// fmt.Println(superString.specials)

	return superString.subCount
}

func main() {
	fmt.Println("")
}

// Sample 2 (aaaa) and explanation are invalid
// According to the Problem conditions, the answer should be 12 but not 10:
// 1) all the characters are the same:
// a, a, a, a, aa, aa, aa, aaa, aaa
