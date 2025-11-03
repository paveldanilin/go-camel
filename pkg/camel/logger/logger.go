package logger

import "context"

type LogLevel int

const (
	LogLevelError = iota + 1
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

type Logger interface {
	Log(ctx context.Context, level LogLevel, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
}
