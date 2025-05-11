# slog 库引入 QuickStart

`slog`是一个简单灵活的Go日志库，提供了标准输出和文件输出支持，并具有丰富的配置选项。

## 特性

- 支持标准输出和文件输出
- 支持日志级别：Debug、Info、Warn、Error、Fatal
- 支持按大小轮转日志文件
- 支持按时间轮转日志文件（秒、分、时、天、周）
- 支持日志文件压缩
- 支持自定义时间格式
- 使用Option模式，配置灵活
- **支持全局单例模式，无需管理日志实例**

## 安装

```bash
go get github.com/hildam/logs/slog
```

## 基本使用

### 导入包

```go
import "github.com/hildam/logs/slog"
```

## 全局单例模式（推荐）

使用全局单例模式无需管理日志实例，直接初始化后调用包级函数：

```go
package main

import (
    "github.com/hildam/logs/slog"
)

func main() {
    // 初始化全局标准输出日志
    err := slog.Init(
        slog.WithLevel("debug"),
        slog.WithTimeFormat(time.RFC3339),
    )
    if err != nil {
        panic(err)
    }
    
    // 直接使用包级函数
    slog.Info("这是一条信息日志")
    slog.Debug("这是一条调试日志")
    slog.Warn("这是一条警告日志: %s", "警告内容")
    slog.Error("这是一条错误日志: %v", map[string]string{"error": "发生错误"})
    slog.Fatal("这是一条致命错误日志: %v", "程序将终止")
    
    // 也可以使用io.Writer接口
    slog.Write([]byte("直接写入内容"))
}
```

### 初始化文件日志

```go
package main

import (
    "github.com/hildam/logs/slog"
)

func main() {
    // 初始化全局文件日志
    err := slog.InitFile("./app.log",
        slog.WithLevel("info"),
        slog.WithMaxSize(100*1024*1024),
        slog.WithMaxBackups(10),
        slog.WithCompress(),
    )
    if err != nil {
        panic(err)
    }
    
    // 直接使用包级函数
    slog.Info("这是写入文件的日志")
}
```

### 使用已创建的Logger实例

您可以使用`slog.UseLogger`函数将已经创建的Logger实例设置为全局Logger。这个函数会重置初始化状态，允许重新初始化全局Logger。

```go
package main

import (
    "github.com/hildam/logs/slog"
)

func main() {
    // 创建一个文件Logger实例
    fileLogger, err := slog.NewFileLogger("./app.log",
        slog.WithLevel("debug"),
        slog.WithRotateInterval("1d"),
    )
    if err != nil {
        panic(err)
    }
    
    // 使用该实例作为全局Logger
    slog.UseLogger(fileLogger)
    
    // 现在可以直接使用包级函数，它们会调用fileLogger
    slog.Info("这条日志会写入app.log文件")
    
    // 也可以直接使用fileLogger
    fileLogger.Debug("这条日志也会写入app.log文件")
    
    // 稍后如果需要，可以切换到标准输出
    stdLogger, err := slog.NewStdoutLogger()
    if err != nil {
        panic(err)
    }
    
    // 切换全局Logger
    slog.UseLogger(stdLogger)
    
    // 此时包级函数会输出到标准输出
    slog.Info("这条日志会输出到控制台")
}
```

### 高级用法：设置自定义Logger实例

```go
package main

import (
    "github.com/hildam/logs/slog"
)

func main() {
    // 创建自定义Logger
    customLogger, err := slog.NewFileLogger("./custom.log",
        slog.WithRotateInterval("1d"),
    )
    if err != nil {
        panic(err)
    }
    
    // 设置为全局Logger
    slog.SetLogger(customLogger)
    
    // 使用包级函数（实际会调用customLogger）
    slog.Info("这条日志会写入custom.log")
}
```

### 获取当前全局Logger

```go
// 获取当前的全局Logger实例
logger := slog.GetLogger()
```

> **注意**：如果尚未初始化全局Logger，`GetLogger()`会尝试使用默认参数创建一个标准输出Logger。如果这一步也失败，则会panic。这样设计是为了确保日志系统问题能够被及时发现并处理，而不是被默默忽略。

### 创建标准输出日志（实例模式）

```go
package main

import (
    "github.com/hildam/logs/slog"
    "time"
)

func main() {
    // 创建默认配置的标准输出日志
    logger, err := slog.NewStdoutLogger()
    if err != nil {
        panic(err)
    }
    
    // 记录不同级别的日志
    logger.Info("这是一条信息日志")
    logger.Warn("这是一条警告日志: %s", "警告内容")
    logger.Error("这是一条错误日志: %v", map[string]string{"error": "发生错误"})
}
```

### 创建文件输出日志（实例模式）

#### 按大小轮转

```go
// 创建按文件大小轮转的日志
logger, err := slog.NewFileLogger("./app.log",
    slog.WithLevel("info"),            // 设置日志级别
    slog.WithMaxSize(100*1024*1024),   // 设置单个日志文件最大大小（100MB）
    slog.WithMaxBackups(10),           // 设置最大保留日志文件数量
    slog.WithCompress(),               // 启用日志文件压缩
)
if err != nil {
    panic(err)
}

logger.Info("这是写入文件的日志")
```

#### 按时间轮转

```go
// 创建按时间轮转的日志（每天轮转一次）
logger, err := slog.NewFileLogger("./time_rotate.log",
    slog.WithRotateInterval("1d"),     // 每天轮转
    slog.WithMaxBackups(30),           // 保留30个备份
)
if err != nil {
    panic(err)
}

logger.Info("这是按时间轮转的日志")
```

### 使用io.Writer接口

slog库实现了`io.Writer`接口，可以直接用于需要Writer的场景：

```go
logger, _ := slog.NewFileLogger("./app.log")
logger.Write([]byte("直接写入的内容"))
```

## 配置选项

`slog`库采用Option模式进行配置，提供以下选项：

### 日志级别

```go
// 设置为debug级别（记录所有级别的日志）
slog.WithLevel("debug")

// 设置为info级别（仅记录info、warn、error级别的日志）
slog.WithLevel("info")
```

### 文件大小

```go
// 设置日志文件最大大小为50MB
slog.WithMaxSize(50*1024*1024)
```

### 时间轮转

```go
// 每小时轮转一次
slog.WithRotateInterval("1h")

// 每周轮转一次
slog.WithRotateInterval("1w")

// 每天轮转一次
slog.WithRotateInterval("1d")

// 每5分钟轮转一次
slog.WithRotateInterval("5m")

// 每30秒轮转一次
slog.WithRotateInterval("30s")
```

### 日志备份

```go
// 保留最近5个备份文件
slog.WithMaxBackups(5)
```

### 缓冲区大小

```go
// 设置缓冲区大小为4KB
slog.WithBufferSize(4*1024)

// 设置为0，立即写入文件
slog.WithBufferSize(0)
```

### 时间格式

```go
// 使用RFC3339格式
slog.WithTimeFormat(time.RFC3339)

// 使用自定义格式
slog.WithTimeFormat("2006-01-02 15:04:05.000")
```

### 文件权限

```go
// 设置日志文件权限
slog.WithFilePerm(0644)
```

### 打印初始化信息

```go
// 启用初始化后打印配置信息
slog.WithPrintAfterInitialized()
```

### 颜色输出

```go
// 启用颜色输出（默认为关闭）
slog.WithColor(true)

// 关闭颜色输出
slog.WithColor(false)
```

### 调用者位置跟踪

控制日志中显示的调用者位置信息：

```go
// 默认为1，显示的是调用日志函数的位置（用户代码）
slog.WithCallerSkip(1)

// 设为0，显示slog库内部调用位置
slog.WithCallerSkip(0)

// 设为2或更高，跳过更多调用层，显示更上层的调用者
slog.WithCallerSkip(2)
```

## 完整示例

### 全局单例模式示例

```go
package main

import (
    "github.com/hildam/logs/slog"
    "time"
)

func main() {
    // 初始化全局Logger
    err := slog.Init(
        slog.WithLevel("debug"),
        slog.WithTimeFormat(time.RFC3339),
    )
    if err != nil {
        panic(err)
    }
    
    // 在主函数中记录日志
    slog.Info("应用启动")
    
    // 调用其他函数
    processData()
}

func processData() {
    // 在其他函数中直接使用，无需传递Logger实例
    slog.Debug("开始处理数据")
    // 处理逻辑...
    slog.Info("数据处理完成")
}
```

### 实例模式示例

```go
package main

import (
    "github.com/hildam/logs/slog"
    "time"
)

func main() {
    // 创建标准输出日志（debug级别）
    stdLogger, err := slog.NewStdoutLogger(
        slog.WithLevel("debug"),
        slog.WithTimeFormat(time.RFC3339),
    )
    if err != nil {
        panic(err)
    }
    
    stdLogger.Info("标准输出日志示例")
    stdLogger.Debug("这条debug日志会被记录")
    
    // 创建文件日志（按大小轮转）
    fileLogger, err := slog.NewFileLogger("./logs/app.log",
        slog.WithLevel("info"),
        slog.WithMaxSize(50*1024*1024),
        slog.WithMaxBackups(10),
        slog.WithCompress(),
        slog.WithBufferSize(8*1024),
        slog.WithFilePerm(0644),
        slog.WithPrintAfterInitialized(),
    )
    if err != nil {
        panic(err)
    }
    
    fileLogger.Info("文件日志示例: %s", "写入成功")
    fileLogger.Debug("这条debug日志不会被记录，因为级别设置为info")
    
    // 创建文件日志（按时间轮转）
    timeLogger, err := slog.NewFileLogger("./logs/time.log",
        slog.WithLevel("info"),
        slog.WithRotateInterval("1d"),
        slog.WithMaxBackups(30),
    )
    if err != nil {
        panic(err)
    }
    
    timeLogger.Info("按时间轮转的日志示例")
}
```

## 注意事项

1. 文件大小轮转和时间轮转是两种互斥的模式，只能选择其一
2. 默认的日志级别是`info`
3. 默认的轮转间隔是`1w`（一周）
4. 默认的最大文件大小是100MB
5. 默认启用日志压缩
6. 默认的文件权限是0700
7. **默认关闭颜色输出，可通过WithColor选项开启**
8. **默认的调用栈跳过级别是1，会显示用户代码的调用位置**
9. **全局单例模式下，`Init`和`InitFile`只有第一次调用会生效**
10. **如果日志初始化失败（如权限问题、磁盘已满等），系统会直接panic，而非忽略错误**
11. **`UseLogger`函数会重置初始化状态，允许后续再次调用`Init`或`InitFile`** 