package sandbox

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Limit the number of goroutines that can run at the same time
func Test_CountingSemaphore(t *testing.T) {
	maxGoroutines := 5
	semaphore := make(chan struct{}, maxGoroutines)

	start := time.Now()

	var cnt atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Simulate a task
			fmt.Printf("Running task %d\n", i)
			time.Sleep(1 * time.Second)
			cnt.Add(1)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, int32(20), cnt.Load())
	assert.Equal(t, 4, int(time.Since(start).Seconds())) // 20/5 = 4
}

func Test_Channels(t *testing.T) {

	wg := sync.WaitGroup{}
	ch := make(chan int)

	for i := range 4 {
		wg.Add(1)
		go func(i int) {
			time.Sleep(1 * time.Second)
			defer wg.Done()
			ch <- i
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	respCnt := 0
	for response := range ch {
		fmt.Println(response)
		respCnt++
	}

	assert.Equal(t, 4, respCnt)
}
