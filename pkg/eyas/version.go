package eyas

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-cmi/gobase/essential/i18n"
)

type ComponentVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SoftwareVersion struct {
	BuildTime  string             `json:"build_time"`
	Prod       string             `json:"prod"`
	Eyas       string             `json:"eyas"`
	Components []ComponentVersion `json:"components"`
}

// ParseBuildInfo parses the key-value formatted BUILDINFO file content
func ParseBuildInfo(content []byte) *SoftwareVersion {
	version := &SoftwareVersion{
		Components: make([]ComponentVersion, 0),
	}

	lines := strings.Split(string(content), "\n")
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch key {
		case "build time", "build_time":
			version.BuildTime = val
		case "prod":
			version.Prod = val
		case "eyas":
			version.Eyas = val
		default:
			// 之后的都是组件版本
			origKey := strings.TrimSpace(parts[0])
			version.Components = append(version.Components, ComponentVersion{
				Name:    origKey,
				Version: val,
			})
		}
	}
	return version
}

func CurrentVersion() (v *SoftwareVersion, err error) {
	buildfile := filepath.Join(GetWorkingDir(), "BUILDINFO")
	contentByte, err := os.ReadFile(buildfile)
	if err != nil {
		return nil, errors.New(i18n.Sprintf("open version file failed"))
	}

	return ParseBuildInfo(contentByte), nil
}
