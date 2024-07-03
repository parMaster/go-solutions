package main

import (
	"log"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func workerFn(w int, jobCh <-chan int, resCh chan<- int) {
	for job := range jobCh {
		log.Printf("worker %d working on %d", w, job)
		resCh <- job * job
	}
}

// primitive worker pool
func runPool(jobNum int) int {
	poolSize := runtime.NumCPU()

	jobCh := make(chan int, jobNum)
	resCh := make(chan int, jobNum)

	for w := range poolSize {
		go workerFn(w, jobCh, resCh)
	}

	for i := range jobNum {
		jobCh <- i
	}
	close(jobCh)

	doneCh := make(chan struct{})
	results := atomic.Int32{}
	go func() {
		defer close(doneCh)
		for range jobNum {
			log.Println("result returned:", <-resCh)
			results.Add(1)
		}
	}()
	<-doneCh

	log.Println("job done")

	return int(results.Load())
}

func Test_WP(t *testing.T) {
	assert.Equal(t, 10, runPool(10))
}

// once again to remember, shorter version
func runJobs(num int) int {

	workersNum := runtime.NumCPU()

	jobsCh := make(chan int, num)
	resCh := make(chan int, num)

	for w := range workersNum {
		go func(wId int, jobsCh <-chan int, resCh chan<- int) {
			for job := range jobsCh {
				log.Printf("worker %d, job %d", wId, job)
				time.Sleep(10 * time.Millisecond)
				resCh <- job * job
			}
		}(w, jobsCh, resCh)
	}

	for job := range num {
		jobsCh <- job
	}
	close(jobsCh)

	received := 0
	for range num {
		received++
		log.Println("received result", <-resCh)
	}

	log.Println("job done")
	return received
}

func Test_AnotherWP(t *testing.T) {
	assert.Equal(t, 10, runJobs(10))
	assert.Equal(t, 1000, runJobs(1000))
}
