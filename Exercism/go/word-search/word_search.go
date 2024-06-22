package wordsearch

import "errors"

var dir = [8][2]int{
	{0, 1},   // right
	{1, 1},   // down right
	{1, 0},   // down
	{1, -1},  // down left
	{0, -1},  // left
	{-1, -1}, // up left
	{-1, 0},  // up
	{-1, 1},  // up right
}

func Solve(words []string, puzzle []string) (map[string][2][2]int, error) {
	result := map[string][2][2]int{}

	foundEverything := true
	for _, word := range words {
		found := false
		for i, row := range puzzle {
			for j := range row {

				for _, d := range dir {

					for k, c := range word {
						if i+k*d[0] < 0 || i+k*d[0] >= len(puzzle) || j+k*d[1] < 0 || j+k*d[1] >= len(puzzle[i]) {
							break
						}

						if puzzle[i+k*d[0]][j+k*d[1]] != byte(c) {
							break
						}

						if k == len(word)-1 {
							found = true
							result[word] = [2][2]int{{j, i}, {j + k*d[1], i + k*d[0]}}
						}

					}

				}

			}
		}
		if !found {
			foundEverything = false
			result[word] = [2][2]int{{-1, -1}, {-1, -1}}
		}
	}

	if !foundEverything {
		return result, errors.New("no words found")
	}

	return result, nil
}
