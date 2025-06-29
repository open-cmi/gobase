package eyas

import (
	"os"
	"path/filepath"
	"strings"
)

type SoftwareVersion struct {
	BuildTime     string `json:"build_time"`
	Prod          string `json:"prod"`
	BuildPlatform string `json:"build_platform"`
}

func CurrentVersion() (v *SoftwareVersion, err error) {
	buildfile := filepath.Join(GetRootPath(), "BUILDINFO")
	contentByte, err := os.ReadFile(buildfile)
	if err != nil {
		return nil, err
	}
	var version SoftwareVersion
	arrs := strings.Split(string(contentByte), "\n")
	for index, line := range arrs {
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		if index == 0 {
			values := strings.SplitN(line, ":", 2)
			version.BuildTime = strings.Trim(values[1], " ")
		} else if index == 1 {
			values := strings.Split(line, ":")
			version.Prod = strings.Trim(values[1], " ")
		} else if index == 2 {
			values := strings.Split(line, ":")
			version.BuildPlatform = strings.Trim(values[1], " ")
		}
	}
	return &version, nil
}
