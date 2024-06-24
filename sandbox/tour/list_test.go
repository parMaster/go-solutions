package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func NewList[T any]() *List[T] {
	l := List[T]{}
	return &l
}

func (l *List[T]) Push(val T) {
	nl := List[T]{val: val}
	l.next = &nl
}

func (l *List[T]) LastVal() T {
	curr := l
	for {
		if curr.next == nil {
			return curr.val
		} else {
			curr = curr.next
		}
	}
}

func Test_PushPop(t *testing.T) {
	l := NewList[string]()

	l.Push("abc")
	l.Push("def")

	val := l.LastVal()
	assert.Equal(t, "def", val)
}

// func main() {
// 	l := NewList[string]()
// 	l.Push("abc")
// 	l.Push("def")

// 	val := l.LastVal()
// 	if val != "def" {
// 		panic("wrong value")
// 	} else {
// 		fmt.Println("OK!")
// 	}
// }
