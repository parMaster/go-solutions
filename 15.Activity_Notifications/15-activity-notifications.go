package main

import (
	"math/rand"
)

func activityNotifications(e []int32, d int32) int32 {

	if len(e) == int(d) {
		return 0
	}

	return 0
}

func main() {
	//
}

func median(a []int32) float32 {
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

	for i, _ := range a {
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
