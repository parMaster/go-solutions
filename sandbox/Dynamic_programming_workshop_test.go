package sandbox

// https://www.youtube.com/watch?v=oBt53YbR9Kk
// Dynamic Programming - Learn to Solve Algorithmic Problems & Coding Challenges
// freeCodeCamp.org @ YouTube

import (
	"fmt"
	"strings"
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

// Traditional Fib O(2^N)
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

// Tribonacci
func Trib(n int) int {
	switch n {
	case 1:
		return 0
	case 2:
		return 0
	case 3:
		return 1
	}
	return Trib(n-1) + Trib(n-2) + Trib(n-3)
}

func Test_Tribonacci(t *testing.T) {
	assert.Equal(t, 0, Trib(1))
	assert.Equal(t, 0, Trib(2))
	assert.Equal(t, 1, Trib(3))
	assert.Equal(t, 1, Trib(4))
	assert.Equal(t, 2, Trib(5))
	assert.Equal(t, 4, Trib(6))
	assert.Equal(t, 7, Trib(7))
	assert.Equal(t, 13, Trib(8))
	assert.Equal(t, 24, Trib(9))
	assert.Equal(t, 44, Trib(10))
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

// Problem VI
// canConstruct(abcdef, [ab, abc, cd, def, abcd]) == true
// Can abcdef be constructed from the elements
// canConstruct(skateboard, [bo, rd, ate, t, ska, sk, boar]) == false
func canConstruct(s string, elements []string) bool {
	if len(s) == 0 {
		return true
	}

	for _, v := range elements {
		if strings.HasPrefix(s, v) {
			if canConstruct(strings.Replace(s, v, "", 1), elements) {
				return true
			}
		}
	}

	return false
}

func canConstructM(s string, elements []string, memo map[string]bool) bool {
	elem, ok := memo[s]
	if ok {
		return elem
	}

	if len(s) == 0 {
		return true
	}

	for _, v := range elements {
		if strings.HasPrefix(s, v) {
			res := canConstructM(strings.Replace(s, v, "", 1), elements, memo)
			memo[s] = res
			if res {
				return true
			}
		}
	}

	memo[s] = false
	return false
}

func Test_canConstruct(t *testing.T) {

	type Tests struct {
		expected bool
		s        string
		elements []string
	}

	testPairs := []Tests{
		{true, "abcdef", []string{"ab", "abc", "cd", "def", "abcd"}},
		{false, "skateboard", []string{"bo", "rd", "ate", "t", "ska", "sk", "boar"}},
		{true, "enterapotentpot", []string{"a", "p", "ent", "enter", "ot", "o", "t"}},
	}

	for _, p := range testPairs {
		assert.Equal(t, p.expected, canConstruct(p.s, p.elements))
	}

	hardCase := Tests{
		false, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"},
	}
	testPairs = append(testPairs, hardCase)

	for _, p := range testPairs {
		memo := make(map[string]bool)
		assert.Equal(t, p.expected, canConstructM(p.s, p.elements, memo))
	}
}

// canConstruct notes
// TIME COMPLEXITY is O(n^m):
// n - wordBank.lenght and branching factor, because every word is checked on each level
// n*n*n*n... every time
// m - target length - tree height
// canConstruct will be called n^m times, thus the time complexity is O(n^m)
//
// Complexity contributers:
// Branching factor (N) - how FAST is it growing
// Tree height 		(M) - how TALL it is
// Total_number_of_Iterations = (Branching_factor ^ Tree_height)
//
// Other costly operations:
// Iteratively creating a subslice every call - will make it more expensive, like:
// O(n^m * m)
// m - operations to iteratively create a subslice, every time, so n^m times
//
// SPACE COMPLEXITY is O(m*m) ??
//
// Memoized time complexity
// O(n*m * m)
// O(m^2) space complexity

// Problem VII
// countConstruct (target, wordBank)
// how many ways to construct a target string with wordBank?
// words can be used as many times as needed
// Example :
// target = "purple" , wordBank = {purp, p, ur, le, purpl}
// purp -> le -> "" 1
// p -> urple -> ur -> ple -> le -> "" 1
// ur -> nil
// le - nil
// purpl -> e -> nil
// Answer = 2
func countConstructM(target string, wordBank []string, memo map[string]int) int {
	elem, ok := memo[target]
	if ok {
		return elem
	}

	if len(target) == 0 {
		return 1
	}

	var countConstruct int = 0
	for _, v := range wordBank {
		if strings.HasPrefix(target, v) {
			memo[target] = countConstructM(strings.Replace(target, v, "", 1), wordBank, memo)
			countConstruct += memo[target]
		}
	}

	memo[target] = countConstruct
	return countConstruct
}

func Test_countConstruct(t *testing.T) {

	type Tests struct {
		expected int
		s        string
		elements []string
	}

	testPairs := []Tests{
		{1, "abcdef", []string{"ab", "abc", "cd", "def", "abcd"}},
		{0, "skateboard", []string{"bo", "rd", "ate", "t", "ska", "sk", "boar"}},
		{4, "enterapotentpot", []string{"a", "p", "ent", "enter", "ot", "o", "t"}},
	}

	hardCase := Tests{
		0, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"},
	}
	hardCaseSolvable := Tests{
		73859288608, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee", "t"},
	}
	testPairs = append(testPairs, hardCase)
	testPairs = append(testPairs, hardCaseSolvable)

	for _, p := range testPairs {
		memo := make(map[string]int)
		assert.Equal(t, p.expected, countConstructM(p.s, p.elements, memo))
	}
}

// Problem VIII
// allConstruct (target, wordBank)
// same, but returns 2D array with all answers
// i.e. 	allConstruct("abcdef", []string{"ab", "c", "cd", "def", "ef", "abcd")
// Should return:
// ab, cd, ef
// ab, c, def
// abc, def
// abcd, ef
// 4 in total, as countConstruct counted before

type QItem struct {
	target string
	found  []string
}

func AllConstruct(target string, wordBank []string) (result [][]string) {

	result = make([][]string, 0)
	var queue = []QItem{{target: target, found: []string{}}}

	proceed := true
	for proceed {
		proceed = false
		newQueue := []QItem{}
		for _, qi := range queue {
			if qi.target == "" {
				result = append(result, qi.found)
				continue
			}
			for _, v := range wordBank {
				if strings.HasPrefix(qi.target, v) {
					newQueue = append(newQueue, QItem{target: strings.Replace(qi.target, v, "", 1), found: append(qi.found, v)})
					proceed = true
				}
			}
		}

		queue = make([]QItem, 0)
		for _, qi := range newQueue {
			queue = append(queue, QItem{target: qi.target, found: qi.found})
		}
	}

	return
}

func Test_AllConstruct(t *testing.T) {
	assert.EqualValues(t, [][]string{{"abc", "def"}}, AllConstruct("abcdef", []string{"ab", "abc", "cd", "def", "abcd"}))
	assert.EqualValues(t, [][]string{{"abc", "def"}, {"abcd", "ef"}, {"ab", "cd", "ef"}}, AllConstruct("abcdef", []string{"ab", "abc", "cd", "def", "ef", "abcd"}))
	assert.EqualValues(t, [][]string{{"abcd", "ef"}, {"ab", "c", "def"}, {"ab", "cd", "ef"}}, AllConstruct("abcdef", []string{"ab", "c", "cd", "def", "ef", "abcd"}))
	assert.EqualValues(t, [][]string{{"purp", "le"}, {"p", "ur", "p", "le"}}, AllConstruct("purple", []string{"purp", "p", "ur", "le", "purpl"}))

	assert.EqualValues(t, [][]string{}, AllConstruct("eeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"}))
}

// Recursive AllConstruct solution
func AllConstructRec(target string, wordBank []string, current []string) (result [][]string) {
	result = make([][]string, 0)
	if target == "" {
		return [][]string{current}
	}

	for _, v := range wordBank {
		if strings.HasPrefix(target, v) {

			newTarget := strings.Replace(target, v, "", 1)

			newCurrent := []string{}
			newCurrent = append(newCurrent, current...)

			newCurrent = append(newCurrent, v)

			res := AllConstructRec(newTarget, wordBank, newCurrent)
			result = append(result, res...)
		}
	}

	return
}

func Test_AllConstructRec(t *testing.T) {
	assert.EqualValues(t, [][]string{{"abc", "def"}}, AllConstructRec("abcdef", []string{"ab", "abc", "cd", "def", "abcd"}, []string{}))
	assert.EqualValues(t, [][]string{{"ab", "cd", "ef"}, {"abc", "def"}, {"abcd", "ef"}}, AllConstructRec("abcdef", []string{"ab", "abc", "cd", "def", "ef", "abcd"}, []string{}))
	assert.EqualValues(t, [][]string{}, AllConstructRec("eeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"}, []string{}))
}

// Memoized recursive AllConstruct solution
func AllConstructRecMemoized(target string, wordBank []string, memo map[string][][]string) (result [][]string) {
	elem, ok := memo[target]
	if ok {
		return elem
	}

	result = make([][]string, 0)
	if target == "" {
		return [][]string{{}}
	}

	for _, v := range wordBank {
		if strings.HasPrefix(target, v) {

			newTarget := strings.Replace(target, v, "", 1)

			res := AllConstructRecMemoized(newTarget, wordBank, memo)
			for _, r := range res {
				r = append(r, v)
				result = append(result, r)
			}
		}
	}

	memo[target] = result
	return
}

func Test_AllConstructRecMemoized(t *testing.T) {
	assert.EqualValues(t, [][]string{{"def", "abc"}}, AllConstructRecMemoized("abcdef", []string{"ab", "abc", "cd", "def", "abcd"}, map[string][][]string{}))
	assert.EqualValues(t, [][]string{}, AllConstructRecMemoized("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"}, map[string][][]string{}))
	assert.EqualValues(t, [][]string{{"ef", "cd", "ab"}, {"def", "abc"}, {"ef", "abcd"}}, AllConstructRecMemoized("abcdef", []string{"ab", "abc", "cd", "def", "ef", "abcd"}, map[string][][]string{}))
}

// Problem IX
// Fib Tabulation
// Building a table iteratively
// fib(6) -> 8
// 0	1	2	3	4	5	6
// take add add
// 0	1	1	2	3	5	8
func fibT(n int) int {
	t := make([]int, n+2)

	t[0] = 0
	t[1] = 1

	for i := 0; i < n; i++ {
		t[i+1] += t[i]
		t[i+2] += t[i]
	}

	return t[n]
}

func Test_FibT(t *testing.T) {
	assert.Equal(t, 8, fibT(6))
	assert.Equal(t, 13, fibT(7))
	assert.Equal(t, 21, fibT(8))
	assert.Equal(t, 12586269025, fibT(50))
}

// Problem X
// Tabulated gridTraveller
type TabGrid struct {
	solutions [][]int
}

func (g *TabGrid) gridTraveller(n, m int) int {
	g.solutions = make([][]int, n+2)
	for i := 0; i <= n+1; i++ {
		g.solutions[i] = make([]int, m+2)
	}
	g.solutions[1][1] = 1

	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			g.solutions[i+1][j] += g.solutions[i][j]
			g.solutions[i][j+1] += g.solutions[i][j]
		}
	}

	return g.solutions[n][m]
}

func Test_GridTravellerTabulized(t *testing.T) {

	var g TabGrid
	assert.Equal(t, 1, g.gridTraveller(1, 1))
	assert.Equal(t, 3, g.gridTraveller(2, 3))
	assert.Equal(t, 3, g.gridTraveller(3, 2))
	assert.Equal(t, 6, g.gridTraveller(3, 3))
}

/**
*	Tabulation Recipe
	- No "brute force first" like with memoization.

	- Visualize the problem as a table.
	- Size the table based On the inputs. Watch out for "off by one" errors
	- initialize the table with default values
	- seed the trivial answer into the table (0,1 for Fib)
	- iterate through the table
	- fill further positions based on the current position
*/

/*
  - Problem Xi - canSum tabulation
    canSum(7, [5,3,4]) -> true

    initializing table of size 7+1
    [F, F, F, F, F, F, F, F]

    seed a basic solution canSub(0) always == true
    [T, F, F, F, F, F, F, F]

    iteration from 0 to target, look ahead at every Element array index and mark it True
    0  1  2  3  4  5  6  7
    [T, F, F, T, T, T, F, F]

    Then we'll skip 1 and 2 because it's already false branches, and finally:
    after we chacked the bounds
    [T, F, F, T, T, T, T, T]

    m[7] is the result

    Tabulized solution
    O(M*N) time complexity
    O(M) space somplexity
*/
func canSumTabulized(target int, elements []int) bool {

	m := make([]bool, target+1)

	for i := 0; i <= target; i++ {
		m[i] = false
	}
	m[0] = true

	for i := 0; i <= target; i++ {
		if m[i] == false {
			continue
		}

		for j := 0; j < len(elements); j++ {
			if i+elements[j] <= target {
				m[int(i+elements[j])] = true
			}
		}
	}

	return m[target]
}

func Test_canSumTab(t *testing.T) {
	testPairs := []struct {
		expected  bool
		targetSum int
		numbers   []int
	}{
		{true, 7, []int{2, 3}},
		{true, 7, []int{5, 3, 4}},
		{false, 7, []int{2, 4}},
		{true, 8, []int{2, 3, 5}},
		{false, 300, []int{7, 14}},
	}

	for _, p := range testPairs {
		assert.Equal(t, p.expected, canSumTabulized(p.targetSum, p.numbers))
	}

}

/*
* Problem XII - howSum tabulation

howSum(7, [5,3,4]) -> [4,3]

so, return an array
seed value: howSum(0, [...]) is always an empty array
*/
func howSumTabulized(target int, elements []int) []int {

	m := make([][]int, target+1)

	for i := 0; i <= target; i++ {
		m[i] = nil
	}
	m[0] = []int{}

	for i := 0; i <= target; i++ {

		if m[i] == nil {
			continue
		}

		for _, elem := range elements {

			if i+elem <= target {

				m[i+elem] = append(m[i], elem)
			}
		}
	}

	return m[target]
}

func Test_howSumTabulized(t *testing.T) {

	assert.Equal(t, []int{}, howSumTabulized(0, []int{5, 3, 4}))

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
		assert.Equal(t, p.expected, howSumTabulized(p.targetSum, p.numbers))
	}
}

/*
Problem XIII - bestSum(target, elements{}) tabulation
*/
func bestSumTabulized(target int, elements []int) []int {

	// initialize
	m := make([][]int, target+1)
	for i := 0; i <= target; i++ {
		m[i] = nil
	}

	// seed
	m[0] = []int{}

	for i := 0; i <= target; i++ {
		if m[i] == nil {
			continue
		}

		for _, elem := range elements {

			// bounds check
			if i+elem <= target {

				// Cgecking target function - minimize the route (result length)
				if m[i+elem] == nil || len(append(m[i], elem)) < len(m[i+elem]) {
					m[i+elem] = append(m[i], elem)
				}
			}
		}
	}

	return m[target]
}

func Test_bestSumTabulized(t *testing.T) {

	testPairs := []struct {
		expected  []int
		targetSum int
		numbers   []int
	}{
		{[]int{7}, 7, []int{5, 3, 4, 7}},
		{[]int{3, 5}, 8, []int{2, 3, 5}},
		{[]int{4, 4}, 8, []int{1, 4, 5}},
		{[]int{25, 25, 25, 25}, 100, []int{1, 2, 5, 25}}, //timeout
	}

	for _, p := range testPairs {
		assert.Equal(t, p.expected, bestSumTabulized(p.targetSum, p.numbers))
	}
}

/*
Problem IX - tabulating canConstruct
canConstruct(abcdef, {ab, abc, cd, def, abcd}) -> true
*/
func canConstructTabulized(target string, elements []string) bool {

	m := make([]bool, len(target)+1)
	for i := 0; i <= len(target); i++ {
		m[i] = false
	}

	// m[i] == true - means that it is possible to construct a string target[0:i-1], so
	// true at m[len(target)] means target can be constructed

	// seed value - empty string can always be constructed whatewer the elements array
	m[0] = true
	subTarget := target

	for i := 0; i <= len(target); i++ {

		if !m[i] {
			continue
		}

		subTarget = target[i:]

		for _, elem := range elements {
			if strings.HasPrefix(subTarget, elem) {
				m[i+len(elem)] = true
			}
		}
	}

	return m[len(target)]
}

func Test_canConstructTabulized(t *testing.T) {

	type Tests struct {
		expected bool
		s        string
		elements []string
	}

	testPairs := []Tests{
		{true, "abcdef", []string{"ab", "abc", "cd", "def", "abcd"}},
		{false, "skateboard", []string{"bo", "rd", "ate", "t", "ska", "sk", "boar"}},
		{true, "enterapotentpot", []string{"a", "p", "ent", "enter", "ot", "o", "t"}},
	}

	hardCase := Tests{
		false, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"},
	}
	testPairs = append(testPairs, hardCase)

	for _, p := range testPairs {
		assert.Equal(t, p.expected, canConstructTabulized(p.s, p.elements))
	}
}

/*
Problem X - countConstruct Tabulation
countConstruct( purple , { purp, p, ur, le, purpl }) -> 2
*/
func countConstructTab(target string, elements []string) int {

	m := make([]int, len(target)+1)
	for i := 0; i <= len(target); i++ {
		m[i] = 0
	}

	// m[i] == 2 - means that there are 2 ways to construct a string target[0:i-1], so
	// number at m[len(target)] is a solution for _target_ string

	// seed value - empty string can always be constructed whatewer the elements array
	m[0] = 1
	subTarget := target

	for i := 0; i <= len(target); i++ {

		if m[i] == 0 {
			continue
		}

		subTarget = target[i:]

		for _, elem := range elements {
			if strings.HasPrefix(subTarget, elem) {
				m[i+len(elem)] += m[i]
			}
		}
	}

	return m[len(target)]
}

func Test_countConstructTab(t *testing.T) {

	type Tests struct {
		expected int
		s        string
		elements []string
	}

	testPairs := []Tests{
		{1, "abcdef", []string{"ab", "abc", "cd", "def", "abcd"}},
		{0, "skateboard", []string{"bo", "rd", "ate", "t", "ska", "sk", "boar"}},
		{4, "enterapotentpot", []string{"a", "p", "ent", "enter", "ot", "o", "t"}},
	}

	hardCase := Tests{
		0, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"},
	}
	hardCaseSolvable := Tests{
		73859288608, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee", "t"},
	}
	testPairs = append(testPairs, hardCase)
	testPairs = append(testPairs, hardCaseSolvable)

	for _, p := range testPairs {
		assert.Equal(t, p.expected, countConstructTab(p.s, p.elements))
	}
}

/*
Problem XII - allConstruct Tabulation

allConstruct("abcdef", []string{"ab", "abc", "cd", "def", "abcd")
Should return:
[

	[ab, cd, ef]
	[ab, c, def]
	[abc, def]
	[abcd, ef]

]
*/
func allConstructTab(target string, elements []string) [][]string {

	m := make([][][]string, len(target)+1)

	for i := 0; i <= len(target); i++ {
		m[i] = [][]string{}
	}
	m[0] = [][]string{{}}

	for i := 0; i <= len(target); i++ {

		if len(m[i]) == 0 {
			continue
		}

		subTarget := target[i:]

		for _, elem := range elements {

			if strings.HasPrefix(subTarget, elem) {

				futurePos := i + len(elem)
				for _, v := range m[i] {
					m[futurePos] = append(m[futurePos], append(v, elem))
				}
			}
		}
	}

	return m[len(target)]
}

func Test_allConstructTab(t *testing.T) {

	type Tests struct {
		expected [][]string
		s        string
		elements []string
	}

	testPairs := []Tests{
		{[][]string{{"abc", "def"}, {"ab", "c", "def"}, {"abcd", "ef"}, {"ab", "cd", "ef"}}, "abcdef", []string{"ab", "abc", "cd", "def", "abcd", "ef", "c"}},
		{[][]string{{}}, "", []string{"cat", "dog"}},
	}

	// Timeout
	// Tabulazed solution is EXPONENTIAL actually

	// hardCase := Tests{
	// 	[][]string{{}}, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeet", []string{"e", "ee", "eee", "eeee"},
	// }
	// testPairs = append(testPairs, hardCase)

	for _, p := range testPairs {
		assert.Equal(t, p.expected, allConstructTab(p.s, p.elements))
	}
}
