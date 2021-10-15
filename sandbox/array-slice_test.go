package sandbox

import (
	"fmt"
	"testing"
)

func TestArraySlice(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	fmt.Println(a)
	fmt.Println("a) cap:", cap(a), " len:", len(a))

	b := append(a[:2], a[3:]...)
	fmt.Println(b)
	fmt.Println("b) cap:", cap(b), " len:", len(b))
	fmt.Println(a)
	fmt.Println("a) cap:", cap(a), " len:", len(a))

	t.Error(
		"Test: ",
	)
	// fucking insane
	// [1 2 3 4 5]
	// a) cap: 5  len: 5
	// [1 2 4 5]
	// b) cap: 5  len: 4
	// [1 2 4 5 5]
	// a) cap: 5  len: 5
}
