package stringset

import (
	"fmt"
	"strings"
)

// Implement Set as a collection of unique string values.
//
// For Set.String, use '{' and '}', output elements as double-quoted strings
// safely escaped with Go syntax, and use a comma and a single space between
// elements. For example, a set with 2 elements, "a" and "b", should be formatted as {"a", "b"}.
// Format the empty set as {}.

// Define the Set type here.

type Set struct {
	m map[string]any
}

func New() Set {
	m := map[string]any{}
	return Set{
		m: m,
	}
}

func NewFromSlice(l []string) Set {
	m := map[string]any{}
	for _, el := range l {
		if _, ok := m[el]; !ok {
			m[el] = struct{}{}
		}
	}
	return Set{
		m: m,
	}
}

func (s Set) String() string {
	elements := []string{}
	for el := range s.m {
		elements = append(elements, fmt.Sprintf("\"%s\"", el))
	}

	return fmt.Sprintf("{%s}", strings.TrimSpace(strings.Join(elements, ", ")))
}

func (s Set) Elements() []string {
	res := []string{}
	for el := range s.m {
		res = append(res, el)
	}
	return res
}

func (s Set) Len() int {
	return len(s.m)
}

func (s Set) IsEmpty() bool {
	return len(s.m) == 0
}

func (s Set) Has(elem string) bool {
	_, ok := s.m[elem]
	return ok
}

func (s Set) Add(elem string) {
	s.m[elem] = struct{}{}
}

// returns true if s1 is a subset of s2
func Subset(s1, s2 Set) bool {

	if s1.IsEmpty() {
		return true
	} else if !s1.IsEmpty() && s2.IsEmpty() {
		return false
	}

	for _, el := range s1.Elements() {
		if !s2.Has(el) {
			return false
		}
	}

	return true
}

// sets are disjoint if they share no elements
func Disjoint(s1, s2 Set) bool {
	if s1.IsEmpty() || s2.IsEmpty() {
		return true
	}
	for _, el := range s1.Elements() {
		if s2.Has(el) {
			return false
		}
	}
	return true
}

// returns true if sets are equal
func Equal(s1, s2 Set) bool {

	// empty sets are equal
	if s1.IsEmpty() && s2.IsEmpty() {
		return true
	}
	// XOR
	if (s1.IsEmpty() && !s2.IsEmpty()) || (!s1.IsEmpty() && s2.IsEmpty()) {
		return false
	}

	// lenght are equal?
	if s1.Len() != s2.Len() {
		return false
	}

	// checking elements
	for _, el := range s1.Elements() {
		if !s2.Has(el) {
			return false
		}
	}

	return true
}

func Intersection(s1, s2 Set) Set {
	res := New()
	for _, el := range s1.Elements() {
		if s2.Has(el) {
			res.Add(el)
		}
	}
	return res
}

// Difference (or Complement) of a set is a set of all elements that are only in the first set
func Difference(s1, s2 Set) Set {
	res := New()
	for _, el := range s1.Elements() {
		if !s2.Has(el) {
			res.Add(el)
		}
	}
	return res
}

func Union(s1, s2 Set) Set {
	res := New()
	for _, el := range s1.Elements() {
		res.Add(el)
	}
	for _, el := range s2.Elements() {
		res.Add(el)
	}
	return res
}
