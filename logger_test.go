package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

// TestNewLogger tests the creation of a new logger with default settings.
func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("Expected logger instance, got nil")
	}
	logger.Info("Test", zap.String("message", "Logger initialized successfully"))
}

// TestWithFileName tests setting a custom log file name.
func TestWithFileName(t *testing.T) {
	logger := NewLogger(WithFileName("test.log"))
	if logger.fileName != "test.log" {
		t.Fatalf("Expected fileName to be 'test.log', got '%s'", logger.fileName)
	}
	logger.Info("Test", zap.String("message", "Custom file name set"))
}

// TestWithMaxSize tests setting a custom max size for log files.
func TestWithMaxSize(t *testing.T) {
	logger := NewLogger(WithMaxSize(50))
	if logger.maxSize != 50 {
		t.Fatalf("Expected maxSize to be 50, got %d", logger.maxSize)
	}
	logger.Info("Test", zap.String("message", "Custom max size set"))
}

// TestWithMaxBackup tests setting a custom max backup count.
func TestWithMaxBackup(t *testing.T) {
	logger := NewLogger(WithMaxBackup(5))
	if logger.maxBackup != 5 {
		t.Fatalf("Expected maxBackup to be 5, got %d", logger.maxBackup)
	}
	logger.Info("Test", zap.String("message", "Custom max backup set"))
}

// TestWithMaxAge tests setting a custom max age for log files.
func TestWithMaxAge(t *testing.T) {
	logger := NewLogger(WithMaxAge(10))
	if logger.maxAge != 10 {
		t.Fatalf("Expected maxAge to be 10, got %d", logger.maxAge)
	}
	logger.Info("Test", zap.String("message", "Custom max age set"))
}

// TestWithCompress tests enabling log file compression.
func TestWithCompress(t *testing.T) {
	logger := NewLogger(WithCompress(true))
	if !logger.compress {
		t.Fatal("Expected compress to be true, got false")
	}
	logger.Info("Test", zap.String("message", "Log compression enabled"))
}

// TestWithLevel tests setting different log levels.
func TestWithLevel(t *testing.T) {
	logger := NewLogger(WithLevel(DEBUG))
	if logger.level != DEBUG {
		t.Fatalf("Expected level to be 'debug', got '%s'", logger.level)
	}
	logger.Debug("Test", zap.String("message", "Debug level log"))
	logger.Info("Test", zap.String("message", "Info level log"))
	logger.Warn("Test", zap.String("message", "Warn level log"))
	logger.Error("Test", zap.String("message", "Error level log"))
}

// TestDump tests the Dump method for debugging.
func TestDump(t *testing.T) {
	logger := NewLogger()
	testData := map[string]string{"key": "value"}
	logger.Dump(testData, "TestDump")
}

// TestLogIf tests the LogIf method for logging errors.
func TestLogIf(t *testing.T) {
	logger := NewLogger()
	var err error = nil
	logger.LogIf(err) // Should not log anything

	err = os.ErrNotExist
	logger.LogIf(err) // Should log the error
}

// TestLogWarnIf tests the LogWarnIf method for logging warnings.
func TestLogWarnIf(t *testing.T) {
	logger := NewLogger()
	var err error = nil
	logger.LogWarnIf(err) // Should not log anything

	err = os.ErrPermission
	logger.LogWarnIf(err) // Should log the warning
}

// TestJSONLogging tests logging JSON data.
func TestJSONLogging(t *testing.T) {
	logger := NewLogger()
	testData := map[string]string{"key": "value"}
	logger.DebugJSON("TestModule", "TestData", testData)
	logger.InfoJSON("TestModule", "TestData", testData)
	logger.WarnJSON("TestModule", "TestData", testData)
	logger.ErrorJSON("TestModule", "TestData", testData)
}
