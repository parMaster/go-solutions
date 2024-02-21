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

func BenchmarkCalc0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := calc0("mil_c10000_l1000000.csv")
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkCalc1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := calc1("mil_c10000_l1000000.csv")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func floatStrToInt(s string) int {

	if i, ok := floatMap[s]; ok {
		return i
	} else {
		panic("floatStrToInt: " + s)
	}
}

func TestFloatStrToInt(t *testing.T) {
	initFloatMap()

	require.Equal(t, 0, floatStrToInt("0.0"))
	require.Equal(t, 69, floatStrToInt("6.9"))
	require.Equal(t, 42, floatStrToInt("4.2"))
	require.Equal(t, 420, floatStrToInt("42.0"))
	require.Equal(t, 0, floatStrToInt("0.0"))
	require.Equal(t, 0, floatStrToInt("-0.0"))
	require.Equal(t, -69, floatStrToInt("-6.9"))
	require.Equal(t, -42, floatStrToInt("-4.2"))
	require.Equal(t, -420, floatStrToInt("-42.0"))

}
