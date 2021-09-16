package main

import (
	"fmt"
	"math/cmplx"
)

func main() {

	var (
		ToBe   bool       = false
		MaxInt uint64     = 1<<64 - 1
		z      complex128 = cmplx.Sqrt(-5 + 12i)
	)

	fmt.Printf("ToBe = %t and type is %T \n\r", ToBe, ToBe)
	fmt.Printf("MaxInt = %d and type is %T \n\r", MaxInt, MaxInt)
	fmt.Printf("z = %g and type is %T \n\r\n\r", z, z)

	// autoassign type
	var autoTypedInt = 42
	fmt.Println(autoTypedInt)
	fmt.Printf("autoTypedInt type is %T \n\r", autoTypedInt)

	var autoTypedInt64 = 1<<63 - 1 // 63 because signed int is autoassigned
	fmt.Println(autoTypedInt64)
	fmt.Printf("autoTypedInt64 type is %T \n\r", autoTypedInt)

	// define typed maxInt
	var maxInt uint64 = 1<<64 - 1
	fmt.Println(maxInt)

	// shortMaxInt := 1<<64 - 1 - error
	shortMaxInt := uint64(1<<64 - 1)
	fmt.Println(shortMaxInt)

	sum := 0
	i := 10 // function scope `i`
	for i := 0; i < 10; i++ {
		sum += i
		fmt.Println(i) // for loop scope `i`
	}
	fmt.Println(sum)
	fmt.Println(i)
}
