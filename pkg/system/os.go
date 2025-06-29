package system

import (
	"os"
	"runtime"
)

func GetOSRelease() string {
	if runtime.GOOS == "linux" {
		_, err := os.Stat("/usr/bin/yum")
		if err == nil {
			// centos系，包含centos、openalinos等
			return "centos"
		}

		_, err = os.Stat("/usr/bin/apt")
		if err == nil {
			// debian系，包含debian, ubuntu, kali等
			return "debian"
		}
		// 如果不是以上系统，则默认是openwrt，因为其他系统暂时没见过
		return "openwrt"
	}
	return runtime.GOOS
}
