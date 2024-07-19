package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/logger"

	"github.com/ethereum/go-ethereum/log"
)

var (
	_ logger.Interface = Logger{}

	SlowThresholdMilliseconds = 200
)

type Logger struct {
	log log.Logger
}

func NewLogger(log log.Logger) Logger {
	return Logger{log.New("module", "db")}
}

func (l Logger) LogMode(lvl logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, data...))
}

func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, data...))
}

func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, data...))
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsedMs := time.Since(begin).Milliseconds()

	// omit any values for batch inserts as they can be very long
	sql, rows := fc()
	if i := strings.Index(strings.ToLower(sql), "values"); i > 0 {
		sql = fmt.Sprintf("%sVALUES (...)", sql[:i])
	}

	if elapsedMs < 200 {
		l.log.Debug("database operation", "duration_ms", elapsedMs, "rows_affected", rows, "sql", sql)
	} else {
		l.log.Warn("database operation", "duration_ms", elapsedMs, "rows_affected", rows, "sql", sql)
	}
}
