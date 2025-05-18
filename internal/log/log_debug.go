//go:build !release

package log

import (
	"log/slog"
	"os"
)

func Init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}
