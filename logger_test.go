package logger

import (
	"context"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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
	ctx := WithTraceID(context.Background(), "trace-123")

	log.WithContext(ctx).Info("test with trace id")
}

func TestTraceIDFromContext(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")
	if traceID := TraceIDFromContext(ctx); traceID != "trace-123" {
		t.Fatalf("expected trace-123, got %q", traceID)
	}
	if traceID := TraceIDFromContext(nil); traceID != "" {
		t.Fatalf("expected empty trace ID, got %q", traceID)
	}
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
	core, observed := observer.New(zapcore.DebugLevel)
	log := &Logger{Logger: zap.New(core)}
	gormLog := NewGormLogger(log, 100*time.Millisecond)

	ctx := WithTraceID(context.Background(), "trace-gorm-123")
	gormLog.Info(ctx, "test info with trace: %s", "value")
	gormLog.Warn(ctx, "test warn with trace")
	gormLog.Error(ctx, "test error with trace")
	gormLog.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT * FROM users WHERE id = 1", 1
	}, nil)

	entries := observed.All()
	if len(entries) != 4 {
		t.Fatalf("expected 4 log entries, got %d", len(entries))
	}
	for _, entry := range entries {
		if traceID, ok := entry.ContextMap()["trace_id"]; !ok || traceID != "trace-gorm-123" {
			t.Fatalf("expected trace_id on %q, got %#v", entry.Message, entry.ContextMap())
		}
	}
	elapsed, ok := entries[3].ContextMap()["elapsed"].(string)
	if !ok || !strings.HasSuffix(elapsed, "ms") {
		t.Fatalf("expected elapsed in milliseconds, got %#v", entries[3].ContextMap()["elapsed"])
	}
}
