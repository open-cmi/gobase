package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-cmi/gobase/pkg/confparser"
	"github.com/open-cmi/gobase/pkg/eyas"
)

const (
	ModeDebug   = "debug"
	ModeRelease = "release"
)

var (
	gConfCtx  *confparser.Context
	gConfFile string
	gMode     = ModeRelease
)

// GetMode returns the global runtime mode ("debug" or "release").
func GetMode() string {
	return gMode
}

// IsDebug reports whether the global runtime mode is ModeDebug.
func IsDebug() bool {
	return gMode == ModeDebug
}

func parseMode(raw json.RawMessage) error {
	var mode string
	if err := json.Unmarshal(raw, &mode); err != nil {
		return err
	}
	if mode == "" {
		gMode = ModeRelease
		return nil
	}
	if mode != ModeDebug && mode != ModeRelease {
		return fmt.Errorf("invalid mode %q: must be %q or %q", mode, ModeDebug, ModeRelease)
	}
	gMode = mode
	return nil
}

func saveMode() json.RawMessage {
	raw, _ := json.Marshal(gMode)
	return raw
}

// ResolveConfigFile returns configfile if non-empty, or falls back to
// filepath.Join(eyas.GetConfDir(), "config.json") when it exists.
func ResolveConfigFile(configfile string) (string, error) {
	if configfile != "" {
		return configfile, nil
	}
	confDir := eyas.GetConfDir()
	if confDir != "" {
		defaultConfig := filepath.Join(confDir, "config.json")
		if fi, err := os.Stat(defaultConfig); err == nil && !fi.IsDir() {
			return defaultConfig, nil
		}
		return "", fmt.Errorf("config file must not be empty (default %s not found)", defaultConfig)
	}
	return "", errors.New("config file must not be empty")
}

// Init config init
func Init(configfile string) error {
	return gConfCtx.Load(configfile)
}

func Save() {
	if gConfCtx != nil {
		gConfCtx.Save()
	}
}

// RegisterConfig register config
func RegisterConfig(name string, parseFunc func(json.RawMessage) error, saveFunc func() json.RawMessage) error {

	if gConfCtx == nil {
		gConfCtx = confparser.NewContext()
	}

	var opt confparser.Option
	opt.Name = name
	opt.ParseFunc = parseFunc
	opt.SaveFunc = saveFunc
	return gConfCtx.Register(&opt)
}

func SetConfigFile(configFile string) {
	gConfFile = configFile
}

func init() {
	_ = RegisterConfig("mode", parseMode, saveMode)
}
