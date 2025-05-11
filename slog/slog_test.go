package slog

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	l := newLogger()
	assert.NotNil(t, l)
	assert.Equal(t, "info", l.level)
	assert.Equal(t, time.DateTime, l.timeFormat)
	assert.Equal(t, "1w", l.rotateInterval)
	assert.True(t, l.compress)
	assert.Equal(t, os.FileMode(0700), l.filePerm)
	assert.False(t, l.enableColor) // 默认不启用颜色

	// 测试应用Options
	l = newLogger(
		WithLevel("debug"),
		WithMaxSize(200),
		WithRotateInterval("1d"),
		WithBufferSize(1024),
		WithMaxBackups(5),
		WithTimeFormat(time.RFC3339),
		WithFilePerm(0644),
		WithPrintAfterInitialized(),
		WithColor(true),
	)
	assert.Equal(t, "debug", l.level)
	assert.Equal(t, 200, l.maxSize)
	assert.Equal(t, "1d", l.rotateInterval)
	assert.Equal(t, 1024, l.bufferSize)
	assert.Equal(t, uint(5), l.maxBackups)
	assert.Equal(t, time.RFC3339, l.timeFormat)
	assert.Equal(t, os.FileMode(0644), l.filePerm)
	assert.True(t, l.printAfterInitialized)
	assert.True(t, l.enableColor)
}

func TestOptions(t *testing.T) {
	// 测试WithLevel
	l := &logger{}
	WithLevel("debug")(l)
	assert.Equal(t, "debug", l.level)

	// 测试无效level
	l = &logger{}
	WithLevel("invalid")(l)
	assert.Equal(t, "info", l.level)

	// 测试WithMaxSize
	l = &logger{}
	WithMaxSize(0)(l)
	assert.Equal(t, 100*1024*1024, l.maxSize)

	l = &logger{}
	WithMaxSize(1024)(l)
	assert.Equal(t, 1024, l.maxSize)

	// 测试其他Option
	l = &logger{}
	WithRotateInterval("2h")(l)
	assert.Equal(t, "2h", l.rotateInterval)

	l = &logger{}
	WithBufferSize(2048)(l)
	assert.Equal(t, 2048, l.bufferSize)

	l = &logger{}
	WithMaxBackups(10)(l)
	assert.Equal(t, uint(10), l.maxBackups)

	l = &logger{}
	WithCompress()(l)
	assert.True(t, l.compress)

	l = &logger{}
	WithTimeFormat(time.RFC3339)(l)
	assert.Equal(t, time.RFC3339, l.timeFormat)

	l = &logger{}
	WithFilePerm(0644)(l)
	assert.Equal(t, os.FileMode(0644), l.filePerm)

	l = &logger{}
	WithPrintAfterInitialized()(l)
	assert.True(t, l.printAfterInitialized)

	// 测试WithColor
	l = &logger{}
	WithColor(true)(l)
	assert.True(t, l.enableColor)

	l = &logger{}
	assert.False(t, l.enableColor) // 默认为false
}

func TestNewStdoutLogger(t *testing.T) {
	logger, err := NewStdoutLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// 测试不同级别的日志记录
	logger.Info("test info: %+v, name: %+v, age: %+v", "test info", "man!", 12343)
	logger.Error("test error: %+v", "test error")
	logger.Warn("test warn: %+v", "test warn")
	logger.Debug("test debug: %+v", "test debug") // 默认不会记录debug日志
	logger.Fatal("test fatal: %+v", "test fatal") // 记录fatal级别日志

	// 测试debug级别
	logger, err = NewStdoutLogger(WithLevel("debug"))
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	logger.Debug("test debug with debug level")

	// 测试Write方法
	n, err := logger.Write([]byte("test write"))
	assert.NoError(t, err)
	assert.NotZero(t, n)
}

func TestNewFileLogger(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logger_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	// 测试基于大小的日志轮转
	logger, err := NewFileLogger(logFile,
		WithMaxSize(1024),
		WithMaxBackups(3),
		WithLevel("debug"),
		WithPrintAfterInitialized(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// 写入日志
	logger.Info("test info")
	logger.Error("test error")
	logger.Warn("test warn")
	logger.Debug("test debug")

	// 测试Write接口
	message := []byte("test direct write")
	n, err := logger.Write(message)
	assert.NoError(t, err)
	assert.Equal(t, len(message), n)

	// 检查文件是否创建
	_, err = os.Stat(logFile)
	assert.NoError(t, err)

	// 测试基于时间的日志轮转
	logger, err = NewFileLogger(logFile,
		WithRotateInterval("1h"),
		WithMaxBackups(2),
		WithBufferSize(0), // 立即写入
	)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// 写入日志
	logger.Info("test time rotate info")
}

func TestHandlerRotateTime(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rotate_time_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "rotate_time.log")

	// 测试各种时间间隔
	intervals := []string{"1s", "5m", "2h", "1d", "2w"}
	for _, interval := range intervals {
		l := &logger{
			level:          "info",
			rotateInterval: interval,
			maxBackups:     2,
			compress:       true,
			filePerm:       0644,
		}

		handler, err := handlerRotateTime(l, logFile)
		assert.NoError(t, err)
		assert.NotNil(t, handler)
	}

	// 测试无效的间隔
	invalid := []string{"", "x", "1", "1y"}
	for _, interval := range invalid {
		l := &logger{rotateInterval: interval}
		_, err := handlerRotateTime(l, logFile)
		assert.Error(t, err)
	}
}

func TestHandlerRotateFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rotate_file_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "rotate_file.log")

	l := &logger{
		level:      "info",
		maxSize:    1024,
		maxBackups: 2,
		compress:   true,
		bufferSize: 0,
		filePerm:   0644,
	}

	handler, err := handlerRotateFile(l, logFile)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	// 测试写入
	writer := handler.Writer()
	_, err = writer.Write([]byte("test handlerRotateFile"))
	assert.NoError(t, err)
}

func TestParseLevels(t *testing.T) {
	// 测试info级别
	levels := parseLevels("info")
	assert.Contains(t, levels, slog.InfoLevel)
	assert.Contains(t, levels, slog.WarnLevel)
	assert.Contains(t, levels, slog.ErrorLevel)
	assert.NotContains(t, levels, slog.DebugLevel)

	// 测试debug级别
	levels = parseLevels("debug")
	assert.Contains(t, levels, slog.InfoLevel)
	assert.Contains(t, levels, slog.WarnLevel)
	assert.Contains(t, levels, slog.ErrorLevel)
	assert.Contains(t, levels, slog.DebugLevel)

	// 测试无效级别（默认为info)
	levels = parseLevels("invalid")
	assert.Contains(t, levels, slog.InfoLevel)
	assert.Contains(t, levels, slog.WarnLevel)
	assert.Contains(t, levels, slog.ErrorLevel)
	assert.NotContains(t, levels, slog.DebugLevel)
}

func TestGenLogFormatter(t *testing.T) {
	// 测试不启用颜色（默认）
	formatter := genLogFormatter(time.RFC3339, false)
	assert.NotNil(t, formatter)
	assert.False(t, formatter.EnableColor)
	assert.True(t, formatter.FullDisplay)
	assert.Equal(t, time.RFC3339, formatter.TimeFormat)

	// 测试启用颜色
	formatter = genLogFormatter(time.RFC3339, true)
	assert.NotNil(t, formatter)
	assert.True(t, formatter.EnableColor)
	assert.True(t, formatter.FullDisplay)
	assert.Equal(t, time.RFC3339, formatter.TimeFormat)
}
