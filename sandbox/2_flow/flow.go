package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// Sqrt implements iterative square root
func Sqrt(x float64) float64 {

	z := float64(1)

	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Println(z)
	}

	return z
}

func returnA() string {

	fmt.Println("returnA Called!")

	return "A"
}

func main() {

	sum := 0
	for i := 0; i < 10; i++ {
		// fmt.Println(i)
		sum += i
	}

	fmt.Println(sum)

	// init and post statements are optional
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)

	// for instead of while )))
	sum = 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println(sum)

	for rand.Intn(91) < 90 { // <- it's a WHILE
		fmt.Println("nope")
	}

	if sum > 100 {
		fmt.Println(sum)
	}

	if sum := 999; sum < 1000 { // evaluate sum := 999 before sum < 1000
		fmt.Println(sum)
	} else {
		fmt.Println(sum) // sum is available inside ELSE statement
	}

	sqrt2 := Sqrt(2)
	fmt.Println("Sqrt result = ", sqrt2)

	//Switch case
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}

	// var_dump in Go ??
	fmt.Printf("\n\r %+v", sqrt2) // %v the value in a default format. when printing structs, the plus flag (%+v) adds field names
	fmt.Printf("\n\r %#v", sqrt2) // %#v a Go-syntax representation of the value

	fmt.Println()

	i := "0"
	switch i {
	case "0":
		fmt.Println("i == 0, returnA not called")
	case returnA(): // Not evalueted nor called
	}

	switch i {
	case "1":
	case returnA(): // Evalueted and called
	}

	defer fmt.Println("One more thing")

	// Same as "switch true" - can be a clean way to write long if-then-else chains.
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}

}
