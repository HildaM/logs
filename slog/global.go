package slog

import (
	"fmt"
	"sync"
)

var (
	// defaultLogger 是全局Logger实例
	defaultLogger LoggerWithWriter
	// once 确保全局Logger只被初始化一次
	once sync.Once
	// mu 保护defaultLogger的并发访问
	mu sync.RWMutex
)

// Init 初始化全局标准输出Logger
// 示例：slog.Init()
func Init(opts ...Option) error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewStdoutLogger(opts...)
	})
	return err
}

// InitFile 初始化全局文件Logger
// 示例：slog.InitFile("./app.log", slog.WithLevel("debug"))
func InitFile(path string, opts ...Option) error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewFileLogger(path, opts...)
	})
	return err
}

// SetLogger 设置自定义Logger为全局Logger
func SetLogger(logger LoggerWithWriter) {
	if logger == nil {
		panic("logger cannot be nil")
	}

	mu.Lock()
	defer mu.Unlock()
	defaultLogger = logger
}

// GetLogger 获取当前全局Logger
func GetLogger() LoggerWithWriter {
	mu.RLock()
	defer mu.RUnlock()

	// 如果尚未初始化，则创建一个默认的标准输出Logger
	if defaultLogger == nil {
		// 这里不使用once，因为用户可能会调用SetLogger替换Logger
		var err error
		defaultLogger, err = NewStdoutLogger()
		if err != nil {
			// 如果默认日志器都创建失败，说明系统环境有严重问题
			panic(fmt.Sprintf("failed to create default logger: %v", err))
		}
	}

	return defaultLogger
}

// 包级方法，直接调用全局Logger

// Info 记录Info级别日志
func Info(msg string, args ...interface{}) {
	GetLogger().Info(msg, args...)
}

// Error 记录Error级别日志
func Error(msg string, args ...interface{}) {
	GetLogger().Error(msg, args...)
}

// Warn 记录Warn级别日志
func Warn(msg string, args ...interface{}) {
	GetLogger().Warn(msg, args...)
}

// Debug 记录Debug级别日志
func Debug(msg string, args ...interface{}) {
	GetLogger().Debug(msg, args...)
}

// Fatal 记录Fatal级别日志
func Fatal(msg string, args ...interface{}) {
	GetLogger().Fatal(msg, args...)
}

// Write 实现io.Writer接口
func Write(p []byte) (int, error) {
	return GetLogger().Write(p)
}
