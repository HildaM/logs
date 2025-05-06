package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gookit/slog"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger slog gorm logger实现
type GormLogger struct {
	logger.Config
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info log info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel <= logger.Silent {
		return
	}
	slog.Infof(utils.FileWithLineNum()+" "+msg, data...)
}

// Warn log warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel <= logger.Silent {
		return
	}
	slog.Warnf(utils.FileWithLineNum()+" "+msg, data...)
}

// Error log error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel <= logger.Silent {
		return
	}
	slog.Errorf(utils.FileWithLineNum()+" "+msg, data...)
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 如果日志级别低于或等于Silent级别(即不输出任何日志)，则直接返回
	if l.LogLevel <= logger.Silent {
		return
	}

	// 计算SQL执行耗时
	elapsed := time.Since(begin)
	costTime := float64(elapsed.Nanoseconds()) / 1e6 // 将纳秒转换为毫秒

	switch {
	// 错误日志处理: 当有错误发生且日志级别>=Error且(错误不是RecordNotFound或配置了不忽略此类错误)
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		// 通过回调函数获取SQL语句和影响的行数
		sql, rows := fc()
		if rows == -1 {
			rows = 0 // 标准化行数，-1表示不适用的情况
		}
		slog.Errorf(utils.FileWithLineNum()+"error=%s [cost:%.3fms] [rows:%v] %s", err, costTime, rows, sql)

	// 慢查询日志处理: 当SQL执行时间超过阈值且阈值不为0且日志级别>=Warn
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			rows = 0
		}
		slog.Warnf(utils.FileWithLineNum()+"slow=%s [cost:%.3fms] [rows:%v] %s", slowLog, costTime, rows, sql)

	// 普通信息日志: 当日志级别为Info时记录所有SQL
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			rows = 0
		}
		// 使用Trace级别记录详细SQL信息，包含执行位置、耗时、影响行数和SQL语句
		slog.Tracef(utils.FileWithLineNum()+"[cost:%.3fms] [rows:%v] %s", costTime, rows, sql)
	}
}
