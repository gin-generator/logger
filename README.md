# logger

[中文文档](README_zh.md)

A structured logging library for Go, built on top of [Uber Zap](https://github.com/uber-go/zap) and [Lumberjack](https://github.com/natefinch/lumberjack). Provides file rotation, trace ID propagation, and GORM integration out of the box.

## Features

- JSON structured logging via Zap
- Automatic log file rotation (size, backup count, age)
- Optional gzip compression of rotated files
- Trace ID propagation through `context.Context`
- GORM integration with slow query detection
- Functional options pattern for clean configuration

## Installation

```bash
go get github.com/gin-generator/logger
```

## Quick Start

```go
import "github.com/gin-generator/logger"

// Console output only (default)
log := logger.New()
log.Info("server started", zap.String("addr", ":8080"))

// Write to file with rotation
log := logger.New(
    logger.WithFile("logs/app.log"),
    logger.WithLevel("debug"),
)
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithLevel(string)` | Set log level: `debug`, `info`, `warn`, `error`, `dpanic`, `panic`, `fatal` | `info` |
| `WithLevelValue(zapcore.Level)` | Set log level using a `zapcore.Level` constant | `zapcore.InfoLevel` |
| `WithFile(filename)` | Enable file output at the given path | disabled |
| `WithRotation(maxSize, maxBackups, maxAge)` | File rotation: max MB, max backup files, max days | `100, 7, 30` |
| `WithCompress(bool)` | Gzip-compress rotated log files | `false` |
| `WithConsole(bool)` | Enable or disable stdout output | `true` |
| `WithEncoder(zapcore.Encoder)` | Override the log encoder | JSON encoder |

`WithRotation` must be called **before** `WithFile` to take effect.

## Usage Examples

### File output with rotation

```go
log := logger.New(
    logger.WithLevel("debug"),
    logger.WithRotation(50, 5, 10), // 50 MB, 5 backups, 10 days
    logger.WithCompress(true),
    logger.WithFile("logs/app.log"),
)
```

### File only (no console)

```go
log := logger.New(
    logger.WithConsole(false),
    logger.WithFile("logs/app.log"),
)
```

### Trace ID via context

Set `logger.TraceIDKey` in the context and call `WithContext` to attach the trace ID to every log entry automatically.

```go
ctx := context.WithValue(req.Context(), logger.TraceIDKey, "abc-123")
log.WithContext(ctx).Info("request received", zap.String("path", "/api/users"))
// Output: {"level":"info","message":"request received","trace_id":"abc-123","path":"/api/users",...}
```

### Custom encoder

```go
consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
log := logger.New(logger.WithEncoder(consoleEncoder))
```

## GORM Integration

`NewGormLogger` wraps the logger to implement `gorm.logger.Interface`. It logs:

- SQL errors at `error` level (excluding `ErrRecordNotFound`)
- Slow queries at `warn` level
- Normal queries at `debug` level

```go
import (
    "time"
    "github.com/gin-generator/logger"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

log := logger.New(logger.WithLevel("debug"))
gormLog := logger.NewGormLogger(log, 100*time.Millisecond) // slow threshold: 100ms

db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
    Logger: gormLog,
})
```

Pass `0` as the slow threshold to use the default of 200ms.

## Log Output Format

Default JSON output:

```json
{"time":"2026-05-08 12:00:00","level":"info","caller":"main.go:12","message":"server started","addr":":8080"}
```

With trace ID:

```json
{"time":"2026-05-08 12:00:01","level":"info","caller":"handler.go:34","message":"request received","trace_id":"abc-123","path":"/api/users"}
```

## License

MIT
