package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringStack(t *testing.T) {

	s := NewStackV2()

	s.push("a")
	s.push("b")
	s.push("c")

	assert.Equal(t, "c", s.pop())
	assert.Equal(t, []interface{}{"a", "b"}, s.items)
	assert.Equal(t, "b", s.pop())
	assert.Equal(t, "a", s.pop())

	s.push("x")
	s.push("y")
	s.push("z")

	assert.Equal(t, "z", s.peek())
	assert.Equal(t, "z", s.peek())
	assert.Equal(t, "z", s.peek())

	assert.Equal(t, "x", s.popFirst())
	assert.Equal(t, "z", s.peek())
	assert.Equal(t, "y", s.popFirst())
	assert.Equal(t, "z", s.popFirst())

	assert.Equal(t, true, s.isEmpty())

	assert.Equal(t, "", s.popFirst())
	assert.Equal(t, "", s.pop())

}
