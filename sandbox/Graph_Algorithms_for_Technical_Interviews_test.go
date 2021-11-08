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
		if !history[currentNode.(string)] {
			f(currentNode)
			history[currentNode.(string)] = true
		}
		for _, v := range g[currentNode] {
			_, ok := history[v.(string)]
			if !ok {
				s.push(v)
			}
		}
	}
}

// undirected path returns true if there is a path from to
// note: traverses all the nodes always, even if the path is found right away
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

// Flattens the graph to a map
func (g anyGraph) nodesList() map[string]bool {
	result := make(map[string]bool)
	for k, v := range g {
		result[k.(string)] = false
		if len(v) != 0 {
			for _, vv := range v {
				result[vv.(string)] = false
			}
		}
	}
	return result
}

// Problem 4. Connected components count
func (g anyGraph) countComponents() int {
	flatGraph := g.nodesList()

	result := 0
	for node, visited := range flatGraph {
		if !visited {
			result++
			g.traverseOnce(node, func(currentNode interface{}) {
				flatGraph[currentNode.(string)] = true
			})
		}
	}
	return result
}

func Test_countComponents(t *testing.T) {

	g := edgesListGraph{
		{"a", "b"},

		{"c", "d"},
		{"c", "e"},
		{"c", "f"},
		{"c", "g"},

		{"o", "o"},
	}

	g.asAnyGraph().traverseOnce("o", func(currentNode interface{}) {
		fmt.Println(currentNode.(string))
	})

	g.asAnyGraph().traverseOnce("a", func(currentNode interface{}) {
		fmt.Println(currentNode.(string))
	})

	g.asAnyGraph().traverseOnce("d", func(currentNode interface{}) {
		fmt.Println(currentNode.(string))
	})

	assert.Equal(t, 3, g.asAnyGraph().countComponents())
}

// Problem 5. Largest component
func (g anyGraph) largestComponent() int {
	flatGraph := g.nodesList()

	var largest int = 0
	for node, visited := range flatGraph {
		if !visited {
			result := 0
			g.traverseOnce(node, func(currentNode interface{}) {
				result++
				flatGraph[currentNode.(string)] = true
			})
			largest = max(largest, result)
		}
	}
	return largest
}

func Test_largestComponent(t *testing.T) {

	twoComponentGraph := anyGraph{
		"0": {"1", "5", "8"},
		"1": {"0"},
		"5": {"0", "8"},
		"8": {"0", "5"},

		"4": {"2", "3"},
		"2": {"3", "4"},
		"3": {"2", "4"},
	}

	twoComponentGraph.traverseOnce("8", func(currentNode interface{}) {
		fmt.Print(currentNode.(string))
	})

	assert.Equal(t, 4, twoComponentGraph.largestComponent())
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func (g anyGraph) flatten() []string {
	var flatGraph []string
	g.traverseOnce("w", func(currentNode interface{}) {
		flatGraph = append(flatGraph, currentNode.(string))
	})
	return flatGraph
}

// Problem 6. Shortest path
// from scratch, not using traverseOnce etc.
func (g anyGraph) shortestPath(from, to string) int {

	type queueElement struct {
		node   string
		lenght int
	}

	history := map[string]bool{from: true}

	s := NewStackV2()
	s.push(queueElement{from, 0})
	for !s.isEmpty() {

		current := s.pop().(queueElement)

		history[current.node] = true

		if current.node == to {
			return current.lenght
		}

		for _, neighbour := range g[current.node] {

			_, visited := history[neighbour.(string)]
			if visited {
				continue
			}

			s.push(queueElement{neighbour.(string), current.lenght + 1})
		}
	}
	return -1
}

func Test_shortestPath(t *testing.T) {

	g := edgesListGraph{
		{"w", "x"},
		{"x", "y"},
		{"y", "z"},
		{"z", "v"},
		{"w", "v"},
	}

	flatGraph := g.asAnyGraph().flatten()

	assert.Equal(t, []string{"w", "x", "v", "y", "z"}, flatGraph)

	assert.Equal(t, 2, g.asAnyGraph().shortestPath("w", "z"))
	assert.Equal(t, 2, g.asAnyGraph().shortestPath("x", "z"))
	assert.Equal(t, 0, g.asAnyGraph().shortestPath("z", "z"))
	assert.Equal(t, 2, g.asAnyGraph().shortestPath("z", "w"))
	assert.Equal(t, -1, g.asAnyGraph().shortestPath("zzz", "w"))

}

type gridGraph [][]string

// Problem 7. Island count
// connectedComponent analogue
func (g gridGraph) islandCount() int {

	var result int
	visited := make([][]bool, len(g))

	for i, row := range g {
		visited[i] = make([]bool, len(row))
	}

	for i, row := range g {
		for j, _ := range row {
			if !visited[i][j] && g[i][j] == "L" {
				// new island found
				result++
				g.explore(i, j, visited)
			}

		}
	}

	return result
}

func (g gridGraph) explore(i, j int, visited [][]bool) bool {

	if 0 > i || i >= len(g) ||
		0 > j || j >= len(g[0]) {
		return false
	}

	if g[i][j] == "W" {
		return false
	}

	if visited[i][j] {
		return false
	}

	visited[i][j] = true

	g.explore(i+1, j, visited)
	g.explore(i-1, j, visited)
	g.explore(i, j+1, visited)
	g.explore(i, j-1, visited)
	return true
}

func Test_landCount(t *testing.T) {

	testGrid := gridGraph{
		{"W", "L", "W", "W", "W"},
		{"W", "L", "W", "W", "W"},
		{"W", "W", "W", "L", "W"},
		{"W", "W", "L", "L", "W"},
		{"L", "W", "W", "L", "L"},
		{"L", "L", "W", "W", "W"},
	}
	assert.Equal(t, 3, testGrid.islandCount())

	test2 := gridGraph{
		{"L", "W", "W", "L", "W"},
		{"L", "W", "W", "L", "L"},
		{"W", "L", "W", "L", "W"},
		{"W", "W", "W", "W", "W"},
		{"W", "W", "L", "L", "L"},
	}
	assert.Equal(t, 4, test2.islandCount())

	test3 := gridGraph{
		{"L", "L", "L"},
		{"L", "L", "L"},
		{"L", "L", "L"},
	}
	assert.Equal(t, 1, test3.islandCount())

	test4 := gridGraph{
		{"W", "W"},
		{"W", "W"},
		{"W", "W"},
	}
	assert.Equal(t, 0, test4.islandCount())

}
