# logger

[![Go Reference](https://img.shields.io/github/v/release/gin-generator/logger.svg?style=flat-square)](https://pkg.go.dev/github.com/gin-generator/logger)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[中文](README_zh.md)

A customizable, production-ready logging package based on [Uber's zap](https://github.com/uber-go/zap)
and [lumberjack](https://github.com/natefinch/lumberjack). This package supports log rotation, JSON formatting,
structured logging, and various log levels with minimal setup.

## Features

- ⚙️ Configurable log file path, size, age, and backups
- 📦 Automatic log rotation with compression support
- 📄 JSON-formatted structured logs
- 🧠 Level-based logging: Debug, Info, Warn, Error, Fatal
- 🧪 Comprehensive unit tests included

## Installation

```bash
go get github.com/gin-generator/logger
```

## Utility Methods

* Dump(interface{}, ...string) – Quick debug print using json.Marshal
* LogIf(error) – Logs an error if not nil
* DebugString, InfoString, WarnString, ErrorString, FatalString
* DebugJSON, InfoJSON, WarnJSON, ErrorJSON, FatalJSON

## Usage

### Example

[Test case](logger_test.go)
[Custom gorm logger](gorm_test.go)

## License

MIT © [Chaozheng Qin](LICENSE)