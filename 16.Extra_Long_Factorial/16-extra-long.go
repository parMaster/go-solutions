package main

import (
	"fmt"
	"math/big"
)

func extraLongFactorials(n int32) {

	fact := new(big.Int).SetInt64(1)
	for i := 1; i <= int(n); i++ {
		fact = fact.Mul(fact, new(big.Int).SetInt64(int64(i)))
	}
	fmt.Println(fact)
}

func main() {

	extraLongFactorials(300)
}
