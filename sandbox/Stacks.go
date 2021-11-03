package sandbox

type StackV2 struct {
	items []interface{}
}

// Value semantics for factory function
func NewStackV2() StackV2 {
	return StackV2{}
}

// Pointer semantics for everything else

// push ...
func (s *StackV2) push(elem interface{}) {
	s.items = append(s.items, elem)
}

// pop ...
func (s *StackV2) pop() interface{} {
	if len(s.items) == 0 {
		return ""
	}
	result := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return result
}

// peek - pop without popping
func (s *StackV2) peek() interface{} {
	if len(s.items) == 0 {
		return ""
	}
	return s.items[len(s.items)-1]
}

// popFirst - returns first element of Stack - FIFO behaviour for Stack struct
func (s *StackV2) popFirst() interface{} {
	if len(s.items) == 0 {
		return ""
	}
	result := s.items[0]
	s.items = s.items[1:len(s.items)]
	return result
}

// isEmpty ...
func (s *StackV2) isEmpty() bool {
	return len(s.items) == 0
}
