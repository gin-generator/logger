package logger

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger implements gorm.logger.Interface for GORM integration.
// It logs SQL queries, errors, and slow queries using the underlying Logger.
type GormLogger struct {
	logger        *Logger
	slowThreshold time.Duration
}

// NewGormLogger creates a new GORM logger with the specified slow query threshold.
// If slowThreshold is 0, it defaults to 200ms.
func NewGormLogger(logger *Logger, slowThreshold time.Duration) *GormLogger {
	if slowThreshold == 0 {
		slowThreshold = 200 * time.Millisecond
	}
	return &GormLogger{
		logger:        logger,
		slowThreshold: slowThreshold,
	}
}

// LogMode implements gorm.logger.Interface.
// It returns the logger itself as log level filtering is handled by zap.
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info logs informational messages from GORM at debug level.
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger().WithContext(ctx).Sugar().Debugf(msg, args...)
}

// Warn logs warning messages from GORM.
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger().WithContext(ctx).Sugar().Warnf(msg, args...)
}

// Error logs error messages from GORM.
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger().WithContext(ctx).Sugar().Errorf(msg, args...)
}

// Trace logs SQL execution details including query, elapsed time, and affected rows.
// It logs errors at error level, slow queries at warn level, and normal queries at debug level.
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	logger := l.zapLogger().WithContext(ctx)

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		logger.Error("gorm", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold:
		logger.Warn("gorm slow query", fields...)
	default:
		logger.Debug("gorm", fields...)
	}
}

// zapLogger adjusts the caller skip to report the correct file and line number
// in logs, skipping internal GORM frames.
func (l *GormLogger) zapLogger() *Logger {
	clone := l.logger.WithOptions(zap.AddCallerSkip(-1))
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok || strings.HasSuffix(file, "_test.go") {
			break
		}
		if strings.Contains(file, filepath.Join("gorm.io", "gorm")) {
			continue
		}
		return &Logger{Logger: clone.WithOptions(zap.AddCallerSkip(i))}
	}
	return l.logger
}
