package letter

import "sync"

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

// basic fan-out
func ConcurrentFrequency(texts []string) FreqMap {

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

// ConcurrentFrequency counts the frequency of each rune in the given strings,
// by making use of concurrency.
// with wait, if len(texts) is unknown
func ConcurrentFrequency_wg(texts []string) FreqMap {
	fm := FreqMap{}

	done := make(chan any)
	results := make(chan FreqMap)
	wg := sync.WaitGroup{}
	for _, text := range texts {
		wg.Add(1)
		go func(text string) {
			fm := FreqMap{}
			for _, r := range text {
				fm[r]++
			}
			results <- fm
		}(text)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case res := <-results:
			for r, freq := range res {
				fm[r] += freq
			}
			wg.Done()
		case <-done:
			return fm
		}
	}
}
