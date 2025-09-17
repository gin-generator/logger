package logger

import (
	"encoding/json"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
	DEBUG = "debug"
)

type Logger struct {
	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory.
	fileName string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	maxSize int

	// MaxBackups is the maximum number of old log files to retain. The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted).
	maxBackup int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename. Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	maxAge int

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	compress bool

	// The default values for log levels are: info, warn, error, and debug
	level string

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	localTime bool

	once *sync.Once
	Log  *zap.Logger
}

func newDefaultLogger() *Logger {
	return &Logger{
		fileName:  "logs/logs.log",
		maxSize:   100,
		maxBackup: 7,
		maxAge:    30,
		compress:  false,
		level:     INFO,
		localTime: true,
		once:      new(sync.Once),
	}
}

func (l *Logger) custom() {
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(l.level)); err != nil {
		fmt.Println("Logger init error: invalid log level. Please check the level setting.")
		return
	}

	l.once.Do(func() {
		l.Log = zap.New(
			zapcore.NewCore(l.getEncoder(), l.getWriter(), level),
			zap.AddCaller(),                   // Add caller file and line number, internally uses runtime.Caller
			zap.AddCallerSkip(1),              // Skip one layer of caller file (runtime.Caller(1))
			zap.AddStacktrace(zap.ErrorLevel), // Show stacktrace only for Error level
		)
	})
}

func (l *Logger) getWriter() zapcore.WriteSyncer {
	logName := time.Now().Format("logs-2006-01-02.log")
	filename := strings.ReplaceAll(l.fileName, "logs.log", logName)

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    l.maxSize,
			MaxBackups: l.maxBackup,
			MaxAge:     l.maxAge,
			Compress:   l.compress,
			LocalTime:  l.localTime,
		}),
	)
}

func (l *Logger) getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,      // Add "\n" at the end of each log line
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // Log level names in uppercase, e.g., ERROR, INFO
		EncodeTime:     customTimeEncoder,              // Custom time format: 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder, // Execution time in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Short format for Caller, e.g., types/converter.go:17
	})
}

// customTimeEncoder defines a custom-friendly time format
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// Dump is used for debugging. It does not interrupt the program and prints a warning message to the terminal.
// The first parameter is rendered using json.Marshal, and the second parameter is an optional message.
//
//	logger.Dump(user.User{Name:"test"})
//	logger.Dump(user.User{Name:"test"}, "User information")
func (l *Logger) Dump(value interface{}, message ...string) {
	valueString := l.jsonString(value)
	// Check if the second parameter message is passed
	if len(message) > 0 {
		l.Log.Warn("Dump", zap.String(message[0], valueString))
	} else {
		l.Log.Warn("Dump", zap.String("data", valueString))
	}
}

// LogIf logs an error-level message if err != nil
func (l *Logger) LogIf(err error) {
	if err != nil {
		l.Log.Error("Error Occurred:", zap.Error(err))
	}
}

// LogWarnIf logs a warning-level message if err != nil
func (l *Logger) LogWarnIf(err error) {
	if err != nil {
		l.Log.Warn("Error Occurred:", zap.Error(err))
	}
}

// LogInfoIf logs an info-level message if err != nil
func (l *Logger) LogInfoIf(err error) {
	if err != nil {
		l.Log.Info("Error Occurred:", zap.Error(err))
	}
}

// Debug logs detailed debug messages
// Example usage:
//
//	logger.Debug("Database", zap.String("sql", sql))
func (l *Logger) Debug(moduleName string, fields ...zap.Field) {
	l.Log.Debug(moduleName, fields...)
}

// Info logs informational messages
func (l *Logger) Info(moduleName string, fields ...zap.Field) {
	l.Log.Info(moduleName, fields...)
}

// Warn logs warning messages
func (l *Logger) Warn(moduleName string, fields ...zap.Field) {
	l.Log.Warn(moduleName, fields...)
}

// Error logs error messages without interrupting the program. Focus on these logs when reviewing.
func (l *Logger) Error(moduleName string, fields ...zap.Field) {
	l.Log.Error(moduleName, fields...)
}

// Fatal logs messages at the same level as Error() and exits the program with os.Exit(1)
func (l *Logger) Fatal(moduleName string, fields ...zap.Field) {
	l.Log.Fatal(moduleName, fields...)
}

// DebugString logs a debug message with a string value. Example usage:
//
//	logger.DebugString("SMS", "Message content", string(result.RawResponse))
func (l *Logger) DebugString(moduleName, name, msg string) {
	l.Log.Debug(moduleName, zap.String(name, msg))
}

func (l *Logger) InfoString(moduleName, name, msg string) {
	l.Log.Info(moduleName, zap.String(name, msg))
}

func (l *Logger) WarnString(moduleName, name, msg string) {
	l.Log.Warn(moduleName, zap.String(name, msg))
}

func (l *Logger) ErrorString(moduleName, name, msg string) {
	l.Log.Error(moduleName, zap.String(name, msg))
}

func (l *Logger) FatalString(moduleName, name, msg string) {
	l.Log.Fatal(moduleName, zap.String(name, msg))
}

// DebugJSON logs a debug message with an object value, using json.Marshal for encoding. Example usage:
//
//	logger.DebugJSON("Auth", "Read logged-in user", auth.CurrentUser())
func (l *Logger) DebugJSON(moduleName, name string, value interface{}) {
	l.Log.Debug(moduleName, zap.String(name, l.jsonString(value)))
}

func (l *Logger) InfoJSON(moduleName, name string, value interface{}) {
	l.Log.Info(moduleName, zap.String(name, l.jsonString(value)))
}

func (l *Logger) WarnJSON(moduleName, name string, value interface{}) {
	l.Log.Warn(moduleName, zap.String(name, l.jsonString(value)))
}

func (l *Logger) ErrorJSON(moduleName, name string, value interface{}) {
	l.Log.Error(moduleName, zap.String(name, l.jsonString(value)))
}

func (l *Logger) FatalJSON(moduleName, name string, value interface{}) {
	l.Log.Fatal(moduleName, zap.String(name, l.jsonString(value)))
}

func (l *Logger) jsonString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		l.Log.Error("Logger", zap.String("JSON marshal error", err.Error()))
	}
	return string(b)
}
