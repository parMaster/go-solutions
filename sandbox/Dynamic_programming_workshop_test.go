package sandbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PROBLEM I
// Fibonacci recursion
// Intuitive method, written "from memory" actually
func FibR(n1, n2, i, n int) int {
	if i == n {
		return n2
	}
	return FibR(n2, n1+n2, i+1, n)
}

func Test_Fib_recursion(t *testing.T) {
	fmt.Println(FibR(1, 1, 3, 4))
	assert.Equal(t, 12586269025, FibR(1, 1, 2, 50))
	assert.Equal(t, 3, FibR(1, 2, 3, 4))
}

//  Traditional Fib O(2^N)
func Fib(n int) int {
	if n <= 2 {
		return 1
	}
	return Fib(n-1) + Fib(n-2)
}

func Test_TraditionalFib(t *testing.T) {
	assert.Equal(t, 1, Fib(1))
	assert.Equal(t, 1, Fib(2))
	assert.Equal(t, 2, Fib(3))
	assert.Equal(t, 3, Fib(4))
	assert.Equal(t, 5, Fib(5))
	assert.Equal(t, 8, Fib(6))
	assert.Equal(t, 13, Fib(7))
	assert.Equal(t, 21, Fib(8))
	assert.Equal(t, 34, Fib(9))
	assert.Equal(t, 55, Fib(10))
	assert.Equal(t, 89, Fib(11))

	assert.Equal(t, Fib(45), 1134903170) // 5 seconds execution
	//assert.Equal(t, Fib(50), 12586269025) // timeout

	assert.Equal(t, 12586269025, FibR(1, 1, 2, 50))
}

// Memoized Fib
// OOP Golang(?) way
type Memo struct {
	memo map[int]int
}

func (m *Memo) FibM(n int) int {
	elem, ok := m.memo[n]
	if ok {
		return elem
	}
	if n <= 2 {
		return 1
	}
	m.memo[n] = m.FibM(n-1) + m.FibM(n-2)
	return m.memo[n]
}

func Test_MemoizedFib(t *testing.T) {

	var m Memo

	m.memo = map[int]int{}
	// OR
	// m.memo = make(map[int]int)

	// Another way
	// var m *Memo = new(Memo)
	// (*m).memo = make(map[int]int)

	assert.Equal(t, 1, m.FibM(1))
	assert.Equal(t, 1, m.FibM(2))
	assert.Equal(t, 2, m.FibM(3))
	assert.Equal(t, 3, m.FibM(4))
	assert.Equal(t, 5, m.FibM(5))
	assert.Equal(t, 1134903170, m.FibM(45))
	assert.Equal(t, 12586269025, m.FibM(50))
}

// Just to complete the Fib topic
// Fib with closures in Go:

// fibonacci is a function that returns
// a function that returns an int.
func FibC() func() int {
	var f0, f1 int = 0, 1

	return func() int {
		res := f0
		t := f1
		f1 = f0 + f1
		f0 = t
		return res
	}
}

func Test_FibC_Example(t *testing.T) {
	f := FibC()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
	assert.Equal(t, f(), 55)
}

// PROBLEM II
// Grid Traveller Problem
func gridTraveller(i, j int) int {

	if i == 1 && j == 1 {
		return 1
	}

	if i == 0 || j == 0 {
		return 0
	}

	return gridTraveller(i-1, j) + gridTraveller(i, j-1)
}

func Test_GridTraveller(t *testing.T) {
	assert.Equal(t, 1, gridTraveller(1, 1))
	assert.Equal(t, 3, gridTraveller(2, 3))
	assert.Equal(t, 3, gridTraveller(3, 2))
	assert.Equal(t, 6, gridTraveller(3, 3))
	assert.Equal(t, 70, gridTraveller(5, 5))
	assert.Equal(t, 3432, gridTraveller(8, 8))
	// assert.Equal(t, 3432, gridTraveller(8, 8))
	// assert.Equal(t, 2333606220, gridTraveller(18, 18))
	// test timed out after 30s
}

// Same memoization approach
type Grid struct {
	solutions map[int]map[int]int
}

func (g *Grid) gridTravellerM(i, j int) int {
	// Check base cases first
	if i == 1 && j == 1 {
		return 1
	}

	if i == 0 || j == 0 {
		return 0
	}

	// Check memorized solutions
	// Those are diagonally symmetrical
	elem, ok := g.solutions[i][j]
	if ok {
		return elem
	}
	elem, ok = g.solutions[j][i]
	if ok {
		return elem
	}

	// Now this seems ugly AF
	_, ok = g.solutions[i]
	if !ok {
		g.solutions[i] = map[int]int{}
	}

	g.solutions[i][j] = g.gridTravellerM(i-1, j) + g.gridTravellerM(i, j-1)
	return g.solutions[i][j]
}

func Test_GridTravellerMemoized(t *testing.T) {

	g := new(Grid)
	g.solutions = make(map[int]map[int]int)

	assert.Equal(t, 1, g.gridTravellerM(1, 1))
	assert.Equal(t, 3, g.gridTravellerM(2, 3))
	assert.Equal(t, 3, g.gridTravellerM(3, 2))
	assert.Equal(t, 6, g.gridTravellerM(3, 3))
	assert.Equal(t, 70, g.gridTravellerM(5, 5))
	assert.Equal(t, 3432, g.gridTravellerM(8, 8))
	assert.Equal(t, 2333606220, g.gridTravellerM(18, 18))
}

// MEMOIZATION RECIPE

// I. Make it work
// - visualize the problem as a tree
// - implement the tree using recursion
// - test it

// II. Make it efficient
// - add a memo object
// - add a base case to return memo values
// - store return values into the memo

// PROBLEM III
// canSum(targetSum, numbers) that takes in a
// targetSum and an array of numbers as arguments
// returns true if it's possiible to generate targetSum
// using numbers from the array
// We can use an element of the array as many times as needed
// all input numbers >=0

// PROBLEM III
// canSum(targetSum, numbers) that takes in a
// targetSum and an array of numbers as arguments
// returns true if it's possiible to generate targetSum
// using numbers from the array
// We can use an element of the array as many times as needed
// all input numbers >=0

// Brute force solution
func canSum(t int, e []int) bool {

	var res bool = false

	if t == 0 {
		res = true
		return true
	}

	if t < 0 {
		return false
	}

	for _, v := range e {
		res = res || canSum(t-v, e)
	}

	return res
}

// Memoized brute force solution
// No receiver, just a memo map
// since map is a reference type, no * or & needed
func canSumM(target int, e []int, memo map[int]bool) bool {

	elem, ok := memo[target]
	if ok {
		return elem
	}

	if target == 0 {
		return true
	}

	if target < 0 {
		return false
	}

	for _, v := range e {
		memo[target-v] = canSumM(target-v, e, memo)
		if memo[target-v] {
			return true
		}
	}

	memo[target] = false
	return false
}
func Test_canSum(t *testing.T) {

	testPairs := []struct {
		expected  bool
		targetSum int
		numbers   []int
	}{
		{true, 7, []int{2, 3}},
		{true, 7, []int{5, 3, 4, 7}},
		{false, 7, []int{2, 4}},
		{true, 8, []int{2, 3, 5}},
		{false, 300, []int{7, 14}},
	}

	for _, p := range testPairs {
		assert.Equal(t, p.expected, canSum(p.targetSum, p.numbers))
	}

	for _, p := range testPairs {
		// clear the memory map every time
		solutions := make(map[int]bool)
		assert.Equal(t, p.expected, canSumM(p.targetSum, p.numbers, solutions))
	}
}

// Problem IV - howSum
// return any combination of Numbers sum of which is Target

// Brute force implementation
func howSum(target int, numbers []int) []int {
	if target == 0 {
		return []int{}
	}

	if target < 0 {
		return nil
	}

	for _, v := range numbers {
		result := howSum(target-v, numbers)
		if result != nil {
			return append(result, v)
		}
	}

	return nil
}

// Memoized brute force
func howSumM(target int, numbers []int, memo map[int][]int) []int {

	elem, ok := memo[target]
	if ok {
		return elem
	}

	if target == 0 {
		return []int{}
	}

	if target < 0 {
		return nil
	}

	for _, v := range numbers {
		result := howSumM(target-v, numbers, memo)
		memo[target-v] = result
		if result != nil {
			return append(result, v)
		}
	}

	return nil
}

func Test_howSum(t *testing.T) {

	testPairs := []struct {
		expected  []int
		targetSum int
		numbers   []int
	}{
		{[]int{3, 2, 2}, 7, []int{2, 3}},
		{[]int{4, 3}, 7, []int{5, 3, 4, 7}},
		{nil, 7, []int{2, 4}},
		{[]int{2, 2, 2, 2}, 8, []int{2, 3, 5}},
		{nil, 300, []int{7, 14}},
	}

	for _, p := range testPairs {
		assert.Equal(t, p.expected, howSum(p.targetSum, p.numbers))
	}

	for _, p := range testPairs {
		memo := make(map[int][]int)
		assert.Equal(t, p.expected, howSumM(p.targetSum, p.numbers, memo))
	}
}

// Problem V bestSum
// Find a shortest way to get sum of elements equal Target

func bestSum(target int, numbers []int) []int {
	if target == 0 {
		return []int{}
	}
	if target < 0 {
		return nil
	}

	var shortest []int

	for _, n := range numbers {
		remainder := target - n
		remainderCombination := bestSum(remainder, numbers)
		if remainderCombination != nil {
			combination := append(remainderCombination, n)
			if shortest == nil || len(combination) < len(shortest) {
				shortest = combination
			}
		}
	}

	return shortest
}

func bestSumM(target int, numbers []int, memo map[int][]int) []int {

	elem, ok := memo[target]
	if ok {
		return elem
	}

	if target == 0 {
		return []int{}
	}

	if target < 0 {
		return nil
	}

	var shortest []int

	for _, n := range numbers {
		remainder := target - n
		remainderCombination := bestSumM(remainder, numbers, memo)
		if remainderCombination != nil {
			combination := append(remainderCombination, n)
			if shortest == nil || len(combination) < len(shortest) {
				shortest = combination
			}
		}
	}

	memo[target] = shortest
	return shortest
}

func Test_bestSum(t *testing.T) {

	testPairs := []struct {
		expected  []int
		targetSum int
		numbers   []int
	}{
		{[]int{7}, 7, []int{5, 3, 4, 7}},
		{[]int{5, 3}, 8, []int{2, 3, 5}},
		{[]int{4, 4}, 8, []int{1, 4, 5}},
		{[]int{25, 25, 25, 25}, 100, []int{1, 2, 5, 25}}, //timeout
	}

	for _, p := range testPairs {
		if p.targetSum < 50 { // to avoid timeouts
			assert.Equal(t, p.expected, bestSum(p.targetSum, p.numbers))
		}
	}

	for _, p := range testPairs {
		memo := make(map[int][]int)
		assert.Equal(t, p.expected, bestSumM(p.targetSum, p.numbers, memo))
	}
}
