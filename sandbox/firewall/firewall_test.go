package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Match(t *testing.T) {
	rules := [][2]string{
		{"192.168.1.0/24", "ALLOW"},
		{"10.0.0.0/16", "DENY"},
		{"10.1.0.0/16", "ALLOW"},
		{"8.8.8.8", "ALLOW"},
	}

	tests := []struct {
		ip        string
		expResult bool
		expError  error
		note      string
	}{
		{
			"192.168.1.10",
			true,
			nil,
			"should match 192.168.1.0/24 ALLOW",
		},
		{
			"10.0.0.10",
			false,
			nil,
			"should match 10.0.0.0/16 DENY",
		},
		{
			"10.0.10.10",
			false,
			nil,
			"should match 10.0.0.0/16 DENY",
		},
		{
			"10.0.10.10",
			false,
			nil,
			"should match 10.0.0.0/16 DENY",
		},
		{
			"10.1.254.128",
			true,
			nil,
			"should match 10.1.0.0/16 ALLOW",
		},
		{
			"8.8.8.8",
			true,
			nil,
			"should match 8.8.8.8 ALLOW",
		},
		{
			"8.qwe.8.8",
			false,
			fmt.Errorf("can`t match CIDR"),
			"can`t match CIDR error expected",
		},
	}

	fw, err := NewFirewall(rules)
	assert.NoError(t, err)

	for _, test := range tests {
		res, err := fw.Allowed(rules, test.ip)
		assert.Equal(t, test.expResult, res, test.note)
		if test.expError != nil {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func Test_NewFirewall(t *testing.T) {
	fw, err := NewFirewall([][2]string{
		{"192.err.1.0/24", "ALLOW"},
		{"10.0.0.0/16", "DENY"},
	})
	assert.Error(t, err)
	assert.Nil(t, fw)

	fw, err = NewFirewall([][2]string{
		{"192.168.1.0/24", "GRANTED"},
		{"10.0.0.0/16", "DENY"},
	})
	assert.Error(t, err)
	assert.Nil(t, fw)

	fw, err = NewFirewall([][2]string{
		{"192.168.1.0/33", "ALLOW"},
	})
	assert.Error(t, err)
	assert.Nil(t, fw)

	fw, err = NewFirewall([][2]string{
		{"192.168.1.0/-5", "ALLOW"},
	})
	assert.Error(t, err)
	assert.Nil(t, fw)
}

func Test_CidrToBin(t *testing.T) {
	fw, err := NewFirewall([][2]string{})
	assert.NoError(t, err)

	ip, mask, err := fw.cidrToBin("192.168.1.1/24")
	assert.NoError(t, err)
	assert.Equal(t, "11000000101010000000000100000001", fmt.Sprintf("%b", ip))
	assert.Equal(t, "11111111111111111111111100000000", fmt.Sprintf("%b", mask))

	log.Printf("\n%b\n%32b\n", ip, mask)
}

func Test_BasicMatch(t *testing.T) {
	fw, err := NewFirewall([][2]string{})
	assert.NoError(t, err)

	rIp, ruleMask, err := fw.cidrToBin("192.168.1.0/24")
	assert.NoError(t, err)

	ip, err := fw.ipToBin("192.168.1.10")
	assert.NoError(t, err)

	log.Printf("\n%b\n%32b\n%32b\n", rIp, ruleMask, ip)

	ip = ip & ruleMask

	log.Printf("\n%32b\n", ip)

	assert.True(t, ip == rIp)
}
