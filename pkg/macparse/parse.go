package macparse

import (
	"net"
	"strings"
)

func ParseThreeSectionMAC(macStr string) (string, error) {
	s := strings.ReplaceAll(macStr, "-", ".")
	hw, err := net.ParseMAC(s)
	if err != nil {
		return "", err
	}

	return hw.String(), nil
}

func ParseMAC(macStr string) (string, error) {

	if len(macStr) == 14 && strings.Count(macStr, "-") == 2 {
		return ParseThreeSectionMAC(macStr)
	}

	hw, err := net.ParseMAC(macStr)
	if err != nil {
		return "", err
	}

	return hw.String(), nil
}
