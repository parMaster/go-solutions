package main

import (
	"math/rand"
)

func main() {
	//
}

func activityNotifications(e []int32, d int32) int32 {

	var result int32 = 0

	// since 1 <= d <= n
	if len(e) == int(d) {
		return result
	}

	for i := d; i < int32(len(e)); i++ {

		if float32(e[i]) >= 2*median(e[i-d:i]) {
			result++
		}
	}

	return result
}

func median(a []int32) float32 {
	a = quicksort(a)
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
