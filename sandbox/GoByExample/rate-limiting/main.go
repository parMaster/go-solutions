package main

import (
	"fmt"
	"time"
)

func main() {

	// rate limiting
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		fmt.Println("i: ", i, time.Now())
		requests <- i
	}
	close(requests)

	limiter := time.Tick(500 * time.Millisecond)

	for req := range requests {
		<-limiter
		fmt.Println("getting request from the queue", req, time.Now())
	}

	// bursty limiter
	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("reading burstyRequests", req, time.Now())
	}
}
