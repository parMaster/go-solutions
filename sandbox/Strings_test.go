package sandbox

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_strings(t *testing.T) {

	n := 100

	input := strconv.Itoa(n)
	fmt.Printf("strconv.Itoa(n): %s of type %s\n", input, reflect.TypeOf(input))

	input = strconv.FormatInt(int64(n), 10)
	fmt.Printf("strconv.FormatInt(int64(n), 10): %s of type %s\n", input, reflect.TypeOf(input))

	// input = string(n) // build failed with error: conversion from int to string yields a string of one rune, not a string of digits (did you mean fmt.Sprint(x)?)
	// fmt.Printf("string(n): %s of type %s\n", input, reflect.TypeOf(input))

	// EqualFold is case-insensitive
	assert.True(t, strings.EqualFold("sTrInG", "string"))

	// Index is case-sensitive
	assert.Equal(t, 4, strings.Index("chicken", "ken"))
	assert.Equal(t, 0, strings.Index("chicken", "ch"))

	// -1 means not found
	assert.Equal(t, -1, strings.Index("chicken", "cH"))

	split := strings.Split("abcd efg", "")
	assert.Equal(t, []string{"a", "b", "c", "d", " ", "e", "f", "g"}, split)

	split = strings.Split("abcd efg", " ")
	assert.Len(t, split, 2)
	assert.Equal(t, []string{"abcd", "efg"}, split)

	rep := strings.Replace("abcd efg", " ", "-", -1)
	assert.Equal(t, "abcd-efg", rep)

}

func change(s []string) {
	s[0] = "Change_function"
}

func Test_arrays(t *testing.T) {

	a := [4]string{"Zero", "One", "Two", "Three"}
	fmt.Println("a:", a)
	// a: [Zero One Two Three]
	fmt.Printf("addr of a[0]: %p\n", &a[0])

	var S0 = a[0:1]
	fmt.Printf("addr of S0[0]: %p\n", &S0[0])

	fmt.Println(S0)
	// [Zero]
	S0[0] = "S0"

	var S12 = a[1:3]
	fmt.Println(S12)
	// [One Two]

	S12[0] = "S12_0"
	S12[1] = "S12_1"

	fmt.Println("a:", a)
	// a: [S0 S12_0 S12_1 Three]

	// Changes to slice -> changes to array
	change(S12)
	fmt.Println("a:", a)
	// a: [S0 Change_function S12_1 Three]

	// capacity of S0
	fmt.Println("Capacity of S0:", cap(S0), "Length of S0:", len(S0))
	fmt.Printf("addr of S0[0]: %p\n", &S0[0])

	// Adding 4 elements to S0
	S0 = append(S0, "N1")
	fmt.Printf("addr of S0[0]: %p\n", &S0[0])
	S0 = append(S0, "N2")
	fmt.Printf("addr of S0[0]: %p\n", &S0[0])
	S0 = append(S0, "N3")
	fmt.Printf("addr of S0[0]: %p\n", &S0[0])
	a[0] = "-N1"

	S0 = append(S0, "N4")
	// will exceed capacity of S0, so a new underlying array will be created
	fmt.Println("Capacity of S0:", cap(S0), "Length of S0:", len(S0))
	fmt.Printf("addr of S0[0]: %p\n", &S0[0])
	// Not the same underlying array anymore!

	// This change does not go to S0
	a[0] = "-N1-"

	// This change does not go to S12
	a[1] = "-N2-"

	fmt.Println("S0:", S0)
	fmt.Println("a: ", a)

	// S12 still points to the same underlying array
	fmt.Println("S12:", S12)
	fmt.Printf("addr of S12[0] (== a[1]): %p\n", &S12[0])

}
