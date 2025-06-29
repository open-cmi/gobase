package net

import "net"

func GetMacAddress(name string) (string, error) {
	intf, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	return intf.HardwareAddr.String(), nil
}
