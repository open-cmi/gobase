package cmdctl

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/essential/logger"
)

var manager *Manager

// Config process config
type Config struct {
	Name       string `json:"name"`
	ExecStart  string `json:"exec_start"`
	RestartSec int    `json:"restart_sec"`
	StopSignal int    `json:"stop_signal"`
}

type CommandConfig struct {
	Services []Config `json:"services"`
}

var gConf CommandConfig

func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	for _, s := range gConf.Services {
		manager.AddProcess(&s)
		err := manager.StartProcess(s.Name)
		if err != nil {
			logger.Errorf("start process failed: %s\n", err.Error())
			return err
		}
	}
	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	manager = NewManager()

	config.RegisterConfig("process", Init, Save)
}
