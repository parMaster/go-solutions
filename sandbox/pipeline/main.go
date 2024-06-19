package pipeline

import "sync"

// Generator converts any slice to channel of values
func Generator[T any](done <-chan any, inputSlice []T) <-chan T {
	outChan := make(chan T)
	go func() {
		defer close(outChan)
		for _, v := range inputSlice {
			select {
			case <-done:
				return
			case outChan <- v:
			}
		}
	}()
	return outChan
}

// StageFn evaluates fn(value) for every value in the pipeline
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

// Concurrency patterns

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
func OrDone[T any](done <-chan any, c <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
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
func Bridge[T any](done <-chan any, chanStream <-chan <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {

			var stream <-chan T
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
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

// tee-channel - receives values from one channel and passes it to two separate channels
func Tee[T any](done <-chan any, in <-chan T) (_, _ <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
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

// fan-in pattern - merge multiple channels into one
func FanIn[T any](done <-chan any, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	multiplexedStream := make(chan T)

	multiplex := func(c <-chan T) {
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
