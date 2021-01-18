package main

import (
	"math/rand"
	"sort"
)

func main() {
	//

}

func activityNotifications(e32 []int32, d int32) int32 {

	var result int32 = 0
	var e []int

	e32 = append(quicksort(e32[0:d]), e32[d:len(e32)]...)

	// sort.SearchInts takes []int    8-/
	e = make([]int, len(e32))
	for i := range e32 {
		e[i] = int(e32[i])
	}

	// since 1 <= d <= n
	if len(e) == int(d) {
		return result
	}

	for i := d; i < int32(len(e)); i++ {
		if float32(e[i]) >= 2*median(e[i-d:i]) {
			result++
		}

		if i < int32(len(e)-1) {
			e = append(insert(e[i-d+1:i], e[i]), e[i+1:len(e)]...)
		} else {
			e = insert(e[i-d+1:i], e[i])
		}

		i--
	}

	return result
}

func insert(a []int, e int) []int {

	pos := sort.SearchInts(a, e)

	if 0 == pos {
		return append([]int{e}, a...)
	}

	if len(a) == pos {
		return append(a, e)
	}

	return append(a[0:pos], append([]int{e}, a[pos:len(a)]...)...)
}

func median(a []int) float32 {

	if 0 == len(a)%2 {
		return float32((a[int(len(a)/2)])+a[int(len(a)/2)-1]) / 2
	}
	return float32(a[int(len(a)/2)])
}

func quicksort(a []int32) []int32 {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]

	for i := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]

	quicksort(a[:left])
	quicksort(a[left+1:])

	return a
}
