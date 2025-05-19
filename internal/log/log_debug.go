//go:build !release

package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type PrettyHandler struct {
	slog.Handler
}

func levelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "\033[36m[DEBUG]\033[0m"
	case slog.LevelInfo:
		return "\033[32m[INFO]\033[0m"
	case slog.LevelWarn:
		return "\033[33m[WARN]\033[0m"
	case slog.LevelError:
		return "\033[31m[ERROR]\033[0m"
	default:
		return fmt.Sprintf("%v", level)
	}
}

func (h PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	var b strings.Builder
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(&b, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintf(os.Stdout, "%s %s%s\n", levelColor(r.Level), r.Message, b.String())
	return nil
}

func Init() {
	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(PrettyHandler{Handler: baseHandler})
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
