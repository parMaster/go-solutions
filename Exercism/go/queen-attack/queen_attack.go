package queenattack

import (
	"errors"
	"math"
)

func CanQueenAttack(wp, bp string) (bool, error) {

	if len(wp) != 2 || len(bp) != 2 ||
		wp == bp ||
		wp[0] < 'a' || bp[0] < 'a' ||
		wp[0] > 'h' || bp[0] > 'h' ||
		wp[1]-0x30 > 8 || bp[1]-0x30 > 8 ||
		wp[1]-0x30 < 1 || bp[1]-0x30 < 1 {
		return false, errors.New("invalid input")
	}

	if wp[0] == bp[0] || wp[1] == bp[1] || // same col/row
		math.Abs(float64(int(wp[0])-int(bp[0]))) == math.Abs(float64(int(wp[1])-int(bp[1]))) { // share diagonal
		return true, nil
	}

	return false, nil
}
