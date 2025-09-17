package logger

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	l "gorm.io/gorm/logger"
)

type GormOption interface {
	apply(*GormLogger)
}

type gormOptionFunc func(*GormLogger)

func (g gormOptionFunc) apply(l *GormLogger) {
	g(l)
}

type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

// NewGormLogger is called externally. It instantiates a GormLogger object. Example:
func NewGormLogger(logger *Logger, opts ...GormOption) *GormLogger {
	log := &GormLogger{
		ZapLogger:     logger.Log,
		SlowThreshold: 200 * time.Millisecond, // Slow query threshold in milliseconds
	}

	for _, opt := range opts {
		opt.apply(log)
	}
	return log
}

// WithSlowThreshold sets the slow query threshold
func WithSlowThreshold(times time.Duration) GormOption {
	return gormOptionFunc(func(l *GormLogger) {
		l.SlowThreshold = times
	})
}

// LogMode sets the logging level
func (l *GormLogger) LogMode(level l.LogLevel) l.Interface {
	return &GormLogger{
		ZapLogger:     l.ZapLogger,
		SlowThreshold: l.SlowThreshold,
	}
}

// Info logs informational messages
func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Debugf(str, args...)
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Warnf(str, args...)
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Errorf(str, args...)
}

// Trace logs SQL execution details, including execution time, rows affected, and errors
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", microsecondsStr(elapsed)),
		zap.Int64("rows", rows),
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, zap.Error(err))
			l.logger().Error("Database Error", logFields...)
		}
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger().Warn("Database Slow Log", logFields...)
	}
	l.logger().Debug("Database Query", logFields...)
}

// logger is an internal helper method to ensure the accuracy of the Caller information in Zap (e.g., paginator/paginator.go:148)
func (l *GormLogger) logger() *zap.Logger {
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapGormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	// Subtract one level of wrapping and add one level of zap.AddCallerSkip(1) during initialization
	clone := l.ZapLogger.WithOptions(zap.AddCallerSkip(-2))
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapGormPackage):
		default:
			// Return a new zap logger with adjusted caller skip
			return clone.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}

// microsecondsStr formats the elapsed time in milliseconds
func microsecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}
