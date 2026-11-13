package dev

import (
	"os/exec"
	"strings"
	"sync"
)

var (
	cachedDeviceID string
	deviceIDMu     sync.RWMutex
)

func parseIORegProperty(output string, key string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, `"`+key+`"`) {
			arr := strings.SplitN(line, "=", 2)
			if len(arr) == 2 {
				val := strings.Trim(arr[1], "\t\r\n\" ")
				if val != "" {
					return val
				}
			}
		}
	}
	return ""
}

// GetDeviceID func
func GetDeviceID() string {
	deviceIDMu.RLock()
	if cachedDeviceID != "" {
		id := cachedDeviceID
		deviceIDMu.RUnlock()
		return id
	}
	deviceIDMu.RUnlock()

	deviceIDMu.Lock()
	defer deviceIDMu.Unlock()
	if cachedDeviceID != "" {
		return cachedDeviceID
	}

	output, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return ""
	}

	text := string(output)
	deviceID := parseIORegProperty(text, "IOPlatformSerialNumber")
	if deviceID == "" {
		deviceID = parseIORegProperty(text, "IOPlatformUUID")
	}

	if deviceID != "" {
		cachedDeviceID = deviceID
	}
	return deviceID
}
