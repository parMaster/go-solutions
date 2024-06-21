package brackets

// closing brackest to opening brackets pairs
var closingOpeningPairs map[rune]rune = map[rune]rune{
	']': '[',
	')': '(',
	'}': '{',
}

func Bracket(input string) bool {
	// brackets stack
	var bStack string

	for _, c := range input {
		switch c {
		case '[', '(', '{':
			bStack += string(c)
		case ']', ')', '}':
			if len(bStack) == 0 {
				return false
			}
			if rune(bStack[len(bStack)-1]) == closingOpeningPairs[c] {
				bStack = bStack[:len(bStack)-1]
			} else {
				return false
			}
		}
	}

	return len(bStack) == 0
}
