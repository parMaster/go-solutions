/*
	Binary Tree Algorithms for Technical Interviews - Full Course
	https://youtu.be/fAAZixBzIAI

	Vocabulary:

	Tree
	Node
	Parent
	Child, Children
	Root node - has no parent
	Leaf node - has no children

	Definitions:

	BINARY TREE criteria:
	- Every node has at most 2 children (0, 1 or 2)
	- Exactly 1 root
	- Exactly 1 path between root and any node

	Empty tree - tree with zero nodes - special case of binary tree

	Left and Right node - children on the left and right

	Binary tree has no cycles

	Representation:

	type Node struct {
		value	interface{}
		left	Node
		right	Node
	}

*/

package sandbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateTreeManually(t *testing.T) {

	var n Node
	n.value = "a"
	n.left = &Node{value: "b", left: &Node{value: "d"}, right: &Node{value: "e"}}
	n.right = &Node{value: "c", right: &Node{value: "f"}}

	// Or everythingliterally:
	var tree = &Node{value: "a",
		left: &Node{value: "b",
			left:  &Node{value: "d"},
			right: &Node{value: "e"},
		},
		right: &Node{value: "c",
			right: &Node{value: "f"},
		},
	}

	fmt.Println(tree.left.right)
}

/* Problem 1. Depth First Value
- Using Depth First Traversal Algo
- and LIFO Stack

Linear time and space solution
n = # of nodes
Time complexity: O(n)
Spce Complexity: O(n)
*/
func Test_dfv(t *testing.T) {

	var tree = Node{value: "a",
		left: &Node{value: "b",
			left:  &Node{value: "d"},
			right: &Node{value: "e"},
		},
		right: &Node{value: "c",
			right: &Node{value: "f"},
		},
	}

	var res string
	s := NewStack()
	s.push(tree)

	for !s.isEmpty() {

		currentNode := s.pop()

		res = res + string(currentNode.value)
		if currentNode.right != nil {
			s.push(*currentNode.right)
		}

		if currentNode.left != nil {
			s.push(*currentNode.left)
		}
	}

	assert.Equal(t, "abdecf", res)
}

// Extract it in the new method

// linearDFV - performs linear time and space complexity depth for value algo, using LIFO stack
func (n *Node) linearDFV() []string {
	var res []string

	s := NewStack()
	s.push(*n)

	for !s.isEmpty() {
		currentNode := s.pop()

		res = append(res, currentNode.value)
		if currentNode.right != nil {
			s.push(*currentNode.right)
		}

		if currentNode.left != nil {
			s.push(*currentNode.left)
		}
	}
	return res
}

// Problem 2 - Recursive tree traversal
// recursiveDFV ...
func (n *Node) recursiveDFV() []string {

	var leftValues []string
	if n.left != nil {
		leftValues = n.left.recursiveDFV()
	}

	var rightValues []string
	if n.right != nil {
		rightValues = n.right.recursiveDFV()
	}

	return append([]string{n.value}, append(leftValues, rightValues...)...)
}

func Test_DFV(t *testing.T) {

	var testTree = &Node{value: "a",
		left: &Node{value: "b",
			left:  &Node{value: "d"},
			right: &Node{value: "e"},
		},
		right: &Node{value: "c",
			right: &Node{value: "f"},
		},
	}

	assert.Equal(t, []string{"a", "b", "d", "e", "c", "f"}, testTree.linearDFV())

	var emptyTree = &Node{}
	assert.Equal(t, []string{""}, emptyTree.linearDFV())

	assert.Equal(t, []string{"a", "b", "d", "e", "c", "f"}, testTree.recursiveDFV())
}

// Problem 3. Breadth First Traversal
// Using  queue instead of stack (FIFO)
func (n *Node) BreadthFirstTraversal() []string {
	var result []string

	s := NewStack()
	s.push(*n)

	for !s.isEmpty() {

		current := s.popFirst()

		result = append(result, current.value)

		if current.left != nil {
			s.push(*current.left)
		}
		if current.right != nil {
			s.push(*current.right)
		}
	}

	return result
}

// Problem 4. Tree Includes
// Given the value, check if there is a node with such value

// Breadth-first search
func (n *Node) TreeIncludesBFS(value string) bool {

	s := NewStack()
	s.push(*n)

	for !s.isEmpty() {

		current := s.popFirst()

		if current.value == value {
			return true
		}

		if current.left != nil {
			s.push(*current.left)
		}

		if current.right != nil {
			s.push(*current.right)
		}
	}

	return false
}

// Depth-first search
func (n *Node) TreeIncludesDFS(value string) bool {

	s := NewStack()
	s.push(*n)

	for !s.isEmpty() {
		currentNode := s.pop()

		if currentNode.value == value {
			return true
		}

		if currentNode.right != nil {
			s.push(*currentNode.right)
		}

		if currentNode.left != nil {
			s.push(*currentNode.left)
		}
	}

	return false
}

// Depth-first search Recursivly
func (n *Node) TreeIncludesDFS_Recursive(value string) bool {
	result := false

	if n == nil {
		return false
	}

	if n.value == value {
		return true
	}

	if n.right != nil {
		result = result || n.right.TreeIncludesDFS_Recursive(value)
	}
	if n.left != nil {
		result = result || n.left.TreeIncludesDFS_Recursive(value)
	}

	return result
}

// Problem 5. Tree Sum
// return sum of numeric node values
func (n *Node) treeSum_Recursive() int {
	var result int

	result += n.intValue

	if n.right != nil {
		result += n.right.treeSum_Recursive()
	}

	if n.left != nil {
		result += n.left.treeSum_Recursive()
	}

	return result
}

func (n *Node) treeSum() int {
	var result int

	s := NewStack()
	s.push(*n)

	for !s.isEmpty() {

		currentNode := s.popFirst()

		result += currentNode.intValue

		if n.right != nil {
			result += n.right.treeSum()
		}

		if n.left != nil {
			result += n.left.treeSum()
		}
	}

	return result
}

// Problem 6. Tree min value
// Find the smallest value node in the tree
func (n *Node) treeMin_Recursive() int {
	var result int = 1<<32 - 1

	result = min(result, n.intValue)

	if n.right != nil {
		result = min(result, n.right.treeMin_Recursive())
	}

	if n.left != nil {
		result = min(result, n.left.treeMin_Recursive())
	}

	return result
}

// Same, iteratively
func (n *Node) treeMin() int {
	var result int = 1<<32 - 1

	s := NewStack()
	s.push(*n)

	for !s.isEmpty() {
		currentNode := s.popFirst()

		result = min(result, currentNode.intValue)

		if currentNode.left != nil {
			s.push(*currentNode.left)
		}

		if currentNode.right != nil {
			s.push(*currentNode.right)
		}
	}

	return result
}

func Test_Everything(t *testing.T) {

	var testTree = &Node{value: "a", intValue: 3,
		left: &Node{value: "b", intValue: 11,
			left:  &Node{value: "d", intValue: 4},
			right: &Node{value: "e", intValue: 2},
		},
		right: &Node{value: "c", intValue: 4,
			right: &Node{value: "f", intValue: 1},
		},
	}

	var emptyTree = &Node{}

	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, testTree.BreadthFirstTraversal())

	searchTests := []struct {
		tree     Node
		needle   string
		expected bool
	}{
		{tree: *testTree, needle: "e", expected: true},
		{tree: *testTree, needle: "not there", expected: false},
		{tree: *testTree, needle: "", expected: false},
		{tree: *emptyTree, needle: "asd", expected: false},
	}

	for i := range searchTests {
		assert.Equal(t, searchTests[i].expected, searchTests[i].tree.TreeIncludesBFS(searchTests[i].needle))
		assert.Equal(t, searchTests[i].expected, searchTests[i].tree.TreeIncludesDFS(searchTests[i].needle))
		assert.Equal(t, searchTests[i].expected, searchTests[i].tree.TreeIncludesDFS_Recursive(searchTests[i].needle))
	}

	assert.Equal(t, 25, testTree.treeSum_Recursive())
	assert.Equal(t, 25, testTree.treeSum())

	assert.Equal(t, 1, testTree.treeMin_Recursive())
	assert.Equal(t, 1, testTree.treeMin())

	testTree.right.right.intValue = 99
	assert.Equal(t, 2, testTree.treeMin_Recursive())
	assert.Equal(t, 2, testTree.treeMin())
}
