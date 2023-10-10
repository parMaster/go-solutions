package main

// https://gobyexample.com/atomic-counters

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	var ops atomic.Uint64

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		// 50 goroutines will be run
		go func() {

			// each of which will run the counter 1000 times and increment the ops by 1
			for c := 0; c < 1000; c++ {
				ops.Add(1)
			}

			// each waitgroup will finish before goroutine quit
			wg.Done()
		}()
	}

	// it will wait for all 50 WGs to finish
	wg.Wait()

	// never knew there's Load() method for atomic
	// still, it should print 50000 here
	fmt.Println("ops:", ops.Load())

	//
	//

	// Can I reproduce the same with the buffered channel?
	// just to remember workers chapter
	ops.Store(0)

	n := 50
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)

	for i := range ch {
		fmt.Println("i=", i)

		for c := 0; c < 1000; c++ {
			ops.Add(1)
		}
	}

	// yes I can!

	fmt.Println("ops:", ops.Load())
}
