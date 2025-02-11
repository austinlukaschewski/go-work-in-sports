package xlog

import (
	"log/slog"
	"os"
	"sync"
)

var Level = new(slog.LevelVar)

// Singleton logger using slog
var logger = func() *slog.Logger {
	defaultClient.Once.Do(func() {
		defaultClient.Value = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: Level}))
	})

	return defaultClient.Value
}

var defaultClient struct {
	Value *slog.Logger
	Once  sync.Once
}

func Debug(msg string, args ...interface{}) {
	logger().Debug(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger().Error(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger().Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger().Warn(msg, args...)
}
