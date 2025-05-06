package slog

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

// logger 结构体
type logger struct {
	sl *slog.Logger // slog库实例

	// 基础配置
	level                 string      // 日志级别
	maxSize               int         // 日志文件最大大小(默认 100MB)
	maxBackups            uint        // 日志文件最大备份数(默认 12)
	timeFormat            string      // 时间格式(默认 "2006-01-02 15:04:05")
	compress              bool        // 是否压缩(默认 true)
	filePerm              os.FileMode // 日志文件权限
	printAfterInitialized bool        // 是否在初始化后打印日志(默认 true)

	// 日志文件写入
	bufferSize int // 缓冲区大小(默认 8*1024. 如果为0, 则立即写入文件)

	// 日志文件轮转
	rotateInterval string    // 日志文件轮转间隔(e.g. `12h` (12 hours), `1d` (1 day), `1w` (1 week), `1m` (1 month))
	rotateWriter   io.Writer // 日志文件轮转写入器
}

// newLoggger 创建新的日志记录器
func newLogger(opts ...Option) *logger {
	sl := slog.New()
	l := &logger{
		sl:             sl,
		level:          "info",
		timeFormat:     time.DateTime,
		rotateInterval: "1w",
		compress:       true,
		filePerm:       0700,
	}
	// 应用选项
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// NewStdoutLogger 创建命令行标准输出日志
func NewStdoutLogger(opts ...Option) (LoggerWithWriter, error) {
	l := newLogger(opts...)

	// 创建日志格式化方式
	logFormatter := genLogFormatter(l.timeFormat)

	// 设置 slog handler
	h := handler.NewConsoleHandler(parseLevels(l.level))
	h.SetFormatter(logFormatter)
	l.sl.AddHandler(h)

	// 设置其他参数
	l.rotateWriter = h.Output
	l.sl.DoNothingOnPanicFatal() // panic时不处理

	return l, nil
}

// NewFileLogger 创建文件日志输出实例
func NewFileLogger(path string, opts ...Option) (LoggerWithWriter, error) {
	l := newLogger(opts...)
	if l.maxBackups == 0 {
		l.maxBackups = 12
	}

	logFormatter := genLogFormatter(l.timeFormat)

	// 文件大小轮转和时间轮转是两种不同的日志管理机制,设计上选择其一,简化逻辑实现和管理
	// 同时使用两种策略会导致复杂的交互问题（如先达到大小限制还是先达到时间限制？)
	if l.maxSize > 0 {
		h, err := handlerRotateFile(l, path)
		if err != nil {
			return nil, err
		}

		h.SetFormatter(logFormatter)
		l.sl.AddHandler(h)
		l.rotateWriter = h.Writer()
	} else if l.rotateInterval != "" {
		h, err := handlerRotateTime(l, path)
		if err != nil {
			return nil, err
		}

		h.SetFormatter(logFormatter)
		l.sl.AddHandler(h)
		l.rotateWriter = h.Writer()
	}

	l.sl.DoNothingOnPanicFatal()

	// 打印 logger 的内部设置
	if l.printAfterInitialized {
		l.Info("[logger] Level: %s", l.level)
		l.Info("[logger] Max Backups: %d", l.maxBackups)
		l.Info("[logger] Max Size: %d", l.maxSize)
		l.Info("[logger] Rotate Interval: %s", l.rotateInterval)
		l.Info("[logger] Compress: %t", l.compress)
		l.Info("[logger] Buffer Size: %d", l.bufferSize)
		l.Info("[logger] File Perm: %o", l.filePerm)
		l.Info("[logger] Time Format: %s", l.timeFormat)
	}

	return l, nil
}

// genLogFormatter 构建日志格式
func genLogFormatter(timeFormat string) *slog.TextFormatter {
	logTemplate := "[{{datetime}}] [{{level}}] [{{caller}}] {{message}}\n"
	// custom log format
	logFormatter := slog.NewTextFormatter(logTemplate)
	logFormatter.EnableColor = true
	logFormatter.FullDisplay = true
	logFormatter.TimeFormat = timeFormat

	return logFormatter
}

// parseLevels 解析日志级别
func parseLevels(level string) (levels []slog.Level) {
	switch strings.ToLower(level) {
	case "debug":
		levels = append(levels, slog.InfoLevel, slog.WarnLevel, slog.ErrorLevel, slog.DebugLevel)
	default:
		levels = append(levels, slog.InfoLevel, slog.WarnLevel, slog.ErrorLevel)
	}
	return
}

// handlerRotateFile 文件大小轮转handler
func handlerRotateFile(log *logger, logFile string) (*handler.SyncCloseHandler, error) {
	return handler.NewSizeRotateFileHandler(
		logFile,
		log.maxSize,
		handler.WithLogLevels(parseLevels(log.level)),
		handler.WithBuffSize(log.bufferSize),
		handler.WithBackupNum(log.maxBackups),
		handler.WithCompress(log.compress),
		handler.WithFilePerm(log.filePerm),
	)
}

// handlerRotateTime
// rotateInterval: 1w, 1d, 1h, 1m, 1s
func handlerRotateTime(log *logger, logFile string) (*handler.SyncCloseHandler, error) {
	// 参数检查
	if len(log.rotateInterval) < 2 {
		return nil, fmt.Errorf("invalid rotate interval: %s", log.rotateInterval)
	}

	// 时间转换
	lastChar := log.rotateInterval[len(log.rotateInterval)-1]
	lowerLastChar := strings.ToLower(string(lastChar))
	switch lowerLastChar {
	case "w", "d":
		// time.ParseDuration() 不支持 w、d，因此需要转换成 h。
		prefix, err := strconv.Atoi(log.rotateInterval[:len(log.rotateInterval)-1])
		if err != nil {
			return nil, err
		}

		if lowerLastChar == "w" {
			log.rotateInterval = fmt.Sprintf("%dh", prefix*7*24)
		} else {
			log.rotateInterval = fmt.Sprintf("%dh", prefix*24)
		}
	case "h", "m", "s":
		break
	default:
		return nil, fmt.Errorf("unsuppored rotate interval type: %s", lowerLastChar)
	}

	rotateIntervalDuration, err := time.ParseDuration(log.rotateInterval)
	if err != nil {
		return nil, err
	}

	// 0700 权限更严格，只允许文件所有者访问和修改文件，而 0664 允许组成员修改文件，且所有人都能读取文件。
	rotatefile.DefaultFilePerm = 0700
	return handler.NewTimeRotateFileHandler(
		logFile,
		rotatefile.RotateTime(rotateIntervalDuration.Seconds()),
		handler.WithLogLevels(parseLevels(log.level)),
		handler.WithBuffSize(log.bufferSize),
		handler.WithBackupNum(log.maxBackups),
		handler.WithCompress(log.compress),
		handler.WithFilePerm(log.filePerm),
	)
}

// 实现 Logger 接口。
func (l *logger) Info(msg string, args ...interface{}) {
	l.sl.Infof(msg, args...)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.sl.Errorf(msg, args...)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.sl.Warnf(msg, args...)
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.sl.Debugf(msg, args...)
}

// Write 仅支持文件输出模式
func (l *logger) Write(p []byte) (int, error) {
	if l.rotateWriter == nil {
		return 0, nil
	}

	return l.rotateWriter.Write(p)
}
