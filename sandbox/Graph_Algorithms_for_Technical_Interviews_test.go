// Graph Algorithms for Technical Interviews - Full Course
// nodejs examples: https://www.youtube.com/watch?v=tWVWeAqZ0WU

package sandbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// I guess I'll need new graph struct for this one
// Can I use Interfaces this time, can I configure Node, so
// value would be int or string for different kinds of graphs

// Also using simple adjacency list would be nice
/*
	Adjacency list
		{
			a: {b,c},
			b: {d},
			c: {e},
			d: {},
			e: {d},
			f: {b},
		}

	hash table like
	key: string; value: array of strings
*/
type graph map[string][]string

// OR DO AN POLYMORPHISM kind of thing -  graph can be string graph or int graph

// Problem 1. Depth-First Traversal (breadth-first as well)

// Iterative, intuitive method
// in Value Semantics
func (g graph) depthFirstTraversal(source string) []string {
	result := []string{}

	if g[source] == nil {
		return result
	}

	s := NewStackV2()
	s.push(source)

	for !s.isEmpty() {

		// currentNode := fmt.Sprint(s.popFirst()) // - breadth-fist traversal
		currentNode := fmt.Sprint(s.pop()) // - depth-first traversal

		// fmt.Sprint(s.popFirst()) - the only(?) downside of stack of interface{}

		result = append(result, currentNode)

		if len(g[currentNode]) != 0 {
			for _, v := range g[currentNode] {
				s.push(v)
			}
		}

	}

	return result
}

// How nice would it be to make traversal
// to call a user-function for every element
// like filepath.WalkDir !
// traversal(func(){})

func (g graph) traverse(source string, f func(currentNode string)) {

	s := NewStackV2()
	s.push(source)

	for !s.isEmpty() {
		currentNode := fmt.Sprint(s.pop()) // - depth-first traversal

		f(currentNode)
		for _, v := range g[currentNode] {
			s.push(v)
		}
	}
}

type anyGraph map[interface{}][]interface{}

func (g anyGraph) traverse(source interface{}, f func(currentNode interface{})) {

	s := NewStackV2()
	s.push(source)

	for !s.isEmpty() {
		currentNode := s.pop() // - depth-first traversal

		f(currentNode)
		for _, v := range g[currentNode] {
			s.push(v)
		}
	}
}

func Test_Part1(t *testing.T) {

	g := graph{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"e"},
		"d": {"f"},
		"e": {},
		"f": {},
	}

	assert.Equal(t, []string{"b", "c"}, g["a"])

	assert.Equal(t, []string{"a", "c", "e", "b", "d", "f"}, g.depthFirstTraversal("a"))

	g.traverse("a", func(currentNode string) {
		fmt.Println(currentNode)
	}) // This is awesome

	var result []string
	g.traverse("a", func(currentNode string) {
		result = append(result, currentNode)
	}) // This is awesome
	assert.Equal(t, []string{"a", "c", "e", "b", "d", "f"}, result)

	intG := anyGraph{
		2: {3, 6},
		3: {1, 4},
		1: {5},
		5: {7},
	}

	var intResult int
	intG.traverse(2, func(currentNode interface{}) {
		intVal, ok := currentNode.(int) // I know I put a graph of ints - intG, so I use currentValue as int
		if ok {
			intResult += intVal
		}
	}) // This is awesome
	assert.Equal(t, 28, intResult) // silly sum traversal

	v := struct {
		a int
		b string
	}{a: 3, b: "asd"}

	someG := anyGraph{
		"a": {v, v, v},
	}
	someG.traverse("a", func(currentNode interface{}) {
		fmt.Println(currentNode)
	})

}

// Problem 2 - has path
// is there a path from source to destination in a acyclic graph
// I'll utilize a .traverse method of anyGraph type
func (g anyGraph) hasPath(source, destination string) bool {
	var result bool = false
	g.traverse(source, func(currentNode interface{}) {
		strVal, ok := currentNode.(string)
		if ok && strVal == destination {
			result = true
		}
	})
	return result
}

func Test_hasPath(t *testing.T) {

	pathGraph := anyGraph{
		"f": {"g", "i"},
		"g": {"h"},
		"h": {},
		"i": {"g", "k"},
		"j": {"i"},
		"k": {},
	}

	assert.Equal(t, true, pathGraph.hasPath("f", "k"))
	assert.Equal(t, false, pathGraph.hasPath("f", "j"))
}

// How could polymorphic graph of anything be implemented ??
