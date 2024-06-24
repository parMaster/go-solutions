package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type Value struct{}

func JustError() error {
	g := new(errgroup.Group)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.somestupidname.com/",
		"http://www.google.com/",
	}
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			// Fetch the URL.
			resp, err := http.Get(url)
			fmt.Println("fetching", url)
			if err == nil {
				resp.Body.Close()
			}
			return err
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
