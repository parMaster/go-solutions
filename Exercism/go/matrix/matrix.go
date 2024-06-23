package matrix

import (
	"errors"
	"strconv"
	"strings"
)

// Define the Matrix type here.
type Matrix [][]int

func New(s string) (Matrix, error) {
	m := Matrix{}

	if len(s) == 0 {
		return m, errors.New("empty input")
	}

	lastLen := 0
	for _, line := range strings.Split(s, "\n") {
		if len(line) == 0 {
			return m, errors.New("empty line")
		}
		row := []int{}
		intStrs := strings.Split(strings.TrimSpace(line), " ")
		for _, c := range intStrs {
			if lastLen != 0 && lastLen != len(intStrs) {
				return m, errors.New("uneven rows")
			}
			lastLen = len(intStrs)
			val, err := strconv.Atoi(c)
			if err != nil {
				return m, errors.New("not a number")
			}
			row = append(row, val)
		}
		m = append(m, row)
	}

	return m, nil
}

// Cols and Rows must return the results without affecting the matrix.
func (m Matrix) Cols() [][]int {
	res := [][]int{}

	for x := 0; x < len(m[0]); x++ {
		row := []int{}
		for y := 0; y < len(m); y++ {
			row = append(row, m[y][x])
		}
		res = append(res, row)
	}

	return res
}

func (m Matrix) Rows() [][]int {
	res := [][]int{}

	for y := range m {
		row := []int{}
		row = append(row, m[y]...)
		res = append(res, row)
	}

	return res
}

func (m Matrix) Set(row, col, val int) bool {

	if len(m) == 0 || row < 0 || row > len(m)-1 || col < 0 || col > len(m[0])-1 {
		return false
	}

	m[row][col] = val

	return true
}
