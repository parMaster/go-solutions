package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/gammazero/workerpool"

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

// preliminary int results
type loc struct {
	mean int
	min  int
	max  int
	n    int
}

// actual float results
type location struct {
	mean float64
	min  float64
	max  float64
	n    int
}

// dictionary for float strings to int
var floatMap map[string]int

// generates a map of float strings to int
func initFloatMap() {
	floatMap = make(map[string]int, 1999)
	for i := -999; i < 1000; i++ {
		floatMap[fmt.Sprintf("%.1f", float64(i)/10)] = i
	}
	floatMap["-0.0"] = 0
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

// Multi-core solution
func solve(filename string, size int) (results map[string]location, err error) {
	initFloatMap()

	cores := runtime.NumCPU()
	wp := workerpool.New(cores)

	// res is a channel for results from workers
	res := make(chan map[string]loc, cores)
	defer close(res)
	// start workers
	for i := range cores {
		start, stop := i*size/cores, (i+1)*size/cores
		wp.Submit(func() { calcM(filename, start, stop, res) })
	}

	// receiving results
	sums := make(map[string]loc, 10000)
	for a := 1; a <= cores; a++ {

		// merging results
		for city, data := range <-res {
			if v, ok := sums[city]; !ok {
				sums[city] = data
			} else {
				v.n += data.n
				v.mean += v.mean
				v.min = min(v.min, data.min)
				v.max = max(v.max, data.max)
				sums[city] = v
			}
		}
	}

	wp.StopWait()

	// convert to float
	total := 0
	results = make(map[string]location, len(sums))
	for c, v := range sums {
		results[c] = location{
			mean: float64(v.mean) / float64(v.n),
			min:  float64(v.min) / 10,
			max:  float64(v.max) / 10,
			n:    v.n,
		}
		total += v.n
	}

	if total != size {
		log.Println("invalid total", total, size)
	}

	return results, nil
}

func calcM(filename string, start, stop int, results chan<- map[string]loc) error {

	fmt.Println("calcM for ", start, stop)
	sums := make(map[string]loc, 10000)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var line []byte
	var city string
	var temp int
	i := 0

	reader := bufio.NewReader(file)
	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}
		if i < start {
			i++
			continue
		}
		if i >= stop {
			break
		}
		i++

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

	// fmt.Println("done with chunk")
	results <- sums

	return nil
}

func main() {
	// generate test data
	// sudo mkdir /mnt/ramdisk
	// sudo mount -t tmpfs -o rw,size=2G tmpfs /mnt/ramdisk
	// genWrite(10000, 100000000, "/mnt/ramdisk/100mil") // 100 mil rows ~1.3GB on ramdisk, keeping SSD happy

	// testing on 10 millin lines dataset
	cities := 10000
	lines := 10000000
	// generatin 10 mil rows
	if _, err := os.Stat("10mil_c10000_l10000000.csv"); os.IsNotExist(err) {
		genWrite(cities, lines, "10mil")
	}

	var start time.Time
	var results map[string]location
	var err error

	// Multi-core solution
	start = time.Now()
	results, err = solve("10mil_c10000_l10000000.csv", lines)
	if err != nil {
		log.Println(err)
	}
	if len(results) != cities {
		log.Printf("invalid results, expected %d, got %d\n", cities, len(results))
	}
	formatted := formatResults(results)
	fmt.Printf("Done in %.1f seconds (would be %.1f seconds for 1B)\n", time.Since(start).Seconds(), time.Since(start).Seconds()*100)

	os.WriteFile("results.txt", []byte(formatted), 0644)
}
