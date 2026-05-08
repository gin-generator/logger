# logger

[English Documentation](README.md)

基于 [Uber Zap](https://github.com/uber-go/zap) 和 [Lumberjack](https://github.com/natefinch/lumberjack) 封装的 Go 结构化日志库，开箱即用地支持文件轮转、Trace ID 透传和 GORM 集成。

## 特性

- 基于 Zap 的 JSON 结构化日志
- 自动日志文件轮转（按大小、备份数量、保留天数）
- 支持对轮转文件进行 gzip 压缩
- 通过 `context.Context` 透传 Trace ID
- GORM 集成，支持慢查询检测
- 函数式选项模式，配置简洁清晰

## 安装

```bash
go get github.com/gin-generator/logger
```

## 快速开始

```go
import "github.com/gin-generator/logger"

// 仅输出到控制台（默认）
log := logger.New()
log.Info("服务启动", zap.String("addr", ":8080"))

// 输出到文件并自动轮转
log := logger.New(
    logger.WithFile("logs/app.log"),
    logger.WithLevel("debug"),
)
```

## 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithLevel(string)` | 设置日志级别：`debug`、`info`、`warn`、`error`、`dpanic`、`panic`、`fatal` | `info` |
| `WithLevelValue(zapcore.Level)` | 使用 `zapcore.Level` 常量设置日志级别 | `zapcore.InfoLevel` |
| `WithFile(filename)` | 启用文件输出，指定文件路径 | 禁用 |
| `WithRotation(maxSize, maxBackups, maxAge)` | 文件轮转：最大 MB 数、最大备份文件数、最大保留天数 | `100, 7, 30` |
| `WithCompress(bool)` | 对轮转后的日志文件进行 gzip 压缩 | `false` |
| `WithConsole(bool)` | 启用或禁用标准输出 | `true` |
| `WithEncoder(zapcore.Encoder)` | 自定义日志编码器 | JSON 编码器 |

`WithRotation` 必须在 `WithFile` **之前**调用才能生效。

## 使用示例

### 文件输出与轮转

```go
log := logger.New(
    logger.WithLevel("debug"),
    logger.WithRotation(50, 5, 10), // 50 MB，5 个备份，保留 10 天
    logger.WithCompress(true),
    logger.WithFile("logs/app.log"),
)
```

### 仅写文件（关闭控制台输出）

```go
log := logger.New(
    logger.WithConsole(false),
    logger.WithFile("logs/app.log"),
)
```

### 通过 context 透传 Trace ID

在 context 中设置 `logger.TraceIDKey`，调用 `WithContext` 后，每条日志都会自动附带 Trace ID。

```go
ctx := context.WithValue(req.Context(), logger.TraceIDKey, "abc-123")
log.WithContext(ctx).Info("收到请求", zap.String("path", "/api/users"))
// 输出：{"level":"info","message":"收到请求","trace_id":"abc-123","path":"/api/users",...}
```

### 自定义编码器

```go
consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
log := logger.New(logger.WithEncoder(consoleEncoder))
```

## GORM 集成

`NewGormLogger` 将 Logger 包装为实现 `gorm.logger.Interface` 的 GORM 日志器，日志行为如下：

- SQL 错误（排除 `ErrRecordNotFound`）→ `error` 级别
- 慢查询 → `warn` 级别
- 正常查询 → `debug` 级别

```go
import (
    "time"
    "github.com/gin-generator/logger"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

log := logger.New(logger.WithLevel("debug"))
gormLog := logger.NewGormLogger(log, 100*time.Millisecond) // 慢查询阈值：100ms

db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
    Logger: gormLog,
})
```

慢查询阈值传 `0` 时使用默认值 200ms。

## 日志输出格式

默认 JSON 格式：

```json
{"time":"2026-05-08 12:00:00","level":"info","caller":"main.go:12","message":"服务启动","addr":":8080"}
```

携带 Trace ID：

```json
{"time":"2026-05-08 12:00:01","level":"info","caller":"handler.go:34","message":"收到请求","trace_id":"abc-123","path":"/api/users"}
```

## License

MIT
