package main

// https://www.hackerrank.com/challenges/coin-change/problem
func getWays(target int, numbers []int, memo map[int][]int) int {
	if target == 0 {
		return 1
	}

	ways := 0
	for _, v := range numbers {
		if v > target {
			continue
		}
		ways += getWays(target-v, numbers, memo)
	}

	return ways
}

func Test_getWays(t *testing.T) {

	testPairs := []struct {
		expected  int
		targetSum int
		numbers   []int
	}{
		{3, 3, []int{8, 3, 1, 2}},
		{4, 4, []int{1, 2, 3}},
		{5, 10, []int{2, 5, 3, 6}},
	}

	for _, p := range testPairs {
		memo := make(map[int][]int)
		assert.Equal(t, p.expected, getWays(p.targetSum, p.numbers, memo))
	}

}

// func getWays(n int, c []int) int {

// }
