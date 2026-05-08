package logger

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	log := New()
	if log == nil {
		t.Fatal("expected logger instance")
	}
	log.Info("test message", zap.String("key", "value"))
}

func TestWithFile(t *testing.T) {
	log := New(
		WithFile("logs/test.log"),
		WithLevel("debug"),
	)
	log.Debug("debug message")
	log.Info("info message")
}

func TestWithRotation(t *testing.T) {
	log := New(
		WithFile("logs/test.log"),
		WithRotation(50, 5, 10),
		WithCompress(true),
	)
	log.Info("test rotation")
}

func TestWithContext(t *testing.T) {
	log := New()
	ctx := context.WithValue(context.Background(), TraceIDKey, "trace-123")

	log.WithContext(ctx).Info("test with trace id")
}

func TestWithLevelValue(t *testing.T) {
	log := New(WithLevelValue(zapcore.DebugLevel))
	log.Debug("debug message with level constant")
	log.Info("info message")
}

func TestGormLogger(t *testing.T) {
	log := New(WithLevel("debug"))
	gormLog := NewGormLogger(log, 100*time.Millisecond)

	ctx := context.Background()
	gormLog.Info(ctx, "test info: %s", "value")
	gormLog.Warn(ctx, "test warn")
	gormLog.Error(ctx, "test error")
}

func TestGormLoggerWithTraceID(t *testing.T) {
	log := New(
		WithLevel("debug"),
		WithFile("logs/test.log"),
	)
	gormLog := NewGormLogger(log, 100*time.Millisecond)

	ctx := context.WithValue(context.Background(), TraceIDKey, "trace-gorm-123")
	gormLog.Info(ctx, "test info with trace: %s", "value")
	gormLog.Warn(ctx, "test warn with trace")
	gormLog.Error(ctx, "test error with trace")
	gormLog.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users WHERE id = 1", 1
	}, nil)
}
