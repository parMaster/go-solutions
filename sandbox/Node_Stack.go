package sandbox

type Node struct {
	value    string
	intValue int
	left     *Node
	right    *Node
}

// Stack helper struct with push and pop methods
type Stack struct {
	stack []Node
}

func NewStack() Stack {
	return Stack{stack: []Node{}}
}

func (s *Stack) push(node Node) {
	s.stack = append(s.stack, node)
}

func (s *Stack) pop() Node {
	if len(s.stack) == 0 {
		return Node{}
	}
	result := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return result
}

func (s *Stack) isEmpty() bool {
	return len(s.stack) == 0
}

// popFirst - returns first element of Stack - FIFO behaviour for Stack struct
func (s *Stack) popFirst() *Node {
	if len(s.stack) == 0 {
		return &Node{}
	}
	result := s.stack[0]
	s.stack = s.stack[1:len(s.stack)]
	return &result
}

/*

package sandbox

type Node struct {
	value string
	left  *Node
	right *Node
}

// Stack helper struct with push and pop methods
type Stack []*Node

func NewStack() Stack {
	return Stack{}
}

func (s Stack) push(node Node) {
	s = append(s, &node)
}

func (s Stack) pop() Node {
	if len(s) == 0 {
		return Node{}
	}
	result := s[len(s)-1]
	s = s[:len(s)-1]
	return result
}

func (s Stack) isEmpty() bool {
	return len(s) == 0
}

*/
