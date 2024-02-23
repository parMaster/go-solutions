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
	"strconv"
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
		floatMap[fmt.Sprintf("%.1f\n", float64(i)/10)] = i
	}
	floatMap["-0.0\n"] = 0
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

// advance position in the file past the next newline
func advancePos(pos int64, f *os.File) int64 {
	if pos == 0 {
		return 0
	}
	f.Seek(pos, 0)
	r := bufio.NewReader(f)
	advance, _ := r.ReadBytes('\n')
	pos += int64(len(advance))
	return pos
}

// anvance interval margins
func advance(from, to int64, f *os.File) (int64, int64) {
	from = advancePos(from, f)
	to = advancePos(to, f)
	return from, to
}

type interval struct {
	from int64
	to   int64
}

// returns intervals to read from file, considering newlines
func getIntervals(file string, cores int) (margins []interval, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, _ := f.Stat()
	size := info.Size()
	chunkSize := size / int64(cores)

	for i := 0; i < cores; i++ {
		from := chunkSize * int64(i)
		to := chunkSize * int64(i+1)
		from, to = advance(from, to, f)
		margins = append(margins, interval{from, to})
	}

	return margins, nil
}

// Multi-core solution
func solve(filename string, size, cores int) (results map[string]location, err error) {
	initFloatMap()

	intervals, _ := getIntervals(filename, cores)

	wp := workerpool.New(cores)

	// res is a channel for results from workers
	res := make(chan map[string]loc, cores)
	defer close(res)
	// start workers
	for _, interval := range intervals {
		wp.Submit(func() { calcM(filename, interval.from, interval.to, res) })
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

func calcM(filename string, start, stop int64, results chan<- map[string]loc) error {

	fmt.Println("calcM for interval ", start, stop)
	sums := make(map[string]loc, 10000)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var line []byte
	var city string
	var temp int

	r := bufio.NewReader(file)
	file.Seek(start, 0)
	i := start

	for i < stop {
		line, err = r.ReadBytes('\n')
		i += int64(len(line))
		if err == io.EOF {
			break
		}

		semicol := slices.Index(line, ';')
		city = string(line[:semicol])
		temp = floatMap[string(line[semicol+1:])]

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

	results <- sums

	return nil
}

func main() {

	cores := runtime.NumCPU()
	args := os.Args[1:]
	if slices.Index(args, "--cores") != -1 {
		if c, err := strconv.Atoi(args[slices.Index(args, "--cores")+1]); err == nil {
			cores = c
		}
	}
	cities := 10000
	if slices.Index(args, "--cities") != -1 {
		if c, err := strconv.Atoi(args[slices.Index(args, "--cities")+1]); err == nil {
			cities = c
		}
	}
	lines := 10000000
	if slices.Index(args, "--lines") != -1 {
		if l, err := strconv.Atoi(args[slices.Index(args, "--lines")+1]); err == nil {
			lines = l
		}
	}

	fmt.Println("Runnig", cores, "workers on", runtime.NumCPU(), "CPU cores")
	fmt.Println("cities:", cities, "lines:", lines, "(", fmt.Sprintf("%.0f %c", float64(lines)/float64(1000000), 'M'), ")")

	filename := fmt.Sprintf("input_c%d_l%d.csv", cities, lines)

	// generate test data
	// sudo mkdir /mnt/ramdisk
	// sudo mount -t tmpfs -o rw,size=2G tmpfs /mnt/ramdisk
	// genWrite(10000, 100000000, "/mnt/ramdisk/100mil") // 100 mil rows ~1.3GB on ramdisk, keeping SSD happy

	// generatin 'lines' random lines of data
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		genWrite(cities, lines, "input")
	}

	var start time.Time
	var results map[string]location
	var err error

	// Multi-core solution
	start = time.Now()
	results, err = solve(filename, lines, cores)
	if err != nil {
		log.Println(err)
	}
	if len(results) != cities {
		log.Printf("invalid results, expected %d, got %d\n", cities, len(results))
	}
	formatted := formatResults(results)
	fmt.Printf("Done in %.1f seconds (would be %.1f seconds for 1B)\n", time.Since(start).Seconds(), time.Since(start).Seconds()*(1000000000/float64(lines)))

	if slices.Index(args, "--ro") == -1 {
		os.WriteFile("results.txt", []byte(formatted), 0644)
		fmt.Println("results written to results.txt")
	}
}
