package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Value struct{}

func JustError() error {
	g, ctx := errgroup.WithContext(context.Background())
	var urls = []string{
		"http://www.golang.org/",
		"http://www.somestupidname.com/",
		"http://www.somestupidname.com/",
		"http://www.google.com/",
	}
	for _, url := range urls {
		time.Sleep(time.Second)

		// Launch a goroutine to fetch the URL.
		g.Go(func() error {

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Fetch the URL.
				resp, err := http.Get(url)
				fmt.Println("fetching", url)
				if err == nil {
					resp.Body.Close()
				}
				// first returned error that != nil will cancel the context
				return err
			}
		})
	}

	// Wait for all HTTP fetches to complete.
	err := g.Wait()
	if err != nil {
		return err
	}
	fmt.Println("Successfully fetched all URLs.")
	return nil
}

func main() {
	err := JustError()
	fmt.Println(err)

}
