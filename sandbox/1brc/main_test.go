package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateAndWrite(t *testing.T) {
	require.NoError(t, genWrite(10, 10, "ten"))
	require.NoError(t, genWrite(10000, 1000000, "mil"))
}

func TestCalc0(t *testing.T) {

	results, err := calc0("ten_c10_l10.csv")
	require.NoError(t, err)
	require.Len(t, results, 4)

	for c, v := range results {
		fmt.Println(c, v.min, v.mean, v.max, v.n)
	}

	results, err = calc0("mil_c10000_l1000000.csv")
	require.NoError(t, err)
	require.Len(t, results, 10000)

	i := 0
	for c, v := range results {
		if i++; i > 10 {
			break
		}
		fmt.Println(c, v.min, v.mean, v.max, v.n)
	}

}

func BenchmarkCalc0RAM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := calc0("/mnt/ramdisk/10mil_c10000_l10000000.csv")
		if err != nil {
			b.Fatal(err)
		}
	}
}
