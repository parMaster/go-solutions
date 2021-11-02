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

func Test_BFT(t *testing.T) {

	var testTree = &Node{value: "a",
		left: &Node{value: "b",
			left:  &Node{value: "d"},
			right: &Node{value: "e"},
		},
		right: &Node{value: "c",
			right: &Node{value: "f"},
		},
	}

	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, testTree.BreadthFirstTraversal())
}
