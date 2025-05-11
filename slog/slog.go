package slog

import "io"

// Logger 接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// LoggerWithWriter 接口
type LoggerWithWriter interface {
	Logger
	io.Writer
}
