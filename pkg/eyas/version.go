package eyas

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/open-cmi/gobase/essential/i18n"
)

type LongVersion struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
}

type ComponentVersion struct {
	LongVersion
	Name string `json:"name"`
}

type SoftwareVersion struct {
	BuildTime  string             `json:"build_time"`
	Prod       LongVersion        `json:"prod"`
	Eyas       LongVersion        `json:"eyas"`
	Components []ComponentVersion `json:"components"`
}

// v1.0.0-7-ge622ad89、v1.0.0-rc1-39-g4f62f61 或 v1.0.0
var (
	verRe    = regexp.MustCompile(`^v?\d+\.\d+\.\d+(?:-[0-9A-Za-z][0-9A-Za-z.-]*)?$`)
	commitRe = regexp.MustCompile(`-(\d+-g[0-9a-fA-F]+)$`)
)

func ParseLongVersion(version string) LongVersion {
	version = strings.TrimSpace(version)
	parsed := LongVersion{Version: version}
	matches := commitRe.FindStringSubmatchIndex(version)
	if len(matches) > 0 {
		parsed.Version = version[:matches[0]]
		parsed.Commit = version[matches[2]:matches[3]]
	}
	if !verRe.MatchString(parsed.Version) {
		return LongVersion{
			Version: version,
		}
	}
	return parsed
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
			version.Prod = ParseLongVersion(val)
		case "eyas":
			version.Eyas = ParseLongVersion(val)
		default:
			// 之后的都是组件版本
			origKey := strings.TrimSpace(parts[0])
			cv := ParseLongVersion(val)
			var com ComponentVersion
			com.Name = origKey
			com.Version = cv.Version
			com.Commit = cv.Commit
			version.Components = append(version.Components, com)
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
