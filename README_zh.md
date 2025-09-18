# logger
[![Go Reference](https://img.shields.io/github/v/release/gin-generator/logger.svg?style=flat-square)](https://pkg.go.dev/github.com/gin-generator/logger)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](README.md)

# 日志包

一个基于 [Uber zap](https://github.com/uber-go/zap) 和 [lumberjack](https://github.com/natefinch/lumberjack) 的高性能日志库，支持日志切割、压缩、结构化日志输出、可配置项丰富，适用于生产环境。

## 功能特色

- 📁 日志文件路径、大小、备份、保存天数可自定义
- 🔄 支持日志自动切割与压缩
- 📋 JSON 格式结构化日志
- 🔍 支持多级日志：Debug、Info、Warn、Error、Fatal
- ✅ 自带单元测试，保证稳定性


## 安装

```bash
go get github.com/gin-generator/logger
```

## 常用方法
* Dump(interface{}, ...string) – 快速打印调试信息
* LogIf(error) – 错误非空时记录错误日志
* 支持：DebugString, InfoString, WarnString, ErrorString, FatalString
* 支持：DebugJSON, InfoJSON, WarnJSON, ErrorJSON, FatalJSON

## 使用
### 示例
[测试用例](logger_test.go)
[自定义Gorm logger](gorm_test.go)

## License
MIT © [Chaozheng Qin](LICENSE)