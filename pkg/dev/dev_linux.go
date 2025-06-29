package dev

import (
	"os"
	"os/exec"
	"strings"
)

// GetDeviceID func
func GetDeviceID() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	var args []string = []string{"-h", "--output=target", execPath}
	cmd := exec.Command("df", args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	arr := strings.Split(string(output), "\n")
	if len(arr) < 2 {
		return ""
	}
	mountPoint := arr[1]

	var args1 []string = []string{"-r", "-p", "-i", "-o", "NAME,MOUNTPOINT", "--noheadings"}
	cmdLsblk := exec.Command("lsblk", args1...)
	output, err = cmdLsblk.Output()
	if err != nil {
		return ""
	}
	arr = strings.Split(string(output), "\n")
	var device string = ""
	for _, line := range arr {
		lineArr := strings.Split(line, " ")
		if len(lineArr) != 2 {
			continue
		}
		mp := lineArr[1]
		disk := lineArr[0]
		if mp == mountPoint {
			device = disk
			break
		}
	}
	if device == "" {
		return ""
	}

	var args2 []string = []string{"-o", "value", device}
	cmdBlkid := exec.Command("blkid", args2...)
	output, err = cmdBlkid.Output()
	if err != nil {
		return ""
	}
	arr = strings.Split(string(output), "\n")
	if len(arr) < 3 {
		return ""
	}
	uuid := arr[0]
	return uuid
}

// func GetDeviceID() string {
// 	var deviceIDFiles []string = []string{
// 		"/sys/class/dmi/id/product_uuid",
// 		"/sys/block/mmcblk0/device/serial",
// 	}
// 	for _, filep := range deviceIDFiles {
// 		file, err := os.Open(filep)
// 		if err != nil {
// 			continue
// 		}

// 		data, _ := io.ReadAll(file)
// 		deviceid := strings.Trim(string(data), " \r\n\t ")
// 		file.Close()
// 		return deviceid
// 	}
// 	return ""
// }
