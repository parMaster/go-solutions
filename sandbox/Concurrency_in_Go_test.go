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
	"runtime"
	"strconv"
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
var fanIn = func(done <-chan any, channels ...<-chan any) <-chan any {
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

// I'll take a brake from the concurrency patterns for now
// ... to be continued
//
//
//
//

// CONTEXT

// Context package serves two primary purposes:
// - to provide an API for cancelling branches in our call-graph
// - to provide a data-bag to transport request-scoped data through our call-graph

// Cancellation:
// - a goroutine's parent may want to cancel it
// - a goroutime may want to cancel its children
// - any blocking operation with a goroutine nedd to be preemptable so it may be cancelled

func Test_WithTimeout(t *testing.T) {

	fLast := func(ctx context.Context) (string, error) {
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
		result, err := fLast(ctx)
		if err != nil {
			return err
		}
		log.Printf("f1Second received result: %s", result)
		return nil
	}

	f1Minute := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		result, err := fLast(ctx)
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
