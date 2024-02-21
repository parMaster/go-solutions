package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

// genWrite generates a file with random data
func genWrite(locations, lines int, filename string) error {
	file, err := os.OpenFile(fmt.Sprintf("%s_c%d_l%d.csv", filename, locations, lines), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	for range lines {
		_, err := fmt.Fprintf(file, "city%d;%.1f\n", rand.Intn(locations), rand.Float64()*199.8-99.9)
		if err != nil {
			return err
		}
	}
	return file.Close()
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
		i += len(line)
	}
}

type location struct {
	mean float64
	min  float64
	max  float64
	n    int
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

// Attempt to optimize calc0

var floatMap map[string]int

// generates a map of float strings to int
func initFloatMap() {
	floatMap = make(map[string]int, 1999)
	for i := -999; i < 1000; i++ {
		floatMap[fmt.Sprintf("%.1f", float64(i)/10)] = i
	}
	floatMap["-0.0"] = 0
}

// calc1 is an attempt to optimize calc0. Still single-core.
// 1/2 of the time of calc0
func calc1(filename string) (results map[string]location, err error) {

	// warming up the map
	initFloatMap()

	type loc struct {
		mean int
		min  int
		max  int
		n    int
	}

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

func formatResults(r map[string]location) string {
	var sb strings.Builder

	cities := maps.Keys(r)
	slices.Sort(cities)

	for _, c := range cities {
		sb.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f,", c, r[c].min, r[c].mean, r[c].max))
	}
	str := sb.String()
	return "{" + str[:len(str)-1] + "}"
}

func main() {

	// generate test data
	// sudo mkdir /mnt/ramdisk
	// sudo mount -t tmpfs -o rw,size=2G tmpfs /mnt/ramdisk
	// genWrite(10000, 100000000, "/mnt/ramdisk/100mil") // 100 mil rows ~1.3GB on ramdisk, keeping SSD happy

	start := time.Now()
	var i int
	for range 10 {
		i += readLines("10mil_c10000_l10000000.csv")
	}
	fmt.Printf("100 mil lines read in %.1f seconds, i= %d\n", time.Since(start).Seconds(), i)

	var results map[string]location
	var err error

	start = time.Now()
	results, err = calc0("10mil_c10000_l10000000.csv")
	if err != nil {
		log.Println(err)
	}
	if len(results) != 10000 {
		log.Println("invalid results")
	}
	fmt.Printf("Calc0 done in %.1f seconds (would be %.1f seconds for 1B)\n", time.Since(start).Seconds(), time.Since(start).Seconds()*100)

	start = time.Now()
	results, err = calc1("10mil_c10000_l10000000.csv")
	if err != nil {
		log.Println(err)
	}
	if len(results) != 10000 {
		log.Println("invalid results")
	}
	fmt.Printf("Calc1 done in %.1f seconds (would be %.1f seconds for 1B)\n", time.Since(start).Seconds(), time.Since(start).Seconds()*100)

	start = time.Now()
	formatted := formatResults(results)
	fmt.Printf("Formatted in %.3f seconds, %d cities\n", time.Since(start).Seconds(), strings.Count(formatted, ",")+1)
	// fmt.Println(formatted)
}
