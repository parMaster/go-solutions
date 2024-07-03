// context-aware fan-out using errgroup
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func doJob(ctx context.Context, job int) (res int, err error) {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	select {
	case <-ctx.Done():
		log.Println("context cancelled: ctx.Err() ")
		return 0, ctx.Err()
	default:
		if job == 5 {
			err := fmt.Errorf("job == 5")
			log.Println(err)
			return 0, err
		}
		return job * job, nil
	}
}

// context-aware fan-out using errgroup
func Test_ErrgroupFanOut(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	resCh := make(chan int)

	g, ctx := errgroup.WithContext(context.Background())
	for _, job := range jobs {
		g.Go(func() error {
			res, err := doJob(ctx, job)
			if err != nil {
				return err
			}
			resCh <- res
			return nil
		})
	}

	defer log.Println("jobs done")

	go func() {
		err := g.Wait()
		close(resCh)
		if err != nil {
			fmt.Println(fmt.Errorf("errgroup returned error: %w", err))
		}
	}()

	for res := range resCh {
		log.Println("received result:", res)
	}

}
