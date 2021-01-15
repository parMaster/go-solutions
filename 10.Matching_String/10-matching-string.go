package main

func matchingStrings(strings []string, queries []string) []int32 {

	var m map[string]int32
	m = make(map[string]int32)

	for _, str := range strings {
		m[str]++
	}

	result := []int32{}
	for _, q := range queries {
		result = append(result, m[q])
	}

	return result
}

func main() {

	matchingStrings([]string{"aba", "baba", "aba", "xzxb"}, []string{"aba", "xzxb", "ab"})

}
