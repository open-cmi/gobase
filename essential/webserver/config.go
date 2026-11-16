package webserver

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type Server struct {
	Address  string `json:"address"`
	Port     int    `json:"port,omitempty"`
	Proto    string `json:"proto"`
	CertFile string `json:"cert,omitempty"`
	KeyFile  string `json:"key,omitempty"`
}

type StaticRule struct {
	Base  string `json:"base"`
	Alias string `json:"alias,omitempty"`
	Root  string `json:"root,omitempty"`
	Index string `json:"index,omitempty"`
}

type Config struct {
	Debug  bool         `json:"debug"`
	Server []Server     `json:"server"`
	Static []StaticRule `json:"static,omitempty"`
}

var gConf Config

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func GetConfig() *Config {
	return &gConf
}

func init() {
	config.RegisterConfig("webserver", Parse, Save)
}

