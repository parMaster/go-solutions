package sandbox

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ipToBin converts string ip (i.e. "192.168.1.10") to uint32 of binary representation
func ipToBin(strIp string) (uint32, error) {
	ipBytes := strings.Split(strIp, ".")
	var ip uint32
	for i, v := range ipBytes {
		ipByte, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		ip |= uint32(ipByte) << (8 * (3 - i))
	}
	return ip, nil
}

// cidrToBin converts CIDR (i.e. "192.168.1.0/24") to binary IP + subnet mask
func cidrToBin(cidr string) (uint32, uint32, error) {
	parts := strings.Split(cidr, "/")

	if len(parts) < 1 {
		return 0, 0, fmt.Errorf("wrong cidr format")
	}

	ip, err := ipToBin(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("can`t convert to int %s: %w", parts[0], err)
	}

	// no subnet
	if len(parts) == 1 {
		mask := 1<<32 - 1
		return ip, uint32(mask), nil
	}

	maskInt, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("can`t convert to int %s: %w", parts[1], err)
	}
	mask := 1<<maskInt - 1
	mask = mask << (32 - maskInt)
	return ip, uint32(mask), nil
}

// matchCidr returns true if masked ip matches CIDR
func matchCidr(cidr, strIp string) (bool, error) {
	ruleIp, ruleMask, err := cidrToBin(cidr)
	if err != nil {
		return false, err
	}

	ip, err := ipToBin(strIp)
	if err != nil {
		return false, err
	}

	ip = ip & ruleMask

	return ip == ruleIp, nil
}

// matchIp returns true if ip matched against one the rules and that rule is ALLOW
// returns false if ip matched the DENY rule or there's no match
func matchIp(rules [][2]string, ip string) bool {

	for _, rule := range rules {
		matched, err := matchCidr(rule[0], ip)
		if err != nil {
			panic(err)
		}
		if matched && rule[1] == "ALLOW" {
			return true
		}
		if matched && rule[1] == "DENY" {
			return false
		}
	}

	return false
}

func Test_Match(t *testing.T) {
	rules := [][2]string{
		{"192.168.1.0/24", "ALLOW"},
		{"10.0.0.0/16", "DENY"},
		{"10.1.0.0/16", "ALLOW"},
		{"8.8.8.8", "ALLOW"},
	}

	tests := []struct {
		ip   string
		exp  bool
		note string
	}{
		{
			"192.168.1.10",
			true,
			"should match 192.168.1.0/24 ALLOW",
		},
		{
			"10.0.0.10",
			false,
			"should match 10.0.0.0/16 DENY",
		},
		{
			"10.0.10.10",
			false,
			"should match 10.0.0.0/16 DENY",
		},
		{
			"10.0.10.10",
			false,
			"should match 10.0.0.0/16 DENY",
		},
		{
			"10.1.254.128",
			true,
			"should match 10.1.0.0/16 ALLOW",
		},
		{
			"8.8.8.8",
			true,
			"should match 8.8.8.8 ALLOW",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.exp, matchIp(rules, test.ip), test.note)
	}
}

func Test_CidrToBin(t *testing.T) {
	ip, mask, err := cidrToBin("192.168.1.1/24")
	assert.NoError(t, err)
	assert.Equal(t, "11000000101010000000000100000001", fmt.Sprintf("%b", ip))
	assert.Equal(t, "11111111111111111111111100000000", fmt.Sprintf("%b", mask))

	fmt.Printf("%b\n%32b\n", ip, mask)
}

func Test_BasicMatch(t *testing.T) {
	rIp, ruleMask, err := cidrToBin("192.168.1.0/24")
	assert.NoError(t, err)

	ip, err := ipToBin("192.168.1.10")
	assert.NoError(t, err)

	fmt.Printf("%b\n%32b\n%32b\n", rIp, ruleMask, ip)

	ip = ip & ruleMask

	fmt.Printf("%32b\n", ip)

	assert.True(t, ip == rIp)
}
