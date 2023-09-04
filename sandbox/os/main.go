package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func handleSignal(sig os.Signal) {
	fmt.Println("handleSignal() Caught:", sig)
}

func main() {

	fmt.Printf("Process ID: %d\n", os.Getpid())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	go func() {
		for {
			fmt.Print("+")
			time.Sleep(1 * time.Second)
		}
	}()

	sig := <-sigs
	fmt.Println("\nExiting,", sig, "signal received.")
}
