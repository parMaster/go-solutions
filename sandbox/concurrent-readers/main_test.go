package main

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// iterative solution
// ReadCount reads file and returns the number of char characters
func ReadCount(file string, char byte) int {

	f, err := os.OpenFile(file, os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}

	cnt := 0
	b := make([]byte, 1024)
	for {
		rb, err := f.Read(b)
		if err == io.EOF {
			break
		}

		for i := range rb {
			if b[i] == char {
				// fmt.Print(string(b[i]))
				cnt++
			}
		}
	}

	return cnt
}

func ReadCountConcurrent(file string, char byte) int {
	f, err := os.OpenFile(file, os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}

	limit := runtime.NumCPU()
	limiter := make(chan struct{}, limit)

	results := make(chan int)
	sum := 0
	wg := sync.WaitGroup{}
	go func() {
		for res := range results {
			sum += res
			wg.Done()
		}
	}()

	for {
		b := make([]byte, 1024)
		rb, err := f.Read(b)
		if err == io.EOF {
			break
		}

		wg.Add(1)
		go func() {
			limiter <- struct{}{}
			cnt := 0
			for i := range rb {
				if b[i] == char {
					cnt++
				}
			}
			results <- cnt
			<-limiter
		}()
	}

	wg.Wait()
	return sum
}

func Test_ReadCount(t *testing.T) {
	cnt := ReadCount("smc.txt", 's')
	assert.Equal(t, 3110, cnt)
}

func Test_ReadCountConcurrent(t *testing.T) {
	cnt := ReadCountConcurrent("smc.txt", 's')
	assert.Equal(t, 3110, cnt)
}

func BenchmarkCount(b *testing.B) {
	ReadCount("smc.txt", 's')
}

func BenchmarkCountConcurrent(b *testing.B) {
	ReadCountConcurrent("smc.txt", 's')
}

func ReadCountConcurrentJobs(file string, char byte) int {
	f, err := os.OpenFile(file, os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}
	sum := 0

	var counter = func(jobs <-chan []byte) chan int {
		results := make(chan int)

		go func() {
			for {
				if job, more := <-jobs; more {
					cnt := 0
					for i := range job {
						if job[i] == char {
							cnt++
						}
					}
					results <- cnt
				} else {
					close(results)
					break
				}
			}
		}()

		return results
	}

	jobs := make(chan []byte)

	go func(f io.ReadCloser) {
		for {
			buf := make([]byte, 1024)
			rb, err := f.Read(buf)
			if err == io.EOF {
				close(jobs)
				err := f.Close()
				if err != nil {
					log.Printf("error closing reader: %v\n", err)
				}
				break
			}
			jobs <- buf[:rb]
		}
	}(f)

	resCh := counter(jobs)

	for res := range resCh {
		sum += res
	}

	return sum
}

func Test_ReadCountConcurrentJobs(t *testing.T) {
	cnt := ReadCountConcurrentJobs("smc.txt", 's')
	assert.Equal(t, 3110, cnt)
}
