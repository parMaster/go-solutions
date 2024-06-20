package strain

// Implement the "Keep" and "Discard" function in this file.

// You will need typed parameters (aka "Generics") to solve this exercise.
// They are not part of the Exercism syllabus yet but you can learn about
// them here: https://go.dev/tour/generics/1

func Keep[T any](input []T, testFunc func(val T) bool) []T {
	output := []T{}

	for _, val := range input {
		if testFunc(val) {
			output = append(output, val)
		}
	}

	return output
}

func Discard[T any](input []T, testFunc func(val T) bool) []T {
	output := []T{}

	for _, val := range input {
		if !testFunc(val) {
			output = append(output, val)
		}
	}

	return output
}
