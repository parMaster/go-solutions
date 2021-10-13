package snippets

import (
	"fmt"
	"testing"
)

const (
	isAdmin          = 1 << iota // 000001
	isHeadquarters               // 000010
	canSeeFinancials             // 000100
	canSeeAfrica                 // 001000
	canSeeAsia                   // 010000
	canSeeEurope                 // 100000
)

func TestBitMaskOperations(t *testing.T) {

	var roles byte = isAdmin | canSeeFinancials | canSeeEurope
	fmt.Printf("%b\n", roles)

	fmt.Printf("Is Admin? %v\n", isAdmin&roles == isAdmin)
	fmt.Printf("Is HQ? %v\n", isHeadquarters&roles == isHeadquarters)

	cases := []struct {
		name     string
		expr     byte
		expected byte
	}{
		{
			"isAdmin",
			isAdmin & roles,
			isAdmin,
		},
		{
			"canSeeFinancials",
			canSeeFinancials & roles,
			canSeeFinancials,
		},
		{
			"isHeadquarters",
			isHeadquarters & roles,
			0,
		},
		{
			"canSeeAsia",
			canSeeAsia & roles,
			0,
		},
	}

	for _, pair := range cases {
		if pair.expr != pair.expected {
			t.Error(
				"Test: ", pair.name,
				" got", pair.expr,
				" expected", pair.expected,
			)
		}
	}
}

const (
	_  = iota // zero ignored
	KB = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func TestFileSizeConversions(t *testing.T) {
	fileSize := 4000000000.
	fmt.Printf("%.2fGB ", fileSize/GB)
}
