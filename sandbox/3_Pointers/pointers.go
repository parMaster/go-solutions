package main

import (
	"fmt"
	"math/rand"
)

func main() {

	var p *int

	i := 42
	p = &i

	fmt.Println(*p)
	fmt.Println(&i)
	fmt.Println(p)

	type Vertex struct {
		X int
		Y int
	}

	v1 := Vertex{1, 2}
	fmt.Println(v1)
	v1.X = 15
	fmt.Println(v1)

	var v = Vertex{1, 2}
	p1 := &v
	(*p1).X = 1e9
	fmt.Println(v)

	w := 101
	arr := make([]int, w, 200)
	for i = range arr {
		arr[i] = rand.Intn(110)
	}

	for i, val := range arr {
		// Enumerate elements, i - index, val - value
		fmt.Println("Element ", i, " has value = ", val)

	}

	for _, val := range arr {
		// Enumerate elements, val - value
		fmt.Print(val, " ")

	}

	fmt.Println(arr)
	fmt.Println(arr[0:4])
	fmt.Println(arr[0:int(w/10)])

	// if arr[0:w/10] == arr[:w/10]  // can't compare slices??
	{
		fmt.Println(arr[0 : w/10])
		fmt.Println(arr[:w/10])
	}
	fmt.Println(arr[w-10:])
	fmt.Println(len(arr[w-10:]))
	fmt.Println(cap(arr[w/2:]))

	s := []int{2, 3, 5, 7, 11, 13}
	printSlice("s", s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice("s", s)

	// Extend its length.
	s = s[:4]
	printSlice("s", s)

	// Drop its first two values.
	s = s[2:]
	printSlice("s", s)

	printSlice("arr", arr)
	arr = arr[0:cap(arr)]
	printSlice("arr", arr)

}

func printSlice(name string, s []int) {
	fmt.Printf("%s len=%d cap=%d %v\n", name, len(s), cap(s), s)
}
