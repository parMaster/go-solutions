package pipeline

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
