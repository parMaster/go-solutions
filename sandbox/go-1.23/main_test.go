package main

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
https://tip.golang.org/wiki/RangefuncExperiment

Iterator functions used as range expressions:

func(func() bool)
func(func(K) bool)
func(func(K, V) bool)

Calls of the iterator argument function produce the iteration values for the “for-range” loop.

Consider this function for iterating a slice backwards:
*/
func Backward[E any](s []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}

// It can be invoked as:
func Test_Backward(t *testing.T) {
	s := []string{"hello", "world"}
	for i, x := range slices.Backward(s) {
		fmt.Println(i, x)
		break
	}
}

/*
1. The Backward function:
This is a generic function that takes a slice of any type E and returns a function.
The returned function is the iterator itself.

2. The returned iterator function:
return func(yield func(int, E) bool) {
    // ...
}
This anonymous function is the actual iterator.
It takes a yield function as an argument, which is used to produce each item in the iteration.

3. The iteration logic:

for i := len(s) - 1; i >= 0; i-- {
    if !yield(i, s[i]) {
        return
    }
}

This loop iterates over the slice in reverse order. For each element, it calls the yield function
with the current index and value. If yield returns false, the iteration stops.

This is where the "magic" happens. The range clause can now accept a function that returns an
iterator (in this case, Backward(s)). The Go runtime calls this function and uses the returned
iterator to control the loop.

The key to understanding this is that the range clause in Go 1.23 can work with custom iterators.
When you use range with Backward(s), it's roughly equivalent to this behind the scenes:

iterator := Backward(s)
iterator(func(i int, x string) bool {
    // This is the body of the for loop
    fmt.Println(i, x)
    return true
})

The range clause handles calling the iterator function and passing it a yield function that
assigns values to the loop variables and executes the loop body.
This new feature allows for powerful custom iteration behaviors, such as iterating in reverse order,
as shown in this example.

The yield function isn't defined within the code we see. This is actually one of the key aspects
of how the new iterator feature works in Go 1.23:

The yield function is not defined by us, the programmers. It's provided by the Go runtime when
the iterator is used in a for...range loop.
In the definition func(yield func(int, E) bool), we're just specifying the signature of the
yield function that will be passed in. It takes two parameters (an int and a value of type E)
and returns a boolean.
When you use this iterator in a for...range loop, the Go runtime creates and passes in the yield
function automatically. This function does several things:

It assigns values to the loop variables (in this case, i and x)
It runs the body of the for loop
It returns true if the loop should continue, or false if it should stop
(e.g., if a break statement was encountered)



So, in essence, the yield function is the bridge between your custom iterator and the for...range
loop mechanics. You define how to iterate (in this case, backwards over a slice), and the Go
runtime handles how to feed those values into the loop.
When you write:

for i, x := range Backward(s) {
    fmt.Println(i, x)
}

The Go runtime is essentially doing something like this behind the scenes:

Backward(s)(func(index int, value string) bool {
    i, x = index, value  // Assign to loop variables
    fmt.Println(i, x)    // Run loop body
    return true          // Continue iteration
})

This abstraction allows you to define custom iteration behaviors without having to worry about
how the for...range loop internally works. You just need to provide a function that calls yield
for each item in your custom sequence.

*/

/*
 This program would translate inside the compiler to a program more like:

slices.Backward(s)(func(i int, x string) bool {
    fmt.Println(i, x)
    return true
})

*/

/*
Iterators

The new *iter* package provides the basic definitions for working with user-defined iterators.

The slices package adds several functions that work with iterators:

    All - returns an iterator over slice indexes and values.
    Values - returns an iterator over slice elements.
    Backward - returns an iterator that loops over a slice backward.
    Collect - collects values from an iterator into a new slice.
    AppendSeq - appends values from an iterator to an existing slice.
    Sorted - collects values from an iterator into a new slice, and then sorts the slice.
    SortedFunc - is like Sorted but with a comparison function.
    SortedStableFunc - is like SortFunc but uses a stable sort algorithm.
    Chunk - returns an iterator over consecutive sub-slices of up to n elements of a slice.
*/

func TestIterators(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}

	fmt.Println("Original slice: ", s)

	fmt.Println("\n\n --- slices.Values does []int -> iter.Seq[int]: ")
	// iter.Seq[int] yields the slice elements in order, no indexes.
	for v := range slices.Values(s) {
		fmt.Print(v, " ")
	}

	fmt.Println("\n\n --- slices.All does []int -> iter.Seq2[int, int]: ")
	// iter.Seq2[int, int] yields index-value pairs.
	for i, v := range slices.All(s) {
		fmt.Print(i, ":", v, " ")
		assert.Equal(t, s[i], v)
	}

	fmt.Println("\n\n --- slices.Sorted: ")
	// slices.Sorted takes an iterator, so we need to convert the slice to an iterator first
	// returns a copy of the slice sorted in increasing order, without modifying the original slice.
	for _, x := range slices.Sorted(slices.Values(s)) {
		fmt.Print(x, " ")
	}

	fmt.Println("\n\n --- slices.SortedFunc sorting values high->low: ")
	// slices.SortedFunc is similar to slices.Sorted but allows you to provide a custom comparison function.
	// returns a copy of the slice sorted using the provided comparison function.
	for _, x := range slices.SortedFunc(slices.Values(s), func(a, b int) int {
		if a > b {
			return -1
		}
		if a < b {
			return 1
		}
		return 0
	}) {
		fmt.Print(x, " ")
	}

	fmt.Println("\n\n --- Collect(Values([]int)) does []int -> iter.Seq[int] -> []int: ")
	// Collect collects values from an iterator into a new slice and returns it.
	// returns a new slice with the same elements as the original slice.
	collected := slices.Collect(slices.Values(s))
	fmt.Println(collected)
	assert.Equal(t, s, collected)

	fmt.Println("\n --- AppendSeq([]int, Values([]int)): ")
	// AppendSeq appends the values from an iterator to an existing slice and returns the extended slice.
	// returns a new slice with the original slice followed by the elements from the iterator.
	appended := slices.AppendSeq(s, slices.Values(s))
	fmt.Println(appended)
	// append(s, s...) does the same thing
	assert.Equal(t, append(s, s...), appended)

	fmt.Println("\n --- Chunk(s, 3) - iterator over consecutive sub-slices of up to n elements of s: ")
	for v := range slices.Chunk(s, 3) {
		fmt.Println(v)
	}

}

/*
The maps package adds several functions that work with iterators:

    All - returns an iterator over key-value pairs from a map.
    Keys - returns an iterator over keys in a map.
    Values - returns an iterator over values in a map.
    Insert - adds the key-value pairs from an iterator to an existing map.
    Collect - collects key-value pairs from an iterator into a new map and returns it.
*/

func TestMapIterators(t *testing.T) {

	m := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}

	fmt.Print("Original Map:", m)

	fmt.Println("\n\n--- All returs iter.Seq2[K,V]:")
	for k, v := range maps.All(m) {
		fmt.Printf("%s:%d ", k, v)
	}

	fmt.Println("\n\n--- Keys returns iterator over keys iter.Seq[K]:")
	for k := range maps.Keys(m) {
		fmt.Printf("%s ", k)
	}

	fmt.Println("\n\n--- Values returns iterator over keys iter.Seq[V]:")
	for v := range maps.Values(m) {
		fmt.Printf("%d ", v)
	}

	fmt.Println("\n\n--- Iterating over slices.Sorted(Keys) printing ordered map:")
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%s:%d ", k, m[k])
	}

	m2 := map[string]int{"six": 6}
	fmt.Println("\n\n--- Insert(map, iter.Seq2[K,V]) adds sequence to a map:")
	fmt.Println(" m + {\"six\": 6}:")
	maps.Insert(m, maps.All(m2))
	for k, v := range m {
		fmt.Printf("%s:%d ", k, v)
	}

	fmt.Println("\n\n--- Collect(All(map)) does map -> iter.Seq2[K,V] -> new map:")
	for k, v := range maps.Collect(maps.All(m)) {
		fmt.Printf("%s:%d ", k, v)
	}

	fmt.Println()
}
