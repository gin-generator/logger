// Package logger provides a structured logging solution based on zap and lumberjack.
// It supports file rotation, trace ID propagation, and GORM integration.
package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger and provides additional functionality like context-aware logging.
type Logger struct {
	*zap.Logger
}

// WithTraceID returns a child context containing traceID.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// TraceIDFromContext returns the trace ID stored in ctx.
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceID, _ := ctx.Value(TraceIDKey).(string)
	return traceID
}

// New creates a new Logger instance with the given options.
// If no options are provided, it uses default configuration.
func New(opts ...Option) *Logger {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	core := zapcore.NewCore(
		cfg.encoder,
		zapcore.NewMultiWriteSyncer(cfg.outputs...),
		cfg.level,
	)

	return &Logger{
		Logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)),
	}
}

// WithContext returns a zap.Logger with trace ID from context if available.
// The trace ID is extracted using TraceIDKey and added as a field to all logs.
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		return l.With(zap.String("trace_id", traceID))
	}
	return l.Logger
}
