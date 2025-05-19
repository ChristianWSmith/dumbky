//go:build release

package log

func Init() {}

func Debug(msg string, args ...any) {}

func Info(msg string, args ...any) {}

func Warn(msg string, args ...any) {}

func Error(msg string, args ...any) {}
