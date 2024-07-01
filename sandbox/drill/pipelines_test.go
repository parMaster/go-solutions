// reminder on pipelines with done channel. Looks like I got it

package main

import (
	"log"
	"testing"
)

func sliceToChan(done <-chan struct{}, s []int) <-chan int {
	intCh := make(chan int)

	go func() {
		defer close(intCh)

		for _, v := range s {
			select {
			case <-done:
				return
			case intCh <- v:
				continue
			}
		}
	}()

	return intCh
}

func double(done <-chan struct{}, inputCh <-chan int) <-chan int {
	outCh := make(chan int)

	go func() {
		defer close(outCh)
		for v := range inputCh {
			select {
			case <-done:
				return
			case outCh <- (v * 2):
				continue
			}
		}
	}()

	return outCh
}

func take(done <-chan struct{}, inputCh <-chan int, n int) <-chan int {
	outCh := make(chan int)

	go func() {
		defer close(outCh)
		i := 0
		for i < n {
			select {
			case <-done:
				return
			case outCh <- <-inputCh:
				i++
				continue
			}
		}
	}()

	return outCh
}

func Test_pipe(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

	done := make(chan struct{})
	defer close(done)

	outCh := double(done, sliceToChan(done, ints))

	for range len(ints) {
		select {
		case v := <-outCh:
			log.Printf("received: %d\n", v)
		case <-done:
			return
		}
	}

}

func Test_Take(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	done := make(chan struct{})

	outCh := take(done, double(done, sliceToChan(done, ints)), 5)

	for {
		select {
		case v, ok := <-outCh:
			if !ok {
				return
			}
			log.Printf("received: %d\n", v)
		case <-done:
			return
		}
	}
}
