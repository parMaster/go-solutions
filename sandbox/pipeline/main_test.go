package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_genericGenerator(t *testing.T) {
	strSlice := []string{"a", "b", "c"}
	done := make(chan any)
	results := []string{}
	for s := range Generator(done, strSlice) {
		results = append(results, s)
	}
	close(done)
	assert.Equal(t, strSlice, results)

	intSlice := []int{1, 2, 3, 4, 5, 69, 420}
	done = make(chan any)
	defer close(done)
	intResults := []int{}
	for i := range Generator(done, intSlice) {
		intResults = append(intResults, i)
	}
	assert.Equal(t, intSlice, intResults)
}

func Test_StageFn(t *testing.T) {
	done := make(chan any)
	strSlice := []string{"a", "b", "c"}
	results := []string{}

	for s := range Take(done, StageFn(done, Generator(done, strSlice), func(v string) string {
		return v + "!"
	}), len(strSlice)) {
		results = append(results, s)
	}
	close(done)
	assert.Equal(t, []string{"a!", "b!", "c!"}, results)

	done = make(chan any)
	intSlice := []int{1, 2, 3, 4, 5, 69, 420}
	intResults := []int{}
	for i := range Take(done, StageFn(done, Generator(done, intSlice), func(v int) int {
		return v + 1
	}), len(intSlice)) {
		intResults = append(intResults, i)
	}
	close(done)
	assert.Equal(t, []int{2, 3, 4, 5, 6, 70, 421}, intResults)
}

func Test_Repeat(t *testing.T) {
	done := make(chan any)
	results := []int{}

	for val := range Take(done, Repeat(done, 1, 2), 3) {
		results = append(results, val)
	}
	close(done)
	assert.Equal(t, []int{1, 2, 1}, results)
}

func Test_RepeatFn(t *testing.T) {
	done := make(chan any)
	results := []string{}

	for s := range Take(done, RepeatFn(done, func() string {
		return "a"
	}), 3) {
		results = append(results, s)
	}
	close(done)
	assert.Equal(t, []string{"a", "a", "a"}, results)
}

func Test_OrDone(t *testing.T) {
	// close done channel after finishing reading from out
	done := make(chan any)
	defer close(done)
	total := 4
	out := OrDone(done, Take(done, Repeat(done, 1, 2), total))
	i := 0
	for range out {
		i++
	}
	assert.Equal(t, total, i)

	// close done channel before finishing reading from out
	done = make(chan any)
	out = OrDone(done, Take(done, Repeat(done, 1, 2), total))
	i = 0
	for range out {
		i++
		close(done)
	}
	assert.Equal(t, 1, i)
}

func Test_TeeChannel(t *testing.T) {
	done := make(chan any)
	defer close(done)
	total := 4
	out1, out2 := Tee(done, Take(done, Repeat[any](done, 1, 2), total))
	i := 0
	for val1 := range out1 {
		val2 := <-out2
		assert.Equal(t, val1, val2) // receiving identical values into both channels
		i++
	}
	assert.Equal(t, total, i)
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
		assert.Equal(t, i, v)
		i++
	}
	assert.Equal(t, 10, i)
}
