package main

import (
	"fmt"
	"os"
	"sync"
)

func f() (f *os.File) {
	f, err := os.OpenFile("1.txt", os.O_RDWR, 0644)
	f.Truncate(0)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {

	f := f()

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		for i := 8888; i < 9888; i++ {
			f.WriteString(fmt.Sprintf("%d\r\n", i))
			fmt.Println(i)
		}
		wg.Done()
	}()

	wg.Add(1)
	f1 := f
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			f1.WriteString(fmt.Sprintf("%d\n", i))
		}
	}()

	wg.Wait()

}
