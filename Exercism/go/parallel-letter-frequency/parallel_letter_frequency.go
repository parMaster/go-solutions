package letter

import "runtime"

// FreqMap records the frequency of each rune in a given text.
type FreqMap map[rune]int

// Frequency counts the frequency of each rune in a given text and returns this
// data as a FreqMap.
func Frequency(text string) FreqMap {
	frequencies := FreqMap{}
	for _, r := range text {
		frequencies[r]++
	}
	return frequencies
}

// basic fan-out with semaphore rate limit
func ConcurrentFrequency(texts []string) FreqMap {

	limit := runtime.NumCPU()
	limiter := make(chan any, limit)

	results := make(chan FreqMap, len(texts))
	for _, text := range texts {
		go func(text string) {
			limiter <- struct{}{}
			freq := FreqMap{}
			for _, c := range text {
				freq[c]++
			}
			results <- freq
			<-limiter
		}(text)
	}

	fm := FreqMap{}
	for i := 0; i < len(texts); i++ {
		res := <-results
		for r, freq := range res {
			fm[r] += freq
		}
	}

	return fm
}

// basic fan-out
func ConcurrentFrequency_basic(texts []string) FreqMap {

	results := make(chan FreqMap, len(texts))
	for _, text := range texts {
		go func(text string) {
			freq := FreqMap{}
			for _, c := range text {
				freq[c]++
			}
			results <- freq
		}(text)
	}

	fm := FreqMap{}
	for i := 0; i < len(texts); i++ {
		res := <-results
		for r, freq := range res {
			fm[r] += freq
		}
	}

	return fm
}
