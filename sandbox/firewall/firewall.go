package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Firewall struct {
	rules [][2]string
}

func NewFirewall(rules [][2]string) (*Firewall, error) {
	r := [][2]string{}
	copy(r, rules)
	fw := &Firewall{rules: r}

	for _, rule := range rules {
		if len(rule) != 2 {
			return nil, fmt.Errorf("rules format error: wrong len")
		}
		if _, _, err := fw.cidrToBin(rule[0]); err != nil {
			return nil, fmt.Errorf("rules format error: %w", err)
		}
		if rule[1] != "ALLOW" && rule[1] != "DENY" {
			return nil, fmt.Errorf("invalid action: %v", rule[1])
		}
	}

	return fw, nil
}

// Match returns true if masked ip matches CIDR
func (f *Firewall) Match(cidr, strIp string) (bool, error) {
	ruleIp, ruleMask, err := f.cidrToBin(cidr)
	if err != nil {
		return false, err
	}

	ip, err := f.ipToBin(strIp)
	if err != nil {
		return false, err
	}

	ip = ip & ruleMask

	return ip == ruleIp, nil
}

// Allowed returns true if ip matched against one the rules and that rule is ALLOW
// returns false if ip matched the DENY rule or there's no match
func (f *Firewall) Allowed(rules [][2]string, ip string) (bool, error) {

	for _, rule := range rules {
		matched, err := f.Match(rule[0], ip)
		if err != nil {
			return false, fmt.Errorf("can`t match CIDR: %w", err)
		}
		if matched && rule[1] == "ALLOW" {
			return true, nil
		}
		if matched && rule[1] == "DENY" {
			return false, nil
		}
	}

	return false, nil
}

// ipToBin converts string ip (i.e. "192.168.1.10") to uint32 of binary representation
func (f Firewall) ipToBin(strIp string) (ip uint32, err error) {
	ipBytes := strings.Split(strIp, ".")
	for i, v := range ipBytes {
		ipByte, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("convertion error (%v to int): %w", v, err)
		}
		if ipByte > 255 || ipByte < 0 {
			return 0, fmt.Errorf("invalid ip byte (%v) in %s", ipByte, strIp)
		}
		ip |= uint32(ipByte) << (8 * (3 - i))
	}
	return ip, nil
}

// cidrToBin converts CIDR (i.e. "192.168.1.0/24") to uint32 representation of IP and subnet mask
func (f Firewall) cidrToBin(cidr string) (ip uint32, mask uint32, err error) {
	parts := strings.Split(cidr, "/")

	if len(parts) < 1 {
		return 0, 0, fmt.Errorf("wrong cidr format")
	}

	ip, err = f.ipToBin(parts[0])
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
	if maskInt > 32 || maskInt <= 0 {
		return 0, 0, fmt.Errorf("0<=mask<=32 condition is not satisfied: %v", maskInt)
	}
	mask = 1<<maskInt - 1
	mask = mask << (32 - maskInt)
	return ip, uint32(mask), nil
}
