package linkedlist

import (
	"errors"
	"slices"
)

type Element struct {
	data int
	next *Element
}

type List struct {
	head *Element
}

func New(elements []int) *List {
	l := &List{}
	for _, e := range elements {
		l.Push(e)
	}
	return l
}

func (l *List) Size() int {
	if l.head == nil {
		return 0
	}

	current := l.head
	size := 1
	for current.next != nil {
		current = current.next
		size++
	}
	return size
}

func (l *List) Push(element int) {
	new := Element{
		data: element,
		next: nil,
	}

	if l.head == nil {
		l.head = &new
		return
	}

	current := l.head
	for current.next != nil {
		current = current.next
	}
	current.next = &new
}

func (l *List) Pop() (int, error) {
	if l.head == nil {
		return 0, errors.New("empty list")
	}

	// absolute mess, but it works somehow
	current := l.head
	prev := current
	for {
		if current.next == nil {
			v := current.data
			if l.head == current {
				l.head = nil
			}
			prev.next = nil
			current = nil
			return v, nil
		}
		prev = current
		current = current.next
	}
}

func (l *List) Array() []int {
	arr := []int{}
	if l.head == nil {
		return arr
	}

	current := l.head
	for {
		arr = append(arr, current.data)
		if current.next == nil {
			return arr
		}
		current = current.next
	}
}

func (l *List) Reverse() *List {
	arr := l.Array()
	slices.Reverse(arr)
	return New(arr)
}
