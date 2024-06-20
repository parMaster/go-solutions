package perfect

import (
	"errors"
)

var ErrOnlyPositive = errors.New("only positive numbers allowed")

type Classification = byte

const (
	ClassificationDeficient Classification = iota + 1
	ClassificationPerfect
	ClassificationAbundant
)

// Perfect: aliquot sum = number
//
//	6 is a perfect number because (1 + 2 + 3) = 6
//	28 is a perfect number because (1 + 2 + 4 + 7 + 14) = 28
//
// Abundant: aliquot sum > number
//
//	12 is an abundant number because (1 + 2 + 3 + 4 + 6) = 16
//	24 is an abundant number because (1 + 2 + 3 + 4 + 6 + 8 + 12) = 36
//
// Deficient: aliquot sum < number
//
//	8 is a deficient number because (1 + 2 + 4) = 7
//	Prime numbers are deficient
func Classify(n int64) (Classification, error) {
	if n <= 0 {
		return 0, ErrOnlyPositive
	}

	aSum := aliquotSum(n)

	if aSum == n {
		return ClassificationPerfect, nil
	}

	if aSum > n {
		return ClassificationAbundant, nil
	}

	if aSum < n {
		return ClassificationDeficient, nil
	}

	return 0, nil
}

// The aliquot sum is defined as the sum of the factors of a number not including the number itself.
// For example, the aliquot sum of 15 is 1 + 3 + 5 = 9.
func aliquotSum(n int64) int64 {
	var sum int64
	for i := int64(1); i < n; i++ {
		if n%i == 0 {
			sum += i
		}
	}
	return sum
}
