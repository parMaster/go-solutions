package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ListNode[T comparable] struct {
	data T
	next *ListNode[T]
}

type LinkedList[T comparable] struct {
	head *ListNode[T]
}

func NewLinkedList[T comparable]() *LinkedList[T] {
	l := &LinkedList[T]{}
	return l
}

func (list *LinkedList[T]) Add(data T) {
	new := &ListNode[T]{
		data: data,
		next: nil,
	}

	if list.head == nil {
		list.head = new
		return
	}

	current := list.head
	for current.next != nil {
		current = current.next
	}
	current.next = new
}

func (list *LinkedList[T]) Traverse(f func(*ListNode[T])) {
	if list.head == nil {
		return
	}

	current := list.head
	f(current)
	for current.next != nil {
		current = current.next
		f(current)
	}
}

func (list *LinkedList[T]) Size() int {
	size := 0
	if list.head == nil {
		return size
	}

	list.Traverse(func(*ListNode[T]) {
		size++
	})
	return size
}

func Test_AddTraverse(t *testing.T) {

	list := NewLinkedList[string]()

	list.Add("A")
	list.Add("B")
	list.Add("C")

	data := []string{}
	list.Traverse(func(node *ListNode[string]) {
		// fmt.Println(node.data)
		data = append(data, node.data)
	})
	assert.Equal(t, []string{"A", "B", "C"}, data)
}

func Load[T comparable](a []T) *LinkedList[T] {
	list := NewLinkedList[T]()

	for _, v := range a {
		list.Add(v)
	}

	return list
}

func (list *LinkedList[T]) Save() []T {
	a := make([]T, 0)

	list.Traverse(func(node *ListNode[T]) {
		a = append(a, node.data)
	})

	return a
}

func Test_LoadSave(t *testing.T) {
	input := []string{"A", "B", "C"}
	list := Load(input)
	output := list.Save()

	assert.Equal(t, output, input)
}

func (list *LinkedList[T]) Delete(value T) (result bool) {

	if list.head == nil {
		return false
	}
	if list.head.data == value {
		list.head = list.head.next
		return true
	}

	prev := list.head
	list.Traverse(func(node *ListNode[T]) {
		if node.data == value {
			prev.next = node.next
			result = true
		}
		prev = node
	})
	return result
}

func Test_LoadDelete(t *testing.T) {
	// Middle
	list := Load([]string{"A", "B", "C"})
	res := list.Delete("B")
	output := list.Save()

	assert.True(t, res, "delete middle result")
	assert.Equal(t, []string{"A", "C"}, output, "delete middle")

	// Last
	list = Load([]string{"A", "B", "C"})
	res = list.Delete("C")
	assert.True(t, res, "delete last result")
	assert.Equal(t, []string{"A", "B"}, list.Save(), "delete last")

	// First
	list = Load([]string{"A", "B", "C"})
	res = list.Delete("A")
	assert.True(t, res, "delete first result")
	assert.Equal(t, []string{"B", "C"}, list.Save(), "delete first")
}

func Test_LoadSizeDelete(t *testing.T) {
	// Middle
	list := Load([]string{"A", "B", "C"})
	assert.Equal(t, 3, list.Size())
	list.Delete("B")
	assert.Equal(t, 2, list.Size())
}

// ToDo:
// Create list with a cycle
// Cycle detection
