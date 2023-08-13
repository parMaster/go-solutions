package sandbox

import (
	"fmt"
	"testing"
)

func Test_Iter(t *testing.T) {
	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
		// v := v // create a new 'v'.
		go func() {
			fmt.Println(v) // loop variable v captured by func literal
			done <- true
		}()
	}

	// wait for all goroutines to complete before exiting
	for range values {
		<-done
	}
}
