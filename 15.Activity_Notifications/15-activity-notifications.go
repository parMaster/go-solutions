package main

import (
	"math"
	"math/rand"
)

func main() {
	//
}

func activityNotifications(e []int32, d int32) int32 {

	var result int32 = 0
	// countsort first d values
	cs := make([]int, 201)
	for i := 0; i < int(d); i++ {
		cs[int(e[i])]++
	}

	for i := int(d); i < len(e); i++ {

		dm := doubleMedian(cs, d)
		if e[i] >= dm {
			result++
		}

		cs[int(e[i])]++
		cs[e[i-int(d)]]--

	}

	return result
}

// Let's find a double median in countsorted array, since it's a problem`s condition
// and I don't need to use float - 2*(a/2) is always int!
func doubleMedian(cs []int, d int32) int32 {

	var result int32 = 0
	var thereYet int = 0

	halfD := int(math.Floor(float64(d / 2.0)))

	if 0 == d%2 {
		for i := 0; i < 201; i++ {
			thereYet += cs[i]
			if 0 == result && thereYet >= halfD {
				result = int32(i)
			}
			if thereYet >= halfD+1 {
				return result + int32(i)
			}
		}
	}

	for i := 0; i < 201; i++ {
		thereYet += cs[i]
		if thereYet > halfD {
			return int32(2 * i)
		}
	}

	return 0
}

// Not used, just don't want to delete tests, so I'll leave it here

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
