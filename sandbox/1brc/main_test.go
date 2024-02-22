package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateAndWrite(t *testing.T) {
	require.NoError(t, genWrite(10, 10, "ten"))
	require.NoError(t, genWrite(10000, 1000000, "mil"))
}

// readLines gives me a baseline result for IO speed
func readLines(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	i := 0
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return i
		}
		line = line[:len(line)-1]
		if len(line) == 0 {
			panic("empty line")
		}
		i++
	}
}

// get a baseline for the time it takes to read the file
func TestReadLines(t *testing.T) {
	if _, err := os.Stat("10mil_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(10000, 10000000, "10mil")
	}

	n := readLines("10mil_c10000_l10000000.csv")
	require.Equal(t, 10000000, n)
	// 0.31s for 10 mils
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

	for range 10000 {
		n := rand.Float64()*199.8 - 99.9
		n = math.Round(n*10) / 10
		require.Equal(t, int(n*10), floatStrToInt(fmt.Sprintf("%.1f", n)))
	}

}

// calc0 is naive single-core solution of 1brc
func calc0(filename string) (results map[string]location, err error) {

	sums := make(map[string]location)

	file, err := os.Open(filename)
	if err != nil {
		return sums, err
	}
	defer file.Close()

	var city string
	var temp float64

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		el := strings.Split(line[:len(line)-1], ";")
		city = el[0]
		if len(el) != 2 {
			log.Println(line[:len(line)-1])
			log.Println("invalid line")
			continue
		}
		temp, err = strconv.ParseFloat(el[1], 64) // definitely can be optimized since we know the format and range is always -99.9 - 99.9
		if err != nil {
			log.Println(line[:len(line)-1])
			log.Println(err)
			os.Exit(1)
		}

		if v, ok := sums[city]; !ok {
			sums[city] = location{
				mean: sums[city].mean + temp,
				min:  temp,
				max:  temp,
				n:    sums[city].n + 1,
			}
		} else {
			v.n++
			v.mean += temp
			v.min = min(v.min, temp)
			v.max = max(v.max, temp)
			sums[city] = v
		}
	}

	for c, v := range sums {
		v.mean = v.mean / float64(v.n)
		sums[c] = v
	}

	return sums, nil
}

func TestCalc0(t *testing.T) {
	if _, err := os.Stat("10mil_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(10000, 10000000, "10mil")
	}

	results, err := calc0("10mil_c10000_l10000000.csv")
	require.NoError(t, err, "calc0")
	require.Len(t, results, 10000)
}

// calc1 is an attempt to optimize calc0. Still single-core.
// 1/2 of the time of calc0
func calc1(filename string) (results map[string]location, err error) {

	// warming up the map
	initFloatMap()

	sums := make(map[string]loc, 10000)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var line []byte
	var city string
	var temp int

	reader := bufio.NewReader(file)
	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}
		city = string(line[:slices.Index(line, ';')])
		temp = floatMap[string(line[slices.Index(line, ';')+1:])]

		if v, ok := sums[city]; !ok {
			sums[city] = loc{
				mean: sums[city].mean + temp,
				min:  temp,
				max:  temp,
				n:    sums[city].n + 1,
			}
		} else {
			v.n++
			v.mean += temp
			v.min = min(v.min, temp)
			v.max = max(v.max, temp)
			sums[city] = v
		}
	}

	// 0.001-0.003 seconds for 10mil
	results = make(map[string]location, len(sums))
	for c, v := range sums {
		results[c] = location{
			mean: float64(v.mean) / float64(v.n),
			min:  float64(v.min) / 10,
			max:  float64(v.max) / 10,
			n:    v.n,
		}
	}

	return results, nil
}

func TestCalc1(t *testing.T) {
	if _, err := os.Stat("10mil_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(10000, 10000000, "10mil")
	}

	results, err := calc1("10mil_c10000_l10000000.csv")
	require.NoError(t, err, "calc1")
	require.Len(t, results, 10000)
}
