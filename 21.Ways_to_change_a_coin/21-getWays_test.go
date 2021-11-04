// https://www.hackerrank.com/challenges/coin-change/problem
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// getWays from HackerRank
func getWays(n int32, c []int64) int64 {

	dp := make([]int64, n+1)
	dp[0] = 1
	for i, _ := range c {
		for j := c[i]; j <= int64(n); j++ {
			dp[j] += dp[j-c[i]]
		}
	}
	return dp[n]
}

func Test_getWays(t *testing.T) {

	testPairs := []struct {
		expected  int64
		targetSum int32
		numbers   []int64
	}{
		{int64(3), 3, []int64{8, 3, 1, 2}},
		{int64(4), 4, []int64{1, 2, 3}},
		{int64(5), 10, []int64{2, 5, 3, 6}},
		{int64(96190959), 166, []int64{5, 37, 8, 39, 33, 17, 22, 32, 13, 7, 10, 35, 40, 2, 43, 49, 46, 19, 41, 1, 12, 11, 28}},
	}

	for _, p := range testPairs {
		// memo := make(map[int][]int)
		assert.Equal(t, p.expected, getWays(p.targetSum, p.numbers))
	}

}

/*

 0  1  2  3  4  5  6  7  8  9  10
init
[1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]

i == 1, coins[i] == 2
[1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0] j == 2..n; dp[2] += dp[0] // so 1
[1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0] dp[3] += dp[1] // so 0
[1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0] // 1
[1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0]
[1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0]
[1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0]
[1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0]
[1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0]
[1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1]
if there is just one coin = 2 (coins is {2}) , there is 1 way to change 2,4,6,8,10 Euro coins - 1*2, 2*2, 3*2, 4*2, or 5*2

i==2, coins[2] == 5
[1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1] start with 5, dp[5] += dp[5-5]
[1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1] dp[6] += dp[6-5] // 1+=0
[1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1] dp[7] += dp[7-5] // 0+=1 = 1
[1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1] dp[8]
[1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1] dp[9]
[1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 2] dp[10] += dp[10-5] // 1+=1 = 2
i==3, coins[3] == 3
[1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2] dp[3] += dp[3-3] // 0+1
[1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 2]
[1, 0, 1, 1, 1, 2, 1, 1, 1, 1, 2]
[1, 0, 1, 1, 1, 2, 2, 1, 1, 1, 2]
[1, 0, 1, 1, 1, 2, 2, 2, 1, 1, 2]
[1, 0, 1, 1, 1, 2, 2, 2, 3, 1, 2]
[1, 0, 1, 1, 1, 2, 2, 2, 3, 3, 2]
[1, 0, 1, 1, 1, 2, 2, 2, 3, 3, 4]
i==4, coins[4] == 6
[1, 0, 1, 1, 1, 2, 3, 2, 3, 3, 4] dp[6] = dp[6] + dp[6-6] = 2+1
[1, 0, 1, 1, 1, 2, 3, 2, 3, 3, 4]
[1, 0, 1, 1, 1, 2, 3, 2, 4, 3, 4]
[1, 0, 1, 1, 1, 2, 3, 2, 4, 4, 4]
[1, 0, 1, 1, 1, 2, 3, 2, 4, 4, 5] dp[10] = dp[10] + dp[10-6] = dp[10] + dp[4] = 4+1

*/

// works perfectly at fast machines with long stacks
// but results in half of the tests being failed due to time limit
func getWays_recursive(n int32, c []int64) int64 {

	if n == 0 {
		return 1
	}

	if n < 0 {
		return 0
	}

	var ways int64 = 0
	for i, v := range c {
		ways += getWays(n-int32(v), c[i:])
	}

	return ways
}
