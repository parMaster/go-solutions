package letter

import (
	"log"
	"runtime"
	"time"
)

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

func ConcurrentFrequency_selects(texts []string) FreqMap {

	countFreqs := func(done <-chan any, texts <-chan string) (freqMap chan FreqMap, terminated chan any) {
		terminated = make(chan any)
		freqMap = make(chan FreqMap)

		go func() {
			defer close(terminated)
			for {
				select {
				case text := <-texts:
					fm := FreqMap{}
					for _, c := range text {
						fm[c]++
					}
					freqMap <- fm
				case <-done:
					return
				}
			}
		}()

		return freqMap, terminated
	}

	done := make(chan any)
	textsCh := make(chan string)

	freqMapCh, terminated := countFreqs(done, textsCh)

	fm := FreqMap{}

	// timeout circuit breaker
	go func() {
		select {
		case <-terminated:
			close(done)
		case <-time.After(200 * time.Second):
			close(done)
		}
	}()

	go func() {
		for _, text := range texts {
			textsCh <- text
			log.Println("text sent")
		}
	}()

	for {
		select {
		case <-done:
			log.Println("terminated (done)")
			return fm
		case <-terminated:
			log.Println("terminated")
			return fm
		case result := <-freqMapCh:
			for c, freq := range result {
				fm[c] += freq
			}
			log.Println("result received")
		}
		time.Sleep(1 * time.Microsecond)
	}
}
