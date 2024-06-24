package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tour/tree"
)

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func Test_Channels1(t *testing.T) {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c

	log.Println(x, y, x+y)
	assert.Equal(t, x, -5)
	assert.Equal(t, y, 17)
	assert.Equal(t, x+y, 12)
}

// Range and close

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func Test_Channels_Rangeandclose(t *testing.T) {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	sum := 0
	for i := range c {
		sum += i
		log.Println(i, sum)
	}
	assert.Equal(t, 88, sum)
}

// Default Selection

func Test_DefaultSelection(t *testing.T) {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// Tree Walk

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {

	var walk func(node *tree.Tree)
	walk = func(node *tree.Tree) {
		if node == nil {
			return
		}
		walk(node.Left)
		ch <- node.Value
		walk(node.Right)
	}

	walk(t)
	close(ch)
}

func Test_TreeWalk(t *testing.T) {
	cap := 10
	tr := tree.New(cap)

	ch := make(chan int)
	go Walk(tr, ch)

	vals := []int{}
	for val := range ch {
		vals = append(vals, val)
	}

	log.Println(vals)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	vals1, vals2 := []int{}, []int{}
	v1Closed, v2Closed := false, false
	for v1Closed == false || v2Closed == false {
		select {
		case v1, ok := <-ch1:
			if !ok {
				v1Closed = true
				continue
			}
			vals1 = append(vals1, v1)
		case v2, ok := <-ch2:
			if !ok {
				v2Closed = true
				continue
			}
			vals2 = append(vals2, v2)
		}
	}
	return reflect.DeepEqual(vals1, vals2)
}

func Test_TreeSame(t *testing.T) {
	cap := 2
	tr1 := tree.New(cap)
	tr2 := tree.New(cap)

	same := Same(tr1, tr2)
	assert.True(t, same)

	assert.False(t, Same(tree.New(cap), tree.New(cap+1)))
}

// Walk using Channel Owner pattern
func WalkOwner(t *tree.Tree) <-chan int {
	channelOwner := func() <-chan int {
		intStream := make(chan int, 4)
		go func() {
			defer close(intStream)

			var walk func(node *tree.Tree)
			walk = func(node *tree.Tree) {
				if node == nil {
					return
				}
				walk(node.Left)
				intStream <- node.Value
				walk(node.Right)
			}
			walk(t)
		}()
		return intStream
	}
	return channelOwner()
}

func Test_TreeWalkOwner(t *testing.T) {
	cap := 10
	tr := tree.New(cap)

	ch := WalkOwner(tr)

	vals := []int{}
	for val := range ch {
		vals = append(vals, val)
	}

	log.Println(vals)
}

// SameOwner uses WalkOwner
// returns true if t1 and t2 contain the same values.
func SameOwner(t1, t2 *tree.Tree) bool {
	ch1 := WalkOwner(t1)
	ch2 := WalkOwner(t2)

	vals1, vals2 := []int{}, []int{}
	v1Closed, v2Closed := false, false
	for v1Closed == false || v2Closed == false {
		select {
		case v1, ok := <-ch1:
			if !ok {
				v1Closed = true
				continue
			}
			log.Println("v1 received")
			vals1 = append(vals1, v1)
		case v2, ok := <-ch2:
			if !ok {
				v2Closed = true
				continue
			}
			log.Println("v2 received")
			vals2 = append(vals2, v2)
		}
	}
	return reflect.DeepEqual(vals1, vals2)
}

func Test_TreeSameOwner(t *testing.T) {
	cap := 2
	tr1 := tree.New(cap)
	tr2 := tree.New(cap)

	same := SameOwner(tr1, tr2)
	assert.True(t, same)

	assert.False(t, SameOwner(tree.New(cap), tree.New(cap+1)))
}

// func main() {
// 	cap := 2
// 	tr1 := tree.New(cap)
// 	tr2 := tree.New(cap)

// 	if !Same(tr1, tr2) {
// 		panic("should be equal")
// 	}

// 	if Same(tree.New(cap), tree.New(cap+1)){
// 		panic("shouldn't be equal")
// 	}
// }
