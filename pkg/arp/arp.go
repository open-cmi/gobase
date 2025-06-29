package arp

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
)

func GetHardwareAddrByIP(ip string) (net.HardwareAddr, error) {
	var ether net.HardwareAddr
	f, err := os.Open("/proc/net/arp")
	if err != nil {
		return ether, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	// IP address   HW type   Flags   HW address    Mask    Device
	s.Scan() // skip the field descriptions

	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		ipAddr := fields[0]
		hwAddr := fields[3]
		if ip == ipAddr {
			return net.ParseMAC(hwAddr)
		}
	}

	return ether, errors.New("arp item not found")
}
