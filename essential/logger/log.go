package logger

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/pkg/eyas"
	"github.com/open-cmi/gobase/pkg/logx"
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

var gbaseLogger Feature

func GetLogger() Feature {
	return gbaseLogger
}

type Config struct {
	Level string `json:"level"`
	Path  string `json:"path,omitempty"`
}

func Error(v ...interface{})                 { gbaseLogger.Error(v...) }
func Errorf(format string, v ...interface{}) { gbaseLogger.Errorf(format, v...) }
func Warn(v ...interface{})                  { gbaseLogger.Warn(v...) }
func Warnf(format string, v ...interface{})  { gbaseLogger.Warnf(format, v...) }
func Info(v ...interface{})                  { gbaseLogger.Info(v...) }
func Infof(format string, v ...interface{})  { gbaseLogger.Infof(format, v...) }
func Debug(v ...interface{})                 { gbaseLogger.Debug(v...) }
func Debugf(format string, v ...interface{}) { gbaseLogger.Debugf(format, v...) }

var gConf Config

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	level := logx.Info
	switch gConf.Level {
	case "debug":
		level = logx.Debug
	case "info":
		level = logx.Info
	case "warn":
		level = logx.Warn
	case "error":
		level = logx.Error
	}

	logPath := gConf.Path
	if logPath != "" && logPath != "stdout" && strings.HasPrefix(logPath, ".") {
		rp := eyas.GetWorkingDir()
		logPath = filepath.Join(rp, logPath)
	}

	gbaseLogger = logx.NewLogger(&logx.Option{
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
	gbaseLogger = logx.NewLogger(&logx.Option{
		Dir:   "",
		Level: logx.Debug,
	})
	config.RegisterConfig("log", Parse, Save)
}
