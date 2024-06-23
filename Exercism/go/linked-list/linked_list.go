package linkedlist

import "errors"

// Define List and Node types here.
// Note: The tests expect Node type to include an exported field with name Value to pass.

type Node struct {
	next  *Node
	prev  *Node
	Value any
}

type List struct {
	head *Node
}

// create new list from a slice of elements
func NewList(elements ...interface{}) *List {
	l := List{}
	for _, el := range elements {
		l.Push(el)
	}
	return &l
}

// pointer to the next node.
func (n *Node) Next() *Node {
	return n.next
}

// pointer to the previous node.
func (n *Node) Prev() *Node {
	return n.prev
}

// insert value at the front of the list
func (l *List) Unshift(v interface{}) {
	if l.head == nil {
		l.Push(v)
		return
	}

	n := Node{
		Value: v,
		next:  l.head,
	}
	l.head.prev = &n
	l.head = &n
}

// insert value at the back of the list.
func (l *List) Push(v interface{}) {
	last := l.Last()

	n := Node{
		Value: v,
	}

	if last == nil {
		l.head = &n
		return
	}

	n.prev = last
	last.next = &n
}

// remove value from the front of the list.
func (l *List) Shift() (interface{}, error) {
	if l.head == nil {
		return nil, errors.New("no elements")
	}

	val := l.head.Value
	l.head = l.head.next
	if l.head != nil {
		l.head.prev = nil
	}
	return val, nil
}

// remove value from the back of the list.
func (l *List) Pop() (interface{}, error) {
	if l.head == nil {
		return nil, errors.New("no elements")
	}

	last := l.Last()

	val := last.Value

	// only one left
	if last == l.head {
		l.head = nil
		return val, nil
	}

	// something left
	last.prev.next = nil
	last = nil
	return val, nil
}

// reverse the linked list.
func (l *List) Reverse() {
	if l.head == nil {
		return
	}

	curr := l.head
	for {
		curr.next, curr.prev = curr.prev, curr.next // flip pointers
		if curr.prev == nil {
			l.head = curr
			break
		}
		curr = curr.prev
	}
}

// returns a pointer to the first node (head).
func (l *List) First() *Node {
	return l.head
}

// returns a pointer to the last node (tail).
func (l *List) Last() *Node {
	if l.head == nil {
		return nil
	}

	last := l.head

	for last.next != nil {
		last = last.next
	}
	return last
}
