package main

import (
	"fmt"
)

func isValid(s string) string {

	freq := make([]int, int('z')-int('a')+1)

	for i := 0; i < len(s); i++ {
		freq[s[i]-97]++
	}

	cs := make(map[int]int)

	for i := range freq {
		if 0 != freq[i] {
			cs[freq[i]]++
		}
	}

	if 2 < len(cs) {
		return "NO"
	}

	if 2 == len(cs) {

		elem, ok := cs[1]
		if ok && elem == 1 {
			return "YES"
		}

		i1, i2, v1, v2 := 0, 0, 0, 0
		for i, v := range cs {
			if 0 == i1 {
				i1 = i
				v1 = v
			} else {
				if i > i1 {
					i2 = i
					v2 = v
				} else {
					i2 = i1
					i1 = i
					v2 = v1
					v1 = i
				}
			}
		}

		if v2 == 1 && i2-1 == i1 {
			return "YES"
		}

		return "NO"
	}

	return "YES"
}

func main() {

	fmt.Println(isValid("abcdefghhgfedecba"), "Exp: YES")
	fmt.Println(isValid("aaaabbcc"), "Exp: NO")
	fmt.Println(isValid("aabbccddee"), "Exp: YES")
	fmt.Println(isValid("aabbccddeefghi"), "Exp: NO")
	fmt.Println(isValid("abcdefghhgfedecbazz"), "Exp: YES")
	fmt.Println(isValid("aabbccddeefgghhhgffedecbazzz"), "Exp: YES")
	fmt.Println(isValid("aabbccddeeffgghhhgffedecbazzz"), "Exp: NO")
	fmt.Println(isValid("aabbc"), "Exp: YES")
	fmt.Println(isValid("ibfdgaeadiaefgbhbdghhhbgdfgeiccbiehhfcggchgghadhdhagfbahhddgghbdehidbibaeaagaeeigffcebfbaieggabcfbiiedcabfihchdfabifahcbhagccbdfifhghcadfiadeeaheeddddiecaicbgigccageicehfdhdgafaddhffadigfhhcaedcedecafeacbdacgfgfeeibgaiffdehigebhhehiaahfidibccdcdagifgaihacihadecgifihbebffebdfbchbgigeccahgihbcbcaggebaaafgfedbfgagfediddghdgbgehhhifhgcedechahidcbchebheihaadbbbiaiccededchdagfhccfdefigfibifabeiaccghcegfbcghaefifbachebaacbhbfgfddeceababbacgffbagidebeadfihaefefegbghgddbbgddeehgfbhafbccidebgehifafgbghafacgfdccgifdcbbbidfifhdaibgigebigaedeaaiadegfefbhacgddhchgcbgcaeaieiegiffchbgbebgbehbbfcebciiagacaiechdigbgbghefcahgbhfibhedaeeiffebdiabcifgccdefabccdghehfibfiifdaicfedagahhdcbhbicdgibgcedieihcichadgchgbdcdagaihebbabhibcihicadgadfcihdheefbhffiageddhgahaidfdhhdbgciiaciegchiiebfbcbhaeagccfhbfhaddagnfieihghfbaggiffbbfbecgaiiidccdceadbbdfgigibgcgchafccdchgifdeieicbaididhfcfdedbhaadedfageigfdehgcdaecaebebebfcieaecfagfdieaefdiedbcadchabhebgehiidfcgahcdhcdhgchhiiheffiifeegcfdgbdeffhgeghdfhbfbifgidcafbfcd"), "Exp: YES")
}
