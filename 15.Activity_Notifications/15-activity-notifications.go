package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	//

	x := []int{1, 2, 3, 3, 3, 3, 4, 5}
	pos := sort.SearchInts(x, 3)

	fmt.Println(pos)
	fmt.Println(x[:pos])
	fmt.Println(x[pos:len(x)])

}

func activityNotifications(e32 []int32, d int32) int32 {

	var result int32 = 0
	var e []int
	// var max int32

	priorDays := quicksort(e32[0:d])

	// sort.SearchInts takes []int    8-/
	e = make([]int, len(priorDays))
	for i := range priorDays {
		e[i] = int(priorDays[i])
	}

	// since 1 <= d <= n
	if len(e32) == int(d) {
		return result
	}

	for i := d; i < int32(len(e32)); i++ {

		if e32[i] >= int32(2*median(e)) {

			// fmt.Println("Element e32[", i, "] = ", e32[i], " med = ", median(e), " 2*med = ", 2*median(e), " len(e) = ", len(e))

			result++
		}

		e = delete(e, e32[i-d])
		e = insert(e, int(e32[i]))
	}
	// fmt.Println(result)
	return result
}

func delete(a []int, e32 int32) []int {

	e := int(e32)

	pos := sort.SearchInts(a, e)

	if 0 == pos {
		return a[1:len(a)]
	}

	if len(a) == pos {
		return a[:len(a)]
	}

	return append(a[:pos], a[pos+1:]...)
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
