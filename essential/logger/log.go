package logger

import (
	"encoding/json"
	"fmt"

	"path/filepath"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/pkg/clogger"
	"github.com/open-cmi/gobase/pkg/eyas"
)

type Feature interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

var Logger Feature

type Config struct {
	Level string `json:"level"`
	Path  string `json:"path,omitempty"`
}

func Error(v ...interface{}) {
	if Logger != nil {
		Logger.Error(v...)
	} else {
		fmt.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Errorf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if Logger != nil {
		Logger.Warn(v...)
	} else {
		fmt.Println(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Warnf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if Logger != nil {
		Logger.Info(v...)
	} else {
		fmt.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if Logger != nil {
		Logger.Debug(v...)
	} else {
		fmt.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

var gConf Config

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	if Logger != nil {
		return nil
	}

	level := clogger.Info
	switch gConf.Level {
	case "debug":
		level = clogger.Debug
	case "info":
		level = clogger.Info
	case "warn":
		level = clogger.Warn
	case "error":
		level = clogger.Error
	}
	logPath := gConf.Path
	if logPath == "" {
		rp := eyas.GetRootPath()
		logPath = filepath.Join(rp, "data")
	}

	Logger = clogger.NewLogger(&clogger.Option{
		Dir:        logPath,
		Compress:   true,
		Level:      level,
		ReserveDay: 30,
	})

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	gConf.Level = "debug"
	config.RegisterConfig("log", Parse, Save)
}
