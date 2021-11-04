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
type strGraph map[string][]string

// OR DO AN POLYMORPHISM kind of thing -  graph can be string graph or int graph

// Problem 1. Depth-First Traversal (breadth-first as well)

// Iterative, intuitive method
// in Value Semantics
func (g strGraph) depthFirstTraversal(source string) []string {
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

func (g strGraph) traverse(source string, f func(currentNode string)) {

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

	g := strGraph{
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
// is there a path from source to destination in an acyclic graph
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
/*
type graph interface {
	iStrGraph
	iNumGraph
}

type iStrGraph interface {
	traverse(param string)
}
type iNumGraph interface {
	traverse(param int)
}

// doesn't work this way - method signatures are different
// how is it better than using interface{} type?

type lettersGraph map[string][]string
type numbersGraph map[int][]int

func (l *lettersGraph) traverse(param string) {
	fmt.Println("traversing letters graph")
}

func (n *numbersGraph) traverse(param int) {
	fmt.Println("traversing numbers graph")
}

func traverse(g graph, param interface{}) {
	g.traverse(param)
}

func Test_traverseInterfaces(t *testing.T) {

	var numbers graph = &numbersGraph{}
	var letters graph = &lettersGraph{}

	numbers.traverse()
	letters.traverse()
}

// "Surprisingly", animal-cat-dog example didn't help at all
*/

// Problem 3. Undirected path
// Finding path in undirected graph. Every edge is bi-directional
// Mind the cycles

// array of edges to adjacency list conversion

type edge struct {
	from string
	to   string
}

type edgesListGraph []edge

func (g edgesListGraph) asAnyGraph() anyGraph { // love this

	a := anyGraph{}
	for _, v := range g {
		_, ok := a[v.from]
		if ok {
			a[v.from] = append(a[v.from], v.to)
		} else {
			a[v.from] = []interface{}{v.to}
		}

		_, ok = a[v.to]
		if ok {
			a[v.to] = append(a[v.to], v.from)
		} else {
			a[v.to] = []interface{}{v.from}
		}
	}

	return a
}

func Test_conversions(t *testing.T) {

	g := edgesListGraph{
		edge{"i", "j"},
		edge{"k", "i"},
		edge{"k", "j"},
		edge{"m", "k"},
		edge{"k", "l"},
		edge{"o", "n"},
	}

	ag := anyGraph{
		"i": {"j", "k"},
		"j": {"i", "k"},
		"k": {"i", "j", "m", "l"},
		"m": {"k"},
		"l": {"k"},

		"o": {"n"},
		"n": {"o"},
	}

	assert.Equal(t, ag, g.asAnyGraph())

	assert.Equal(t, true, g.undirectedPath("j", "m"))
	assert.Equal(t, true, g.undirectedPath("m", "m"))
	assert.Equal(t, true, g.undirectedPath("o", "n"))
	assert.Equal(t, true, g.undirectedPath("n", "o"))

	assert.Equal(t, false, g.undirectedPath("n", "m"))
	assert.Equal(t, false, g.undirectedPath("m", "n"))
}

// traversing graph once, keeps history of visited nodes and minds the cycles
func (g anyGraph) traverseOnce(source interface{}, f func(currentNode interface{})) {

	s := NewStackV2()
	s.push(source)

	history := make(map[string]bool)

	for !s.isEmpty() {
		currentNode := s.popFirst()
		current := currentNode.(string) // assume there is a string. Will cause panic if there's something else!
		f(currentNode)
		for _, v := range g[currentNode] {
			_, ok := history[v.(string)] // same
			if !ok {
				s.push(v)

				history[current] = true
			}
		}

	}
}

func (g *edgesListGraph) undirectedPath(from, to string) bool {

	var result bool

	g.asAnyGraph().traverseOnce(from, func(currentNode interface{}) {
		current, ok := currentNode.(string)
		if ok && to == current {
			result = true
		}
	})
	return result
}
