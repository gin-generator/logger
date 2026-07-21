package logger

import (
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Option is a function that configures the logger.
type Option func(*config)

// config holds the internal configuration for the logger.
type config struct {
	level      zapcore.Level
	encoder    zapcore.Encoder
	outputs    []zapcore.WriteSyncer
	filename   string
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
	localTime  bool
}

// TraceIDKey is the context key used to store and retrieve trace IDs.
const TraceIDKey = "X-Trace-Id"

// defaultConfig returns a config with sensible defaults.
func defaultConfig() *config {
	return &config{
		level:      zapcore.InfoLevel,
		encoder:    defaultEncoder(),
		outputs:    []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)},
		maxSize:    100,
		maxBackups: 7,
		maxAge:     30,
		localTime:  true,
	}
}

// defaultEncoder returns a JSON encoder with standard configuration.
func defaultEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "message",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.DateTime),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

// WithLevel sets the minimum log level from a string.
// Valid levels are: "debug", "info", "warn", "error", "dpanic", "panic", "fatal".
func WithLevel(level string) Option {
	return func(c *config) {
		var l zapcore.Level
		if err := l.UnmarshalText([]byte(level)); err == nil {
			c.level = l
		}
	}
}

// WithLevelValue sets the minimum log level using zapcore.Level constants.
// Example: WithLevelValue(zapcore.DebugLevel)
func WithLevelValue(level zapcore.Level) Option {
	return func(c *config) {
		c.level = level
	}
}

// WithFile enables file output with the specified filename.
// The file will be rotated based on size, backup count, and age settings.
func WithFile(filename string) Option {
	return func(c *config) {
		c.filename = filename
		c.outputs = append(c.outputs, zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    c.maxSize,
			MaxBackups: c.maxBackups,
			MaxAge:     c.maxAge,
			Compress:   c.compress,
			LocalTime:  c.localTime,
		}))
	}
}

// WithRotation configures log file rotation parameters.
// maxSize is the maximum size in megabytes before rotation.
// maxBackups is the maximum number of old log files to retain.
// maxAge is the maximum number of days to retain old log files.
func WithRotation(maxSize, maxBackups, maxAge int) Option {
	return func(c *config) {
		c.maxSize = maxSize
		c.maxBackups = maxBackups
		c.maxAge = maxAge
	}
}

// WithCompress enables or disables gzip compression of rotated log files.
func WithCompress(compress bool) Option {
	return func(c *config) {
		c.compress = compress
	}
}

// WithConsole enables or disables console output.
// When disabled, logs will only be written to file (if configured).
func WithConsole(enable bool) Option {
	return func(c *config) {
		if !enable {
			c.outputs = nil
		}
	}
}

// WithEncoder sets a custom zapcore.Encoder for log formatting.
func WithEncoder(encoder zapcore.Encoder) Option {
	return func(c *config) {
		c.encoder = encoder
	}
}
