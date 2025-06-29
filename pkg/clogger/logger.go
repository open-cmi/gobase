package clogger

import (
	"log"

	"github.com/longbozhan/timewriter"
)

type Level uint8

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Option struct {
	Dir        string
	Level      Level
	Compress   bool
	ReserveDay int
}

type Logger struct {
	printer map[Level]*log.Logger
	option  Option
}

func NewLogger(option *Option) *Logger {
	var logger Logger

	wr := timewriter.TimeWriter{
		Dir:        option.Dir,        // 日志路径
		Compress:   option.Compress,   // 是否压缩
		ReserveDay: option.ReserveDay, // 老化天数
	}

	logger.printer = make(map[Level]*log.Logger, 0)
	logger.printer[Debug] = log.New(&wr, "[DEBUG] ", log.Ldate|log.Lmicroseconds)
	logger.printer[Info] = log.New(&wr, "[INFO] ", log.Ldate|log.Lmicroseconds)
	logger.printer[Warn] = log.New(&wr, "[WARN] ", log.Ldate|log.Lmicroseconds)
	logger.printer[Error] = log.New(&wr, "[ERROR] ", log.Ldate|log.Lmicroseconds)

	logger.option = *option
	return &logger
}

func (l *Logger) SetLevel(level Level) {
	if level > Error {
		l.option.Level = Error
		return
	}

	l.option.Level = level
}

func (l *Logger) Printf(level Level, format string, args ...interface{}) {
	if level < l.option.Level || level > Error {
		return
	}

	printer := l.printer[level]
	printer.Printf(format, args...)
}

func (l *Logger) Println(level Level, args ...interface{}) {
	if level < l.option.Level || level > Error {
		return
	}

	printer := l.printer[level]
	printer.Println(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if Debug < l.option.Level {
		return
	}

	printer := l.printer[Debug]
	printer.Printf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	if Debug < l.option.Level {
		return
	}

	printer := l.printer[Debug]
	printer.Println(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if Info < l.option.Level {
		return
	}
	printer := l.printer[Info]
	printer.Printf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	if Info < l.option.Level {
		return
	}
	printer := l.printer[Info]
	printer.Println(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if Warn < l.option.Level {
		return
	}
	printer := l.printer[Warn]
	printer.Printf(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	if Warn < l.option.Level {
		return
	}
	printer := l.printer[Warn]
	printer.Println(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if Error < l.option.Level {
		return
	}
	printer := l.printer[Error]
	printer.Printf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	if Error < l.option.Level {
		return
	}
	printer := l.printer[Error]
	printer.Println(args...)
}
