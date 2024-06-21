package encode

import (
	"fmt"
	"strconv"
)

func RunLengthEncode(input string) string {

	if len(input) == 0 {
		return ""
	}

	var last rune = rune(input[0])
	cnt := 1

	enc := ""
	for _, c := range input[1:] {
		if last != c {
			if cnt > 1 {
				enc += fmt.Sprintf("%d", cnt)
			}
			enc += string(last)
			last = c
			cnt = 1
			continue
		}
		cnt++
	}

	if cnt > 1 {
		enc += fmt.Sprintf("%d", cnt)
	}
	enc += string(last)

	return enc
}

func RunLengthDecode(input string) string {
	dec := ""
	numSeq := ""
	for i := 0; i < len(input); i++ {

		if '0' <= input[i] && input[i] <= '9' {
			numSeq += string(input[i])
		} else {

			if numSeq != "" {
				num, _ := strconv.Atoi(numSeq)
				for j := 0; j < num; j++ {
					dec += string(input[i])
				}
				numSeq = ""
			} else {
				dec += string(input[i])
			}

		}

	}

	return dec
}
