package logger

import (
	"context"
	"regexp"
	"time"

	gormLogger "gorm.io/gorm/logger"
	gormUtil "gorm.io/gorm/utils"
)

const (
	fileKey = "file"
	dataKey = "data"
	rowsKey = "rows"
)

// Config config for gormlogger
type Config struct {
	LogSQL      bool
	NotifyQuery string
}

// NewGormLogger build logger that satisfies gorm logger interface
// ref: https://github.com/go-gorm/gorm/blob/1b9cd56c5336ba6e22936c289e586261b75d7b35/logger/logger.go
func NewGormLogger(cfg *Config) gormLogger.Interface {
	// On prod, SQL may contain personal information, and LogSQL should be false
	l := &customLogger{logLevel: gormLogger.Info, logSQL: cfg.LogSQL}
	if cfg.NotifyQuery != "" {
		l.notifyQuery = regexp.MustCompile(cfg.NotifyQuery)
	}
	return l
}

type customLogger struct {
	logLevel    gormLogger.LogLevel
	logSQL      bool
	notifyQuery *regexp.Regexp
}

// LogMode log mode
func (l *customLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.logLevel = level
	return l
}

// Info print info
func (l *customLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Info {
		Infow("[Info] "+msg, fileKey, gormUtil.FileWithLineNum(), dataKey, data)
	}
}

// Warn print warn messages
func (l *customLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Warn {
		Warnw("[Warn]"+msg, fileKey, gormUtil.FileWithLineNum(), dataKey, data)
	}
}

// Error print error messages
func (l *customLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Error {
		Errorw("[Error]"+msg, fileKey, gormUtil.FileWithLineNum(), dataKey, data)
	}
}

// Trace print sql related messages
func (l *customLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), _ error) {
	var (
		elapsed   = time.Since(begin)
		sql, rows = fc()
	)

	if l.notifyQuery != nil && l.notifyQuery.MatchString(sql) {
		// TODO: Add rollbar or slack notification
		logger.Warnf("detected dangerous query: %s", sql)
	}

	if l.logSQL {
		Infow(gormUtil.FileWithLineNum(), "duration(ms)", float64(elapsed.Nanoseconds())/1e6, rowsKey, rows, "sql", sql)
	} else {
		Infow(gormUtil.FileWithLineNum(), "duration(ms)", float64(elapsed.Nanoseconds())/1e6, rowsKey, rows)
	}
}
