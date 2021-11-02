package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_popFirst(t *testing.T) {

	var s Stack

	s.push(Node{value: "a"})
	s.push(Node{value: "b"})
	s.push(Node{value: "c"})
	assert.Equal(t, "c", s.pop().value)
	assert.Equal(t, "b", s.pop().value)
	assert.Equal(t, "a", s.pop().value)

	s.push(Node{value: "a"})
	s.push(Node{value: "b"})
	s.push(Node{value: "c"})
	assert.Equal(t, "a", s.popFirst().value)
	assert.Equal(t, "b", s.popFirst().value)
	assert.Equal(t, "c", s.popFirst().value)
	assert.Equal(t, &Node{}, s.popFirst())
}
