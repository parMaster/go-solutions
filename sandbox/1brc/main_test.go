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
	if _, err := os.Stat("input_c10_l10.csv"); os.IsNotExist(err) {
		require.NoError(t, genWrite(10, 10, "input"))
	}
	if _, err := os.Stat("input_c10000_l1000000.csv"); os.IsNotExist(err) {
		require.NoError(t, genWrite(10000, 1000000, "input"))
	}
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
	if _, err := os.Stat("input_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(10000, 10000000, "input")
	}

	n := readLines("input_c10000_l10000000.csv")
	require.Equal(t, 10000000, n)
	// 0.31s for 10 mils
}

// dictionary for float strings to int
var floatMap map[string]int

// generates a map of float strings to int
func initFloatMap() {
	floatMap = make(map[string]int, 2000)
	for i := -999; i < 1000; i++ {
		floatMap[fmt.Sprintf("%.1f\n", float64(i)/10)] = i
	}
	floatMap["-0.0\n"] = 0
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
	for range 100000 {
		n := rand.Float64()*199.8 - 99.9
		n = math.Round(n*10) / 10
		require.Equal(t, int(n*10), floatStrToInt(fmt.Sprintf("%.1f\n", n)))
	}
}

func parseFloat(b []byte) int {

	neg := false
	n := 0
	b = b[:len(b)-1]
	if b[0] == '-' {
		neg = true
		b = b[1:]
	}
	if len(b) == 4 {
		// 12.3
		n = int(b[0]-0x30)*100 + int(b[1]-0x30)*10 + int(b[3]-0x30)
	} else if len(b) == 3 {
		// 1.2
		n = int(b[0]-0x30)*10 + int(b[2]-0x30)
	}

	if neg {
		return -n
	}
	return n
}

func TestParseFloatNew(t *testing.T) {

	require.Equal(t, 12, parseFloat([]byte("1.2\n")))
	require.Equal(t, 123, parseFloat([]byte("12.3\n")))

	require.Equal(t, -12, parseFloat([]byte("-1.2\n")))
	require.Equal(t, -123, parseFloat([]byte("-12.3\n")))

	for range 100000 {
		n := rand.Float64()*199.8 - 99.9
		n = math.Round(n*10) / 10
		require.Equal(t, int(n*10), parseFloat([]byte(fmt.Sprintf("%.1f\n", n))))
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
	if _, err := os.Stat("input_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(10000, 10000000, "input")
	}

	results, err := calc0("input_c10000_l10000000.csv")
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
				mean: sums[city].mean + int64(temp),
				min:  temp,
				max:  temp,
				n:    sums[city].n + 1,
			}
		} else {
			v.n++
			v.mean += int64(temp)
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
	if _, err := os.Stat("input_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(10000, 10000000, "input")
	}

	results, err := calc1("input_c10000_l10000000.csv")
	require.NoError(t, err, "calc1")
	require.Len(t, results, 10000)
}

// group by city and run the calculation in parallel??
func readCities(filename string) (cities map[string][]int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cities = make(map[string][]int, 10000)

	i := 0
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		city := string(line[:slices.Index(line, ';')])
		temp := floatMap[string(line[slices.Index(line, ';')+1:])]

		if _, ok := cities[city]; !ok {
			cities[city] = []int{temp}
		} else {
			cities[city] = append(cities[city], temp)
		}

		i++
	}

	if i != 10000000 {
		panic("invalid number of lines")
	}

	return cities, nil
}

func TestReadCities(t *testing.T) {

	initFloatMap()

	cities, err := readCities("input_c10000_l10000000.csv")
	// no good! appends killing the performance
	require.NoError(t, err)
	require.Len(t, cities, 10000)
}

// read the file in chunks the proper way so it's possible to calculate the results in parallel
func readSection(file string, from, to int64) (lines []string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	_, err = f.Seek(from, 0)
	if err != nil {
		return nil, err
	}

	i := from
	for i < to {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return lines, err
		}
		i += int64(len(line))
		lines = append(lines, string(line))
	}

	return lines, nil
}

func TestReadSection(t *testing.T) {

	file := "input_c10000_l1000000.csv"

	// testing division into chunks, make sure we get all the lines independently of the number of chunks
	for i := 1; i < 16; i++ {
		cores := i
		intervals, err := getIntervals(file, cores)
		require.NoError(t, err)

		fmt.Println(intervals)

		lines := []string{}
		for _, m := range intervals {
			ll, err := readSection(file, m.from, m.to)
			lines = append(lines, ll...)
			require.NoError(t, err)
		}
		require.Len(t, lines, 1000000)
	}
}
