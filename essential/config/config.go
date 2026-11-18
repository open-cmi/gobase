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

var gConfCtx *confparser.Context

var gConfFile string

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
	err := gConfCtx.Load(configfile)
	return err
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
