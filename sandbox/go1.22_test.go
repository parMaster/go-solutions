package sandbox

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strings"
	"testing"

	"maps"

	"github.com/stretchr/testify/assert"
)

func Test_MinMax(t *testing.T) {
	// t.Fatal("not implemented")

	var x, y int

	y = 5

	m := min(x) // m == x
	assert.Equal(t, m, x)

	m = min(x, y) // m is the smaller of x and y
	assert.Equal(t, m, x)

	m = max(x, y) // m is the larger of x and y
	assert.Equal(t, m, y)

	m = max(x, y, 10) // m is the larger of x and y but at least 10
	assert.Equal(t, m, 10)

	// For numeric arguments, assuming all NaNs are equal, min and max are commutative and associative:
	assert.Equal(t, min(x, y), min(y, x))
	assert.Equal(t, max(x, y), max(y, x))

	//min(x, y, z) == min(min(x, y), z) == min(x, min(y, z))
	z := 10
	assert.Equal(t, min(x, y, z), min(min(x, y), z))
	assert.Equal(t, min(x, y, z), min(x, min(y, z)))

	c := max(1, 2.0, 10) // c == 10.0 (floating-point kind)
	assert.Equal(t, c, 10.0)
	assert.IsType(t, c, float64(10.0))

	f := max(0, float32(x)) // type of f is float32
	assert.Equal(t, f, float32(0))
	assert.IsType(t, f, float32(0))

	// var s []string
	// _ = min(s...)               // invalid: slice arguments are not permitted
	mx := max("", "foo", "bar") // mx == "foo" (string kind)
	assert.Equal(t, mx, "foo")

	// strings
	assert.Equal(t, min("a", "b", "c"), "a")
	assert.Equal(t, max("a", "b", "c"), "c")

	assert.Equal(t, min("foo", "bar"), min("bar", "foo"))
	assert.Equal(t, max("foo", "bar"), max("bar", "foo"))

	slices.Contains([]string{"a", "b", "c"}, "a")
}

func Test_clear(t *testing.T) {
	//

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	assert.Equal(t, m["a"], 1)
	assert.Equal(t, m["b"], 2)
	assert.Equal(t, m["c"], 3)
	assert.Equal(t, len(m), 3)

	clear(m)

	log.Printf("m: %v", m) // "m: map[]" (the map has been cleared)
	v, ok := m["a"]
	log.Printf("m['a']: %v", v) // "m['a']: 0" (note: ok is false)
	assert.Equal(t, ok, false)
	assert.Equal(t, v, 0)

	log.Printf("m['f']: %v", m["f"]) // "m['f']: 0"
	assert.Equal(t, m["a"], 0)
	assert.Equal(t, m["b"], 0)
	assert.Equal(t, m["c"], 0)
	assert.Equal(t, len(m), 0)

	s := []int{1, 2, 3}

	assert.Equal(t, s[0], 1)
	assert.Equal(t, s[1], 2)
	assert.Equal(t, s[2], 3)
	log.Printf("&s[0]: %v", &s[0])

	clear(s)

	log.Printf("s: %v", s) // "s: [0 0 0]"
	log.Printf("&s[0]: %v", &s[0])
	assert.Equal(t, s[0], 0)
	assert.Equal(t, s[1], 0)
	assert.Equal(t, s[2], 0)
	assert.Equal(t, len(s), 3)
	assert.Equal(t, cap(s), 3)
}

// Test for values
func Test_forValues(t *testing.T) {
	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
		// v := v // create a new 'v'. - not necessary since 1.21 with GOEXPERIMENT=loopvar environment variable
		go func() {
			fmt.Println(v) // nolint: // loop variable v captured by func literal
			done <- true
		}()
	}

	// wait for all goroutines to complete before exiting
	for range values {
		<-done
	}
}

// New testing/slogtest package
// The new testing/slogtest package can help to validate slog.Handler implementations.

// New slices package
// The new slices package provides many common operations on slices, using generic functions that work with slices of any element type.

// func BinarySearch(x S, target E) (int, bool)
// func BinarySearchFunc(x S, target T, cmp func(E, T) int) (int, bool)
// func Clip(s S) S
// func Clone(s S) S
// func Compact(s S) S
// func CompactFunc(s S, eq func(E, E) bool) S
// func Compare(s1, s2 S) int
// func CompareFunc(s1 S1, s2 S2, cmp func(E1, E2) int) int
// func Contains(s S, v E) bool
// func ContainsFunc(s S, f func(E) bool) bool
// func Delete(s S, i, j int) S
// func DeleteFunc(s S, del func(E) bool) S
// func Equal(s1, s2 S) bool
// func EqualFunc(s1 S1, s2 S2, eq func(E1, E2) bool) bool
// func Grow(s S, n int) S
// func Index(s S, v E) int
// func IndexFunc(s S, f func(E) bool) int
// func Insert(s S, i int, v ...E) S
// func IsSorted(x S) bool
// func IsSortedFunc(x S, cmp func(a, b E) int) bool
// func Max(x S) E
// func MaxFunc(x S, cmp func(a, b E) int) E
// func Min(x S) E
// func MinFunc(x S, cmp func(a, b E) int) E
// func Replace(s S, i, j int, v ...E) S
// func Reverse(s S)
// func Sort(x S)
// func SortFunc(x S, cmp func(a, b E) int)
// func SortStableFunc(x S, cmp func(a, b E) int)

func Test_BinSearch(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	unsorted := []int{4, 3, 2, 1, 5, 8, 7, 6}

	// IsSorted reports whether s is sorted.
	assert.True(t, slices.IsSorted(s))
	assert.True(t, slices.IsSortedFunc(s, func(a, b int) int { return a - b }))

	// IsSorted reports whether s is sorted.
	assert.False(t, slices.IsSorted(unsorted))
	assert.False(t, slices.IsSortedFunc(unsorted, func(a, b int) int { return a - b }))

	// BinarySearch returns the smallest index i such that x[i] >= target.
	position, ok := slices.BinarySearch(s, 1)
	assert.Equal(t, position, 0)
	assert.Equal(t, ok, true)

	a := []int{1, 2}
	a = append(a, 3)
	assert.Equal(t, a, []int{1, 2, 3})
	assert.Equal(t, len(a), 3)
	assert.Equal(t, cap(a), 4) // cap doubles when appending to a slice

	// Clone returns a copy of s.
	b := slices.Clone(a)
	assert.Equal(t, b, []int{1, 2, 3})
	assert.Equal(t, len(b), 3)
	assert.Equal(t, cap(b), 3) // cap mutates to len when cloning a slice

	// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
	a = slices.Clip(a)
	assert.Equal(t, a, []int{1, 2, 3})
	assert.Equal(t, len(a), 3)
	assert.Equal(t, cap(a), 3)

	assert.Equal(t, a, []int{1, 2, 3})
	assert.Equal(t, len(a), 3)
	assert.Equal(t, cap(a), 3)

	// Contains reports whether x is within s.
	assert.True(t, slices.Contains(s, 1))

	// Max returns the maximum element in s.
	assert.Equal(t, slices.Max(s), 8)

	// Reverse reverses the elements of s.
	slices.Reverse(s)
	assert.Equal(t, s, []int{8, 7, 6, 5, 4, 3, 2, 1})

	// Sort sorts s.
	slices.Sort(s)
	assert.Equal(t, s, []int{1, 2, 3, 4, 5, 6, 7, 8})

	// Compact returns a slice containing all elements in s that satisfy f.
	// The order of elements is preserved.
	bloated := []int{1, 1, 2, 2, 2, 2, 3, 4, 5, 6, 7, 8}
	compacted := slices.Compact(bloated)
	assert.Equal(t, compacted, []int{1, 2, 3, 4, 5, 6, 7, 8})

}

// New maps package
// The new maps package provides several common operations on maps, using generic functions that work with maps of any key or element type.
//
// func Clone(m M) M
// func Copy(dst M1, src M2)
// func DeleteFunc(m M, del func(K, V) bool)
// func Equal(m1 M1, m2 M2) bool
// func EqualFunc(m1 M1, m2 M2, eq func(V1, V2) bool) bool
func Test_Maps(t *testing.T) {

	m := map[string]string{"a": "a", "b": "b", "c": "c"}

	m1 := maps.Clone(m)

	assert.Equal(t, m, m1)

	m1["c"] = "ddd"

	assert.NotEqual(t, m, m1)

	maps.Copy(m1, m)

	assert.Equal(t, m, m1)

	maps.DeleteFunc(m1, func(k string, v string) bool { return k == "c" || v == "ddd" })

	assert.Equal(t, map[string]string{"a": "a", "b": "b"}, m1)

	m2 := maps.Clone(m)

	assert.True(t, maps.Equal(m, m2))
	assert.False(t, maps.Equal(m, m1))

	for k, v := range m2 {
		m2[k] = strings.ToUpper(v)
	}

	// assert.True(t, maps.EqualFunc(m, m2, func(v string, v1 string) bool { return strings.ToLower(v) == strings.ToLower(v1) }))
	// use EqualFold instead!!
	assert.True(t, maps.EqualFunc(m, m2, func(v string, v1 string) bool { return strings.EqualFold(v, v1) }))
}

// New cmp package
// The new cmp package defines the type constraint Ordered
// and two new generic functions Less and Compare that are useful with ordered types.

func Test_Cmp(t *testing.T) {
	// func Compare[T Ordered](x, y T) int
	// Compare returns
	// -1 if x is less than y,
	//  0 if x equals y,
	// +1 if x is greater than y.

	assert.Equal(t, -1, cmp.Compare(1, 2))
	assert.Equal(t, 0, cmp.Compare(1, 1))
	assert.Equal(t, 1, cmp.Compare(2, 1))

	assert.Equal(t, -1, cmp.Compare("abc", "bcd"))
	assert.Equal(t, 0, cmp.Compare("aaa", "aaa"))
	assert.Equal(t, 1, cmp.Compare("bbb", "aaa"))

	// func Less[T Ordered](x, y T) bool
	// Less reports whether x is less than y. For floating-point types,
	// a NaN is considered less than any non-NaN, and -0.0 is not less than (is equal to) 0.0.

	assert.True(t, cmp.Less("aaa", "bbb"))
	assert.False(t, cmp.Less("bbb", "aaa"))

	// func Or[T comparable](vals ...T) T
	// Or returns the first of its arguments that is not equal to the zero value.
	// If no argument is non-zero, it returns the zero value.

	assert.Equal(t, "aaa", cmp.Or("", "aaa", "bbb"))
	assert.Equal(t, "bbb", cmp.Or("", "bbb", "aaa"))

	// slices?
	// assert.Equal(t, []string{"a"}, cmp.Or([]string{}, []string{"a"}))
	// NO :(

	assert.Equal(t, 1, cmp.Or(0, 1, 2))

	type c struct {
		a string
		b int
	}

	// Structs?
	empty := c{}
	a := c{a: "aaa", b: 1}
	b := c{a: "bbb", b: 1}

	assert.Equal(t, a, cmp.Or(empty, a, b))
	// YES!
}
