// Package logger for logging
package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Initialize logger without construction.
	// 	Because we want to puts logs from anywhere without initialization and without injection
	// We don't integrate rollbar to logger
	// 	The information to send rollbar and log is different

	logger = func() *zap.SugaredLogger {
		encodeCfg := zap.NewProductionEncoderConfig()
		encodeCfg.EncodeTime = zapcore.RFC3339TimeEncoder

		_ = godotenv.Load("./.env")
		timestamp := time.Now().Format("2006-01-02T15-04-05")
		logPath := fmt.Sprintf("%v/log_filename-%v-%v.log", os.Getenv("LOG_PATH"), os.Getenv("RUN_MODE"), timestamp)

		lumberjackLogger := &lumberjack.Logger{
			Filename:   logPath, // Base log file name
			MaxSize:    10,      // Max size in MB before rotation
			MaxBackups: 3,       // Max number of old log files to retain
			MaxAge:     28,      // Max number of days to retain old log files
			Compress:   true,    // Compress old log files (gzip)
		}

		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encodeCfg), zapcore.Lock(zapcore.AddSync(lumberjackLogger)), zap.NewAtomicLevel()),
			zapcore.NewCore(zapcore.NewJSONEncoder(encodeCfg), zapcore.Lock(os.Stdout), zap.NewAtomicLevel()),
		)
		l := zap.New(core, zap.AddCallerSkip(1))

		return l.Sugar()
	}()
)

// Logger wraps methods of zap logger for limited usage
type Logger interface {
	// Error prints error log
	Error(args ...any)
	// Errorf prints error log with the specified msg format
	Errorf(format string, args ...any)
	// Errorw prints msg with arguments
	Errorw(msg string, args ...any)

	// Warn prints warn log
	Warn(args ...any)
	// Warnf prints warn log with the specified msg format
	Warnf(format string, args ...any)
	// Warnw prints msg with arguments
	Warnw(msg string, args ...any)

	// Info prints info log
	Info(args ...any)
	// Infof prints info log with the specified msg format
	Infof(format string, args ...any)
	// Infow prints msg with arguments
	Infow(msg string, args ...any)
}

// GetLogger is used for integrating with other pkg to inject it inside of the pkg
func GetLogger() Logger {
	return logger
}

// Info logs a message at level Info on the standard logger.
func Info(args ...any) {
	logger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...any) {
	logger.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...any) {
	logger.Error(args...)
}

// Fatal logs a message at level Error on the standard logger.
func Fatal(args ...any) {
	logger.Fatal(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

// Infow prints msg with arguments
func Infow(msg string, args ...any) {
	logger.Infow(msg, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// Warnw prints msg with arguments
func Warnw(msg string, args ...any) {
	logger.Warnw(msg, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// Errorw prints msg with arguments
func Errorw(msg string, args ...any) {
	logger.Errorw(msg, args...)
}
