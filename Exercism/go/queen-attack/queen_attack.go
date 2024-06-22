package queenattack

import (
	"errors"
	"math"
)

func CanQueenAttack(whitePosition, blackPosition string) (bool, error) {
	// invalid input
	if len(whitePosition) != 2 || len(blackPosition) != 2 ||
		whitePosition == blackPosition ||
		int(whitePosition[0])-97 > 7 || int(blackPosition[0])-97 > 7 ||
		int(whitePosition[1])-48 > 8 || int(blackPosition[1])-48 > 8 ||
		int(whitePosition[1])-48 < 1 || int(blackPosition[1])-48 < 1 {
		return false, errors.New("invalid inpit")
	}

	p1 := [2]int{
		int(whitePosition[0]) - 97,
		int(whitePosition[1]),
	}

	p2 := [2]int{
		int(blackPosition[0]) - 97,
		int(blackPosition[1]),
	}

	// on the same row/col
	if p1[0] == p2[0] || p1[1] == p2[1] {
		return true, nil
	}

	if math.Abs(float64(p1[0]-p2[0])) == math.Abs(float64(p1[1]-p2[1])) {
		return true, nil
	}

	return false, nil
}
