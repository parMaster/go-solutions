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
