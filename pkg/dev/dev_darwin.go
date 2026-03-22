package dev

import (
	"os/exec"
	"strings"
)

// GetDeviceID func
func GetDeviceID() string {
	grep := exec.Command("grep", "IOPlatformSerialNumber")
	ioreg := exec.Command("ioreg", "-l")

	// Get ps's stdout and attach it to grep's stdin.
	pipe, err := ioreg.StdoutPipe()
	if err != nil {
		return ""
	}

	defer pipe.Close()
	grep.Stdin = pipe
	// Run ps first.
	ioreg.Start()

	// Run and get the output of grep.
	output, err := grep.Output()
	if err != nil {
		return ""
	}

	arr := strings.Split(string(output), "=")
	if len(arr) != 2 {
		return ""
	}
	deviceid := strings.Trim(arr[1], "\t\n\" ")
	return deviceid
}
