package cmdctl

import (
	"syscall"
)

// ParentIsRunning parent is running
func ParentIsRunning(ppid int) bool {
	err := syscall.Kill(ppid, 0)
	if err != nil {
		return false
	}
	return true
}
