package allyourbase

import (
	"errors"
	"math"
	"slices"
)

func ConvertToBase(inputBase int, inputDigits []int, outputBase int) ([]int, error) {

	if inputBase < 2 {
		return nil, errors.New("input base must be >= 2")
	}
	if outputBase < 2 {
		return nil, errors.New("output base must be >= 2")
	}

	slices.Reverse(inputDigits)

	val := 0
	for i, d := range inputDigits {
		if d < 0 || d >= inputBase {
			return nil, errors.New("all digits must satisfy 0 <= d < input base")
		}
		val += int(float64(d) * math.Pow(float64(inputBase), float64(i)))
	}

	if val == 0 {
		return []int{0}, nil
	}

	output := []int{}
	for val > 0 {
		reminder := val % outputBase
		output = append(output, reminder)
		val /= outputBase
	}
	slices.Reverse(output)

	return output, nil
	// panic("Please implement the ConvertToBase function")
}
