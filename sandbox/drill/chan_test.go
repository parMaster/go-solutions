// interview questions
package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// single-thread solution
func JoinChannels(chs ...<-chan int) <-chan int {
	resCh := make(chan int)
	go func() {
		defer close(resCh)
		stop := false
		for !stop {
			stop = true
			for _, ch := range chs {
				v, ok := <-ch
				if ok {
					resCh <- v
					stop = false
				}
			}
		}
	}()
	return resCh
}

func Test_JoinChannels(t *testing.T) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		for _, v := range []int{1, 2, 3} {
			a <- v
		}
		close(a)
	}()

	go func() {
		for _, v := range []int{10, 20, 30} {
			b <- v
		}
		close(b)
	}()

	go func() {
		for _, v := range []int{100, 200, 300} {
			c <- v
		}
		close(c)
	}()

	for v := range JoinChannels(a, b, c) {
		fmt.Println(v)
	}
}

// multithread classic fan-in
func JoinChannelsFanIn(chs ...<-chan int) <-chan int {
	resCh := make(chan int)

	wg := &sync.WaitGroup{}

	multiplex := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			resCh <- v
		}
	}

	for _, ch := range chs {
		wg.Add(1)
		go multiplex(ch)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	return resCh
}

func Test_JoinChannelsFanIn(t *testing.T) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		for _, v := range []int{1, 2, 3} {
			a <- v
		}
		close(a)
	}()

	go func() {
		for _, v := range []int{10, 20, 30} {
			b <- v
		}
		close(b)
	}()

	go func() {
		for _, v := range []int{100, 200, 300} {
			c <- v
		}
		close(c)
	}()

	received := 0
	for v := range JoinChannelsFanIn(a, b, c) {
		received++
		fmt.Println(v)
	}

	assert.Equal(t, 9, received)
}

// from example (uglier fan-in impl.)
func JoinChannelsExample(chs ...<-chan int) <-chan int {
	result := make(chan int)

	wg := &sync.WaitGroup{}

	wg.Add(len(chs))

	for _, ch := range chs {
		ch := ch

		go func() {
			defer wg.Done()

			for v := range ch {
				result <- v
			}
		}()
	}

	go func() {
		wg.Wait()

		close(result)
	}()

	return result
}

func Test_JoinChannelsExample(t *testing.T) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		for _, v := range []int{1, 2, 3} {
			a <- v
		}
		close(a)
	}()

	go func() {
		for _, v := range []int{10, 20, 30} {
			b <- v
		}
		close(b)
	}()

	go func() {
		for _, v := range []int{100, 200, 300} {
			c <- v
		}
		close(c)
	}()

	received := 0
	for v := range JoinChannelsExample(a, b, c) {
		received++
		fmt.Println(v)
	}

	assert.Equal(t, 9, received)
}

// tricky question

func worker() chan int {
	ch := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
		close(ch)
	}()
	return ch
}

func Test_WorkerWait(t *testing.T) {

	timeStart := time.Now()
	_, _ = <-worker(), <-worker()
	println(int(time.Since(timeStart).Seconds())) // 6 sec, not 3

	// because it's equivalen to:
	timeStart = time.Now()
	<-worker()
	<-worker()
	println(int(time.Since(timeStart).Seconds())) // 6 sec

	// utilizing fan-in to execute workers concurrently:
	timeStart = time.Now()
	for v := range JoinChannelsExample(worker(), worker(), worker()) {
		println(v)
	}
	println(int(time.Since(timeStart).Seconds())) // 3 sec
}
