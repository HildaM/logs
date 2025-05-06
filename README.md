# Logs

一个简单、灵活的Go日志库，专为小型项目设计。

## 特性

- 支持多种输出方式：标准输出和文件输出
- 支持日志级别：Debug、Info、Warn、Error
- 文件日志支持按大小轮转
- 文件日志支持按时间轮转（支持秒、分、时、天、周）
- 支持日志文件压缩
- 支持自定义时间格式
- 使用Option模式灵活配置

## 安装

```bash
go get github.com/hildam/logs
```

## 使用示例

### 标准输出日志

```go
import "github.com/hildam/logs/slog"

// 创建标准输出日志
logger, err := slog.NewStdoutLogger(
    slog.WithLevel("debug"),                // 设置日志级别
    slog.WithTimeFormat(time.RFC3339),      // 自定义时间格式
)
if err != nil {
    panic(err)
}

// 使用日志
logger.Info("这是一条信息日志: %s", "信息内容")
logger.Debug("这是一条调试日志: %s", "调试内容")
logger.Warn("这是一条警告日志: %s", "警告内容") 
logger.Error("这是一条错误日志: %s", "错误内容")

// 也可以使用io.Writer接口
logger.Write([]byte("直接写入内容"))
```

### 文件输出日志

```go
import "github.com/hildam/logs/slog"

// 创建文件输出日志（按大小轮转）
logger, err := slog.NewFileLogger("./app.log",
    slog.WithLevel("info"),            // 设置日志级别
    slog.WithMaxSize(100*1024*1024),   // 设置单个日志文件最大大小（100MB）
    slog.WithMaxBackups(10),           // 设置最大保留日志文件数量
    slog.WithCompress(),               // 启用日志文件压缩
    slog.WithBufferSize(8*1024),       // 设置缓冲区大小
)
if err != nil {
    panic(err)
}

// 使用日志
logger.Info("文件日志示例: %s", "信息内容")

// 创建文件输出日志（按时间轮转）
loggerTime, err := slog.NewFileLogger("./time_rotate.log",
    slog.WithRotateInterval("1d"),     // 每天轮转
    slog.WithMaxBackups(30),           // 保留30个备份
)
if err != nil {
    panic(err)
}

loggerTime.Info("时间轮转日志示例")
```

## 配置选项

| 选项 | 描述 | 默认值 |
|------|------|--------|
| WithLevel | 设置日志级别 (debug, info, warn, error) | info |
| WithMaxSize | 设置日志文件最大大小（字节） | 100MB |
| WithRotateInterval | 设置日志轮转间隔 (1s, 5m, 2h, 1d, 2w) | 1w |
| WithMaxBackups | 设置最大备份文件数 | 12 |
| WithBufferSize | 设置缓冲区大小（字节，0表示立即写入） | 8KB |
| WithCompress | 启用日志压缩 | 默认启用 |
| WithTimeFormat | 设置时间格式 | time.DateTime |
| WithFilePerm | 设置日志文件权限 | 0700 |
| WithPrintAfterInitialized | 启用初始化后打印日志器配置信息 | 默认禁用 |

## 未来规划

- 支持更多日志库适配
- 添加更多日志格式
- 支持更多输出目标（如网络、数据库等）

## 许可证

MIT 