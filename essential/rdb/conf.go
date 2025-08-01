package rdb

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

// Config cache config
type Config struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
}

var gConf Config

func GetConf() *Config {
	return &gConf
}

// Parse db init
func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	gConf.Host = "127.0.0.1"
	gConf.Port = 25431
	config.RegisterConfig("rdb", Parse, Save)
}
