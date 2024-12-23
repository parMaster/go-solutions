package sandbox

/*
	Testing concurrency patterns and practices while reading
	"Concurrency in Go. Tools and Techniques for Developers" by Katherine Cox-Buday
*/

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
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

/*
When you close a channel in Go, you're signaling that no more values will be sent on that channel.
However, any values that have already been sent on the channel before it was closed are still
available to be received. This is true for both buffered and unbuffered channels.
*/
func Test_Send_Close_Read(t *testing.T) {
	stream := make(chan any)
	go func() {
		stream <- "value"
		close(stream)
	}()

	for {
		val, ok := <-stream
		if !ok {
			break
		}
		fmt.Println(val)
	}

}

// CONCURRENCY PRIMITIVES

// sync package
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

func Test_WaitGroups(t *testing.T) {

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

// I don't really understand this one
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
}

func Test_CondQueue(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]any, 0, 10)
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

// sync.Once - .Do() will be executed not more than once
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
// reusing the instances from the pool instead of instantiating the new ones
func Test_Pool(t *testing.T) {

	myPool := &sync.Pool{
		New: func() any {
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
	}() // just to illustrate next lines

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
	assert.False(t, ok)          // channel closed, ok == false
	assert.Equal(t, "", message) // no message received

	if ok { // closing of closed channel will panic
		close(dataStream)
	}

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
	lockStream := make(chan any)

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
	c1 := make(chan any)
	close(c1)
	c2 := make(chan any)
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

func Test_Confinement(t *testing.T) {

	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}
	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	results := chanOwner()
	consumer(results)
}

// Sending iteration variables out on a channel
func Test_For_Select_Loop(t *testing.T) {

	done := make(chan any)
	stringStream := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		done <- "timeout"
	}()

	go func() {
		for {
			s := <-stringStream
			log.Printf("Received %s", s)
			time.Sleep(520 * time.Millisecond)
		}
	}()

	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done: // will read in 1 second no matter what? like a timeout?
			return
		case stringStream <- s: // blocking write, waiting to stringStream to be empty
		}
	}
}

// Looping infinitely waiting to be stopped
func Test_Loop_Infinitely(t *testing.T) {
	done := make(chan any)

	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()

	i := 0
	for {
		select {
		case <-done:
			log.Println("Done")
			assert.Equal(t, 5, i) // 5 times 200ms in 1 second
			return
		default:
			// i++
			// log.Println("Working")
			// time.Sleep(200 * time.Millisecond)
		}
		// or here:
		i++
		log.Println("Working")
		time.Sleep(200 * time.Millisecond)
	}
}

// Goroutine leak demo
func Test_NoTermination(t *testing.T) {
	worked := false
	doWork := func(strings <-chan string) <-chan any {
		completed := make(chan any)
		go func() { // this goroutine will leak, never stop
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings { // because this read will be blocked forever
				fmt.Println(s)
				worked = true
			}
		}()
		return completed
	}
	doWork(nil)

	assert.False(t, worked)
}

// termination example from the book
func Test_Termination_Example(t *testing.T) {
	doWork := func(
		done <-chan any, strings <-chan string,
	) <-chan any {
		terminated := make(chan any)
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan any)
	terminated := doWork(done, nil)
	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	<-terminated
	fmt.Println("Done.")
}

// termination example with something really being sent on strings chan
func Test_ProperTermination(t *testing.T) {
	doWork := func(done <-chan any, strings <-chan string) <-chan any {
		terminated := make(chan any)
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println("Working on", s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan any)
	strings := make(chan string)
	terminated := doWork(done, strings)
	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	for i := range 10 {
		select {
		case strings <- strconv.Itoa(i):
			time.Sleep(150 * time.Millisecond)
		case <-terminated:
			// Join doWork goroutine with the main goroutine
			fmt.Println("Done.")
			return
			// case <-done:
		}
	}
}

// Book Example
// Goroutine leak caused by Block on write to the channel
func Test_BlockOnWrite(t *testing.T) {
	routineReturned := false
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
			routineReturned = true // unreachable code, goroutine will wait forever

			// should be handled like this:
			// for {
			// 	select {
			// 	case randStream <- rand.Int():
			// 	case <-done:
			// 		return
			// 	}
			// }
		}()
		return randStream
	}
	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	assert.False(t, routineReturned)
}

// Fixing previous example
func Test_FixedBlockOnWrite(t *testing.T) {
	routineReturned := false
	newRandStream := func(done <-chan any) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)

			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					routineReturned = true
					return
				}
			}
		}()
		return randStream
	}
	done := make(chan any)
	defer close(done)
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	go func() {
		<-done
		assert.True(t, routineReturned)
	}()
}

// Example:
// The or-channel. Pretty messed up code, just copy-pasted and tested
func Test_OrChannel(t *testing.T) {
	var or func(channels ...<-chan any) <-chan any
	or = func(channels ...<-chan any) <-chan any {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan any)
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
	assert.True(t, true)
}

// Example:
// return errors along the result from goroutine:
func Test_ReturnErrors(t *testing.T) {
	type Result struct {
		Url      string
		Error    error
		Response *http.Response
	}
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Url: url, Response: resp}
				select {

				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}
	log.Println("Checking urls...")
	done := make(chan interface{})
	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("url: %s, error: %v\n", result.Url, result.Error)
			continue
		}
		fmt.Printf("url: %s, Response: %v\n", result.Url, result.Response.Status)
	}
	close(done)

	// stop if there are 3+ errors:
	log.Println("Checking urls until there are 3 or more errors:")
	done = make(chan interface{})
	errCount := 0
	urls = []string{"a", "https://www.google.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
	close(done)

}

// PIPELINES
//
// Generators
//

func strGenerator(done <-chan any, strSlice []string) <-chan string {
	strChan := make(chan string)
	go func() {
		defer close(strChan)
		for _, s := range strSlice {
			select {
			case <-done:
				return
			default:
				strChan <- s
			}
		}

	}()
	return strChan
}

func Test_strGenerator(t *testing.T) {

	strSlice := []string{"a", "b", "c"}

	done := make(chan any)
	defer close(done)
	for s := range strGenerator(done, strSlice) {
		log.Println(s)
	}
}

// I think I'm gonna start collecting these pipeline functions
// in pipeline package (pipeline/main.go)
func Generator[T any](done <-chan any, strSlice []T) <-chan T {
	outChan := make(chan T)
	go func() {
		defer close(outChan)
		for _, s := range strSlice {
			select {
			case <-done:
				return
			case outChan <- s:
			}
		}
	}()
	return outChan
}

func Test_genericGenerator(t *testing.T) {

	strSlice := []string{"a", "b", "c"}
	done := make(chan any)
	defer close(done)
	for s := range Generator(done, strSlice) {
		log.Println(s)
	}

	intSlice := []int{1, 2, 3, 4, 5, 69, 420}
	done = make(chan any)
	defer close(done)
	for s := range Generator(done, intSlice) {
		log.Println(s)
	}
}

func StageFn[T any](done <-chan any, input <-chan T, fn func(v T) T) <-chan T {
	outChan := make(chan T)
	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case outChan <- fn(<-input):
			}
		}

	}()
	return outChan
}

// Repeat repeats the values until done is closed.
func Repeat[T any](done <-chan any, values ...T) <-chan T {
	outChan := make(chan T)
	go func() {
		defer close(outChan)
		for {
			for _, value := range values {
				select {
				case <-done:
					return
				case outChan <- value:
				}
			}
		}

	}()
	return outChan
}

// RepeatFn repeats the result of fn() until done is closed.
func RepeatFn[T any](done <-chan any, fn func() T) <-chan T {
	outChan := make(chan T)
	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case outChan <- fn():
			}
		}

	}()
	return outChan
}

// generate random ints for 100 microseconds
func Test_genericFnGenerator(t *testing.T) {

	done := make(chan any)
	go func() {
		// time.Sleep(100 * time.Microsecond)
		// close(done)
		// or like this:
		<-time.After(100 * time.Microsecond)
		close(done)
	}()

	for i := range RepeatFn(done, rand.Int) {
		log.Println(i)
	}
}

// Take takes the first num values from the input channel and returns them in a new channel.
func Take[T any](done <-chan any, input <-chan T, num int) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)
		for range num {
			select {
			case <-done:
				return
			case output <- <-input:
			}
		}

	}()
	return output
}

// generate exactly 10 random ints and print half
func Test_10IntsGenerator(t *testing.T) {

	done := make(chan any)
	num := 0
	for i := range Take(done, StageFn(done, RepeatFn(done, rand.Int), func(v int) int { return v / 2 }), 10) {
		log.Println(i)
		num++
	}
	assert.Equal(t, 10, num)
	close(done)

	// Once again but stage-by-stage:

	// done = make(chan any)
	// defer close(done)

	// <-taker<-halver<-generator - why this doesn't work???
	// generator := RepeatFn(done, rand.Int)
	// halver := StageFn(done, generator, func(v int) int { return v / 2 })
	// taker := Take(done, halver, 10)
	// for i := range <-taker {
	// 	log.Println(i)
	// }

}

func Test_GeneratortoStage(t *testing.T) {

	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	results := []int{}
	done := make(chan any)
	for n := range Take(done, StageFn(done, Generator(done, input), func(v int) int {
		return v * 2
	}), 10) {
		results = append(results, n)
	}
	assert.Equal(t, []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, results)
	close(done)

	// <-taker<-halver<-generator

	done = make(chan any)
	defer close(done)
	generator := RepeatFn(done, rand.Int)
	halver := StageFn(done, generator, func(v int) int { return v / 2 })
	taker := Take(done, halver, 10)
	for i := range taker {
		log.Println(i)
		time.Sleep(100 * time.Millisecond)
	}

}

// fan-in pattern - merge multiple channels into one
func fanIn(done <-chan any, channels ...<-chan any) <-chan any {
	var wg sync.WaitGroup
	multiplexedStream := make(chan any)

	multiplex := func(c <-chan any) {
		defer wg.Done()

		for val := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- val:
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go multiplex(ch)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

// some worker function that takes input channel and returns output channel
func primeFinderMock(done <-chan any, input <-chan any) <-chan any {
	outChan := make(chan any)
	go func() {
		defer close(outChan)
		for {
			select {
			case <-done:
				return
			case outChan <- <-input:
			}
		}
	}()
	return outChan
}

// Using fan-in to merge slice of channels into a single channel
func Test_FanIn(t *testing.T) {
	done := make(chan any)
	defer close(done)

	start := time.Now()
	rand := func() any { return rand.Intn(50000000) }
	randIntStream := RepeatFn(done, rand)
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Not Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinderMock(done, randIntStream)
	}
	var i int
	for prime := range Take(done, fanIn(done, finders...), 10) {
		i++
		fmt.Printf("%d\t%d\n", i, prime)
	}
	fmt.Printf("Search took: %v\n", time.Since(start))
}

// or-done-channel (or done channel)
/*
	// allows to make simple loops like this one:
	for val := range OrDone(done, myChan) {
		// Do something with val
	}

	//instead of monstrocity like this:
	loop:
	for {
		select {
		case <-done:
			break loop
		case maybeVal, ok := <-myChan:
			if ok==false {
				return // or maybe break from for
			}
			// Do something with val
		}
	}
*/
func OrDone(done, c <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// Bridge channel
// consume values from a sequence of channels (channel of channels here):
//
//	<-chan <-chan any
//
// destructure the channel of channels into a simple channel:
func Bridge(done <-chan any, chanStream <-chan <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {

			var stream <-chan any
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			for val := range OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func Test_BridgeChannel(t *testing.T) {
	genVals := func() <-chan <-chan any {
		chanStream := make(chan (<-chan any))
		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan any, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}

		}()
		return chanStream
	}
	i := 0
	for v := range Bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
		assert.Equal(t, i, v)
		i++
	}
	assert.Equal(t, 10, i)
}

// tee-channel - receives values from one channel and passes it to two separate channels
func Tee(done <-chan any, in <-chan any) (_, _ <-chan any) {
	out1 := make(chan any)
	out2 := make(chan any)
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range OrDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func Test_TeeChannel(t *testing.T) {
	done := make(chan any)
	defer close(done)
	out1, out2 := Tee(done, Take(done, Repeat[any](done, 1, 2), 4))
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}

// CONTEXT

// Context package serves two primary purposes:
// - to provide an API for cancelling branches in our call-graph
// - to provide a data-bag to transport request-scoped data through our call-graph

// Cancellation:
// - a goroutine's parent may want to cancel it
// - a goroutime may want to cancel its children
// - any blocking operation with a goroutine nedd to be preemptable so it may be cancelled

func Test_WithTimeout(t *testing.T) {

	fLongJob := func(ctx context.Context) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(1 * time.Minute): // simulate long job
			// case <-time.After(500 * time.Millisecond): // Will succeed - no timeout
		}
		return "result", nil
	}

	f1Second := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		result, err := fLongJob(ctx)
		if err != nil {
			return err
		}
		log.Printf("f1Second received result: %s", result)
		return nil
	}

	f1Minute := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		result, err := fLongJob(ctx)
		if err != nil {
			return err
		}
		log.Printf("f1Minute received result: %s", result)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := f1Second(ctx)
		if err != nil {
			log.Printf("error running f1Second: %v\n", err)
			cancel()
			return
		}
		log.Printf("f1Second success\n")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := f1Minute(ctx)
		if err != nil {
			log.Printf("error running f1Minute: %v\n", err)
			cancel()
			return
		}
		log.Printf("f1Minute success\n")
	}()

	wg.Wait()
	log.Println("Done")
}

// tinkering with Cancel and Timeout contexts
func Test_WithTimeout_Callstack(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	top := func(ctx context.Context, to time.Duration, abruptCancel bool) error {
		ctx, cancel := context.WithCancel(ctx) // another context
		defer cancel()

		err := func(ctx context.Context) error { // middle function
			ctx, cancel := context.WithCancel(ctx) // and another context
			defer cancel()

			if abruptCancel {
				cancel()
			}

			err := func(ctx context.Context) error { // last function on the call stack
				ctx, cancel := context.WithTimeout(ctx, time.Second) // and yet another context
				defer cancel()

				select {
				case <-ctx.Done():
					return fmt.Errorf("error: %w\n", ctx.Err())
					// case <-time.After(1 * time.Minute): // simulate long job
				case <-time.After(to): // Will succeed - no timeout
				}

				fmt.Println("last function done")
				return nil
			}(ctx)
			if err != nil {
				return fmt.Errorf("last function error: %w\n", err)
			}

			fmt.Println("middle function done")
			return nil
		}(ctx)
		if err != nil {
			return fmt.Errorf("middle function error: %w\n", err)
		}

		fmt.Println("top function done")
		return nil
	}

	err := top(ctx, 500*time.Millisecond, false) // timeout was not reached, no error
	assert.NoError(t, err, "timeout should't be reached")

	err = top(ctx, 1500*time.Millisecond, false) // timeout reached, error returned
	if err != nil {
		log.Printf("intended error: %v", err)
	}
	assert.Error(t, err, "timeout should definitely be reached")

	err = top(ctx, 1500*time.Millisecond, true) // abrupt cancel in the middle of the callstack
	if err != nil {
		log.Printf("intended abrupt error: %v", err)
	}
	assert.Error(t, err, "abrupt context cancel should terminate everything")

}

// context.WithValue: just avoid using it, it's not worth it :/ It's:
// - type unsafe
// - complex heuristics to determine WHAT should be saved
// - definitely DON'T use it to pass optional parameters!

// CONCURRENCY AT SCALE
//
// Error propagation
// Important to relay some critical information:
// - what happened
// - where and when it occured - call stack is helpful
// - a friendly user-faced messaage
// - how to get more information (stack trace, error id)
// Categories of errors:
// - bug
// - known edge case
//
// Wrap the errors at each level with more context

// Hearthbeats

func doWork(done <-chan interface{}, pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)
		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}
		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}

func Test_Hearthbeat(t *testing.T) {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(done)
	})
	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}

func doWorkFail(done <-chan interface{}, pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)
		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}
		for i := 0; i < 2; i++ {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}

func Test_DoWorkFail(t *testing.T) {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })
	const timeout = 2 * time.Second
	heartbeat, results := doWorkFail(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r)
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

// hearthbeat that happens at the beginning of a unit of work
func Test_Hearthbeat_At_Beginning(t *testing.T) {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)
			for i := 0; i < 10; i++ {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}
				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()
		return heartbeatStream, workStream
	}
	done := make(chan interface{})
	defer close(done)
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

//

func ProcessNumbers(
	done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)
		time.Sleep(2 * time.Second) //simulate the delay before goroutine starts working
		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()
	return heartbeat, intStream
}

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := ProcessNumbers(done, intSlice...)
	<-heartbeat
	i := 0
	for r := range results {
		assert.Equal(t, intSlice[i], r)
		i++
	}
}

// Replicated requests
// get result from a worker that finished first
func Test_ReplicateResult(t *testing.T) {
	doWork := func(
		done <-chan interface{}, id int,
		wg *sync.WaitGroup, result chan<- int,
	) {
		started := time.Now()
		defer wg.Done()

		// Simulate random load
		simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second

		select {
		case <-done:
		case <-time.After(simulatedLoadTime):
		}
		select {
		case <-done:
		case result <- id:
		}
		took := time.Since(started)
		// Display how long handlers would have taken
		if took < simulatedLoadTime {
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %v\n", id, took)
	}

	done := make(chan interface{})

	result := make(chan int)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ { // start 10 handlers to handle our requests.
		go doWork(done, i, &wg, result)
	}
	firstReturned := <-result // grabs the first returned value from the group of handlers.
	close(done)               // cancel all the remaining handlers.
	// 							 This ensures they don’t continue to do unnecessary work.
	wg.Wait()
	fmt.Printf("Received an answer from #%v\n", firstReturned)
}

// Rate Limiting
// token bucket algorithm
// the bucket has a depth (capacity) of - d
// the rate at which tokens replenished - r

func Open() *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1), // 1 request per second
	}
}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func (a *APIConnection) ReadFile(ctx context.Context) error { // Pretend we do work here
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}
func (a *APIConnection) ResolveAddress(ctx context.Context) error { // Pretend we do work here
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func Test_RateLimiting1(t *testing.T) {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	wg.Wait()
}

// Multilimiter

type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit) // sort limiters by restriction ascending
	return &multiLimiter{limiters: limiters}
}

type multiLimiter struct {
	limiters []RateLimiter
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}
func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit() // choose the most restrictive limiter
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

type APIConnectionMulti struct {
	apiLimit RateLimiter
}

func (a *APIConnectionMulti) ReadFile(ctx context.Context) error {
	if err := a.apiLimit.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}
func (a *APIConnectionMulti) ResolveAddress(ctx context.Context) error {
	if err := a.apiLimit.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func OpenMulti() *APIConnectionMulti {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)   // limit per second with no burstiness
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) // 10 per minute with burstiness of 10
	return &APIConnectionMulti{
		apiLimit: MultiLimiter(secondLimit, minuteLimit), // combine two limits
	}
}

func Test_MultiRateLimiting(t *testing.T) {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	apiConnection := OpenMulti()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	wg.Wait()
}

// Limit multiple resources : disk, network, api

type APIConnectionMultiResource struct {
	apiLimit,
	diskLimit,
	networkLimit RateLimiter
}

func OpenMultiResource() *APIConnectionMultiResource {
	return &APIConnectionMultiResource{
		apiLimit: MultiLimiter(
			rate.NewLimiter(Per(2, time.Second), 2),
			rate.NewLimiter(Per(10, time.Minute), 10),
		),
		diskLimit: MultiLimiter(
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: MultiLimiter(
			rate.NewLimiter(Per(3, time.Second), 3),
		)}
}

func (a *APIConnectionMultiResource) ReadFile(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.diskLimit).Wait(ctx)
	if err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}
func (a *APIConnectionMultiResource) ResolveAddress(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.networkLimit).Wait(ctx)
	if err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func Test_MultiResourceRateLimiting(t *testing.T) {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	apiConnection := OpenMultiResource()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	wg.Wait()
}

// rate.Limiter has a few other capabilities ...

// TODO: not reaaly grasped what's going on here >>>

// Healing Unhealthy Goroutines
// steward - monitors a health of a goroutine and restarts it if it's unhealthy
// ward - working gorouting that's monitored

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-or(append(channels[2:], orDone)...):
		}
	}()
	return orDone
}

type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration,
) (heartbeat <-chan interface{})

var newSteward = func(
	timeout time.Duration,
	startGoroutine startGoroutineFn,
) startGoroutineFn {
	return func(
		done <-chan interface{},
		pulseInterval time.Duration) <-chan interface{} {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)
			var wardDone chan interface{}
			var wardHeartbeat <-chan interface{}
			startWard := func() {
				wardDone = make(chan interface{})
				wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2)
			}
			startWard()
			pulse := time.Tick(pulseInterval)
		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
}

func Test_Steward_Simplistic_Ward(t *testing.T) {
	// ward is really simplistic - takes no parameters and returns no arguments
	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("ward: Hello, I'm irresponsible!")
		go func() {
			<-done // not sending any pulses, not doing anything
			log.Println("ward: I am halting.")
		}()
		return nil
	}
	doWorkWithSteward := newSteward(4*time.Second, doWork) // 4 sec timeout for doWork
	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() { // end the test after 9 seconds
		log.Println("main: halting steward and ward.")
		close(done)
	})

	for range doWorkWithSteward(done, 4*time.Second) { // start the steward
		// range over the steward pulses
	}
	log.Println("Done")
}

// let's make some meaningful ward:
var doWorkFn = func(
	done <-chan interface{}, intList ...int,
) (startGoroutineFn, <-chan interface{}) {
	intChanStream := make(chan (<-chan interface{}))
	intStream := Bridge(done, intChanStream)
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration) <-chan interface{} {
		intStream := make(chan interface{})
		heartbeat := make(chan interface{})
		go func() {
			defer close(intStream)
			select {
			case intChanStream <- intStream:
			case <-done:
				return
			}
			pulse := time.Tick(pulseInterval)
			for {
			valueLoop:
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative value: %v\n", intVal)
						return
					}
					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return heartbeat
	}
	return doWork, intStream
}

func Test_Steward_Meaningful_Ward(t *testing.T) {
	log.SetFlags(log.Ltime | log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)
	doWorkWithSteward := newSteward(1*time.Second, doWork)
	doWorkWithSteward(done, 1*time.Hour)

	for intVal := range Take(done, intStream, 6) {
		fmt.Printf("Received: %v\n", intVal)
	}
}

// Some more concurrency patterns drilling, during Ardanlabs Ultimate Go videos:

// Fan-out pattern
func Test_FanOut(t *testing.T) {

	total := 2000

	results := make(chan int, total)
	for i := range total { // spawning 2000 goroutines
		go func() {
			time.Sleep(200 * time.Millisecond)
			results <- i
		}()
	}

	totalReceived := 0
	for range total {
		<-results
		totalReceived++
	}
	assert.Equal(t, totalReceived, total)
}

// fan-out with semaphore rate limit
func Test_RateLimitedFanOut(t *testing.T) {
	total := 2000
	results := make(chan int, total)

	limit := 200
	limiter := make(chan any, limit)

	for i := range total {
		go func() {
			limiter <- true
			time.Sleep(100 * time.Millisecond)
			results <- i
		}()
	}

	totalReceived := 0
	for range total {
		<-results
		<-limiter
		totalReceived++
	}
	assert.Equal(t, totalReceived, total)
}

// wait-for-task, fancy version
func waitForTask(sendTasks, receiveTasks int, sWait, rWait time.Duration) int32 {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan any)
	tasks := make(chan int)
	var totalReceived int32

	// worker
	go func(context.Context) {
		defer close(done)
		for range receiveTasks {
			select {
			case task, ok := <-tasks:
				if !ok {
					return
				}
				atomic.AddInt32(&totalReceived, 1)
				log.Println("task received:", task)
				time.Sleep(rWait) // simulate delay on receive side
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	for i := range sendTasks {
		select {
		case tasks <- i:
			time.Sleep(sWait) // sumulate delay on the send side
		case <-ctx.Done():
			return atomic.LoadInt32(&totalReceived)
		}
	}
	cancel() // don't wait till timeout happens
	<-done
	return totalReceived
}

func Test_WaitForTask(t *testing.T) {

	// testing with different values

	// happy path
	received := waitForTask(5, 5, time.Duration(0), time.Duration(0))
	assert.Equal(t, int32(5), received)

	// waiting for 6 result, but sending 5 tasks
	received = waitForTask(5, 6, time.Duration(0), time.Duration(0))
	assert.Equal(t, int32(5), received)

	// waiting for 5 result, but sending 6 tasks
	// timed out, but correctly received 5 results
	received = waitForTask(6, 5, time.Duration(0), time.Duration(0))
	assert.Equal(t, int32(5), received)

	// 1 second delay after each send
	// only 2-3 tasks processed, 2-3 results received before timeout
	received = waitForTask(6, 5, time.Second, time.Duration(0))
	assert.LessOrEqual(t, int32(2), received)

	// 1 second delay after each receive
	// only 2-3 tasks processed, 2-3 results received before timeout
	received = waitForTask(6, 5, time.Duration(0), time.Second)
	assert.LessOrEqual(t, int32(2), received)

	// 1 second delay after each receive and send
	// only 1-2 tasks processed, 1-2 results received before timeout
	// second task can come before timeout hits
	received = waitForTask(6, 5, time.Duration(0), time.Second)
	assert.LessOrEqual(t, int32(1), received)
}

// Bill's example actually didn't work without waitgroup
func Test_BasicPooling(t *testing.T) {

	totalTasks := 100
	var received atomic.Int32

	wg := sync.WaitGroup{}

	g := runtime.NumCPU()
	tasks := make(chan int)
	for wi := range g {
		go func(worker int) {
			for task := range tasks {
				fmt.Println("Worker", wi, "received task", task)
				received.Add(1)
				wg.Done()
			}
			fmt.Println("Worker", wi, "shut down")
		}(wi)
	}

	for task := range totalTasks {
		wg.Add(1)
		tasks <- task
	}
	close(tasks)

	wg.Wait()
	assert.Equal(t, int32(totalTasks), received.Load())
}

// limited amount of goroutines working on a bunch of tasks
func Test_BoundedFanout(t *testing.T) {

	var received atomic.Int32

	workers := runtime.NumCPU()

	wg := sync.WaitGroup{}
	wg.Add(workers)

	tasksCh := make(chan int)
	// receiving tasks, dispatching workers
	for wi := range workers {
		go func(worker int) {
			defer wg.Done()
			for task := range tasksCh {
				received.Add(1)
				log.Println("Task", task, "received by worker", worker)
			}
		}(wi)
	}

	// sending tasks
	totalTasks := 200
	for task := range totalTasks {
		tasksCh <- task
	}
	close(tasksCh)
	wg.Wait()

	assert.Equal(t, int32(totalTasks), received.Load())
}

// Drop pattern - drop tasks after some cap
// Bill's example with time.Sleep is incorrect again
func Test_DropPattern(t *testing.T) {

	var received atomic.Int32

	totalTasks := 2000
	cap := 200

	done := make(chan any)

	taskCh := make(chan int, cap)
	// receiving side
	go func() {
		defer close(done)
		for range taskCh {
			received.Add(1)
			// log.Println("received task")
		}
	}()

	// send side
	for task := range totalTasks {
		select {
		case taskCh <- task:
			//task sent
		default:
			//task dropped
		}
	}

	close(taskCh)
	<-done
	assert.LessOrEqual(t, int32(200), received.Load(), "~200 should be received")
}

// Cancellation pattern

func Test_Cancel_Buffered(t *testing.T) {
	timeout := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskCh := make(chan any, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		taskCh <- "result"
	}()

	select {
	case res := <-taskCh:
		log.Println("result received", res)
	case <-ctx.Done():
		log.Println("timeout reached")
	}

	time.Sleep(time.Second)
	log.Println("---- end")
}

func Test_Cancel_UnBuffered(t *testing.T) {
	timeout := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskCh := make(chan any)
	go func() {
		time.Sleep(200 * time.Millisecond)
		select {
		case taskCh <- "result":
			log.Println("result sent")
		case <-ctx.Done():
			log.Println("timeout reached on send side")
		}
	}()

	select {
	case res := <-taskCh:
		log.Println("result received", res)
	case <-ctx.Done():
		log.Println("timeout reached")
	}

	time.Sleep(time.Second)
	log.Println("---- end")
}
