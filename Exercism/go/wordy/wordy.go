package wordy

import (
	"strconv"
	"strings"
)

//  0. What is 5?
//  1. What is 5 plus 13?
//  2. What is 7 minus 5?
//     What is 6 multiplied by 4?
//     What is 25 divided by 5?
//  3. What is 5 plus 13 plus 6? - no priority, always left-to right
//  4. What is 3 plus 2 multiplied by 3? (15, not 9)
func Answer(question string) (int, bool) {
	// question = "What is 1 plus 1 plus 1?"
	// cleanup input:
	q, _ := strings.CutPrefix(question, "What is ")
	q = strings.Trim(q, "?")
	q = strings.Replace(q, "by ", "", -1) // so operations are always 1 word, and:
	qParts := strings.Split(q, " ")
	var val int
	var err error
	if len(qParts)%2 != 1 {
		return 0, false
	} else {
		// n operations + (n-1) operands
		// first number
		val, err = strconv.Atoi(qParts[0])
		if err != nil {
			return 0, false
		}

		for i := 1; i < len(qParts)-1; i++ {
			// operand parsing
			val2, err := strconv.Atoi(qParts[i+1])
			if err != nil {
				return 0, false
			}

			switch qParts[i] {
			case "plus":
				val += val2
			case "minus":
				val -= val2
			case "multiplied":
				val *= val2
			case "divided":
				val /= val2
			default:
				return 0, false
				// unsopported operation
			}
			i++
		}
	}

	return val, true
}
