package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func main() {

	// generate test data

	// sudo mkdir /mnt/ramdisk
	// sudo mount -t tmpfs -o rw,size=2G tmpfs /mnt/ramdisk
	// genWrite(10000, 100000000, "/mnt/ramdisk/mil")

	genWrite(10000, 1000000, "mil")

	start := time.Now()

	results, err := calc0("mil_c10000_l1000000.csv")
	if err != nil {
		log.Println(err)
	}
	if len(results) != 10000 {
		log.Println("invalid results")
	}

	end := time.Now()
	fmt.Printf("Done in %.1f seconds", end.Sub(start).Seconds())

}
