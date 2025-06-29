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

type Config struct {
	Debug      bool     `json:"debug"`
	Server     []Server `json:"server"`
	StrictAuth bool     `json:"strict_auth"`
}

var gConf Config
var gShouldStartServer bool = false

func Parse(raw json.RawMessage) error {
	gShouldStartServer = true
	err := json.Unmarshal(raw, &gConf)
	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("webserver", Parse, Save)
}
