package logger

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGormLogger(t *testing.T) {
	var dbConfig gorm.Dialector
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		"root",
		"admin@qcz",
		"127.0.0.1",
		"3306",
		"weekend",
		"utf8mb4",
	)
	dbConfig = mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
	})

	logger := NewLogger(
		WithFileName("logs/sql.log"),
		WithLevel(DEBUG),
		WithTimeZone(true),
	)

	_logger := NewGormLogger(logger, WithSlowThreshold(300*time.Millisecond))
	db, err := gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})

	if err != nil {
		t.Fatal("database connection failure")
		return
	}

	err = db.Exec("show tables;").Error
	if err != nil {
		t.Fatal("database query failure")
		return
	}
}
