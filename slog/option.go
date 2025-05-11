package slog

import "os"

// Option 选项函数
type Option func(*logger)

// WithLevel 设置日志级别
func WithLevel(level string) Option {
	switch level {
	case "info", "warn", "error", "debug":
		break
	default:
		level = "info"
	}

	return func(l *logger) {
		l.level = level
	}
}

// WithMaxSize 设置日志文件最大大小
func WithMaxSize(maxSize int) Option {
	if maxSize == 0 {
		maxSize = 100 * 1024 * 1024 // 100 MB
	}

	return func(l *logger) {
		l.maxSize = maxSize
	}
}

// WithRotateInterval 设置日志文件轮转间隔
func WithRotateInterval(interval string) Option {
	return func(l *logger) {
		l.rotateInterval = interval
	}
}

// WithBufferSize 设置日志文件缓冲区大小
func WithBufferSize(bufferSize int) Option {
	return func(l *logger) {
		l.bufferSize = bufferSize
	}
}

// WithMaxBackups 设置日志文件最大备份数
func WithMaxBackups(maxBackups uint) Option {
	return func(l *logger) {
		l.maxBackups = maxBackups
	}
}

// WithCompress 设置日志文件是否压缩
func WithCompress() Option {
	return func(l *logger) {
		l.compress = true
	}
}

// WithTimeFormat 设置日志文件时间格式
func WithTimeFormat(timeFormat string) Option {
	return func(l *logger) {
		l.timeFormat = timeFormat
	}
}

// WithFilePerm 设置日志文件权限
func WithFilePerm(perm os.FileMode) Option {
	return func(l *logger) {
		l.filePerm = perm
	}
}

// WithPrintAfterInitialized 设置是否在初始化后打印日志
func WithPrintAfterInitialized() Option {
	return func(l *logger) {
		l.printAfterInitialized = true
	}
}

// WithColor 设置是否启用颜色输出
// 默认为false，不启用颜色
func WithColor(enable bool) Option {
	return func(l *logger) {
		l.enableColor = enable
	}
}

// WithCallerSkip 设置调用栈跳过级别
// 默认为8，跳过slog库内部调用 (这个参数调试出来的)
// 小于8，则显示slog库内部调用位置
// 设置为0，则从底层调用开始展示
func WithCallerSkip(skip int) Option {
	return func(l *logger) {
		if skip < 0 {
			skip = 0
		}
		l.callerSkip = skip
	}
}
