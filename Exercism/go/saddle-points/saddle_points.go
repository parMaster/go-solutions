package matrix

import (
	"strconv"
	"strings"
)

type Matrix [][]int

type Pair [2]int

// apparently, we're not allowed to return error, even if the input is empty )))
// absolutely pathetic
func New(s string) (*Matrix, error) {
	m := Matrix{}
	if len(s) == 0 {
		return &m, nil
	}
	lines := strings.Split(strings.TrimSpace(s), "\n")

	for _, line := range lines {
		row := []int{}
		lineValues := strings.Split(strings.TrimSpace(line), " ")
		for _, c := range lineValues {
			v, err := strconv.Atoi(string(c))
			if err != nil {
				panic(err)
			}
			row = append(row, v)
		}
		m = append(m, row)
	}
	return &m, nil
}

func (m Matrix) checkSaddle(y, x int) bool {

	// check north to south
	for i := 0; i < len(m); i++ {
		if m[i][x] < m[y][x] && i != y {
			return false
		}
	}

	// check west to east
	for i := 0; i < len(m[0]); i++ {
		if m[y][i] > m[y][x] && i != x {
			return false
		}
	}

	return true
}

func (m *Matrix) Saddle() []Pair {

	res := []Pair{}

	for i := range *m {
		for j := range (*m)[i] {

			if m.checkSaddle(i, j) {
				res = append(res, [2]int{i + 1, j + 1})
			}

		}
	}

	return res
}
