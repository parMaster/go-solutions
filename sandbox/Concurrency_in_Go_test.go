package sandbox

/*
	Testing concurrency patterns and practices while reading
	"Concurrency in Go. Tools and Techniques for Developers" by Katherine Cox-Buday
*/

import (
	"fmt"
	"log"
	"math"
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

// CONCURRENCY PRIMITIVES

// sync package
//

func Test_SimulateWaitGroups(t *testing.T) {

	var wg, cnt atomic.Int32

	for i := range 10 {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			log.Println("Task", i)
			cnt.Add(1)
		}()
	}

	for wg.Load() > 0 {
		time.Sleep(100 * time.Millisecond)
	}

	assert.Equal(t, int32(0), wg.Load())
	assert.Equal(t, int32(10), cnt.Load())
	log.Println("Done")
}

func Test_WorkGroups(t *testing.T) {

	var wg sync.WaitGroup
	var cnt atomic.Int32

	numTasks := 10
	wg.Add(numTasks)
	for i := range numTasks {
		go func() {
			defer wg.Done()
			log.Println("Task", i)
			cnt.Add(1)
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(10), cnt.Load())
	log.Println("Done")
}

func Test_Cond(t *testing.T) {

	condition := false
	conditionTrue := func() bool {
		return condition
	}

	c := sync.NewCond(&sync.Mutex{})

	go func() {
		time.Sleep(1 * time.Second)
		condition = true
		c.Signal()
	}()

	c.L.Lock()
	for conditionTrue() == false {
		c.Wait()
	}
	c.L.Unlock()

	assert.True(t, condition)
	// I don't really understand this one
}

func Test_CondQueue(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}
	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue", i)
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

// sync.Once
func Test_Once(t *testing.T) {

	cnt := 0
	inc := func() {
		cnt++
	}

	var once sync.Once

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			once.Do(inc)
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, cnt)
}

// once counts how many times Do() was called
func Test_OnceDo(t *testing.T) {

	once := sync.Once{}
	//
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	once.Do(increment)
	once.Do(decrement)

	// once counts how many times Do() was called, not how many times the function was executed
	assert.Equal(t, 1, count)
}

// sync.Pool
func Test_Pool(t *testing.T) {

	myPool := &sync.Pool{
		New: func() interface{} {
			log.Printf("Creating new instance")
			return struct{}{}
		},
	}
	myPool.Get()             // Creating new instance
	instance := myPool.Get() // Creating new instance
	myPool.Put(instance)     // Releasing instance
	myPool.Get()             // Reusing instance, no new instance created
}

//
// Channels
//

func Test_CloseChan(t *testing.T) {

	dataStream := make(chan string)
	go func() {
		dataStream <- "hello from func"
	}()

	message, ok := <-dataStream
	assert.True(t, ok) // value received
	assert.Equal(t, "hello from func", message)

	go func() {
		log.Println("Closing channel")
		close(dataStream)
	}()

	log.Println("Blocking on read from open channel")
	message, ok = <-dataStream
	log.Println("Channel closed, block released")
	assert.False(t, ok)          // channel closed
	assert.Equal(t, "", message) // no message received

	// reading from closed channel
	// this way multiple downstreams can sit on blocking reading call
	// until a single upstream doesn't close the channel
	log.Println("Reading from Closed channel")
	message, ok = <-dataStream
	log.Println("Closed channel read is not blocking")
	assert.False(t, ok)          // channel closed
	assert.Equal(t, "", message) // no message received

	// closing of closed channel will panic
	// close(dataStream)

	// check before close:
	if _, ok := <-dataStream; ok {
		close(dataStream)
	}
}

func Test_RangingOverChannel(t *testing.T) {

	intStream := make(chan int)

	go func() {
		defer close(intStream)
		for i := range 10 {
			intStream <- i
		}
	}()

	sum := 0
	for received := range intStream {
		log.Println("Received", received)
		sum += received
	}
	assert.Equal(t, 45, sum)
}

// unlocking multiple goroutines
func Test_UnlockingMultiple(t *testing.T) {
	lockStream := make(chan interface{})

	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-lockStream
			log.Println("Worker", i, "unblocked")
		}(i)
	}

	close(lockStream)
	wg.Wait()
}

// Buffered Channel
func Test_BufferedChannel(t *testing.T) {

	intStream := make(chan int, 4)

	go func(n int) {
		defer close(intStream)
		defer log.Println("Producer Done.")
		for i := range n {
			log.Println("Sending", i)
			intStream <- i
		}

	}(10)

	for i := range intStream {
		log.Println("Receiving", i, "|", len(intStream), "in queue")

	}

}

// Channel ownership
// owner responsible for writing to channel and closing channel
// returns read-only channel
// this ensures that:
// - writes won't happen on closed or nil channel
// - close will happen once
func Test_Channel_Owner(t *testing.T) {

	channelOwner := func() <-chan int {
		intStream := make(chan int, 5)

		go func() {
			defer close(intStream)

			for i := range 10 {
				intStream <- i
			}
		}()
		return intStream
	}

	received := 0
	ints := channelOwner()
	for i := range ints {
		log.Println("Received", i)
		received++
	}
	assert.Equal(t, 10, received) // everything's received
}

// Select statement
func Test_Select(t *testing.T) {

	start := time.Now()
	block := make(chan any)
	go func() {
		time.Sleep(1 * time.Second)
		close(block)
	}()

	select { // equal to simple blocking read `<-block`
	case <-block:
		log.Println("Unblocked in", time.Since(start))
	}
	log.Println("Done.")

	//
	// Both channels are ready at the same time
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)
	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)

	// roughly equal chance of being selected, difference is less than 10%
	assert.Less(t, math.Abs(float64(c2Count-c1Count)), float64(100))

	//
	// Timeout
	timedOut := false
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		timedOut = true
		fmt.Println("Timed out.")
	}
	assert.True(t, timedOut)

	//
	// Default clause
	defaulted := false
	var c11 <-chan int
	var c12 <-chan int
	select {
	case <-c11:
	case <-c12:
	default:
		defaulted = true
		log.Println("Defaulted")
	}
	assert.True(t, defaulted)

	//
	// Loop with default
	// run default clause while waiting for the channel to unblock

	done := make(chan any)
	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()

	workCounter := 0
	keepGoing := true
	for keepGoing {
		select {
		case <-done:
			keepGoing = false
		default:
		}

		workCounter++
		time.Sleep(1 * time.Millisecond)
	}

	log.Println("Made", workCounter, "cycles before stopped")
	assert.Less(t, 1, workCounter) // workCounter > 1

	//
	// Block forever
	// select {}
}

// CONCURRENCY PATTERNS

// to be continued ..
