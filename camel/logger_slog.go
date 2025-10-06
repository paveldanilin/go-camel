package camel

import (
	"context"
	"log/slog"
)

type slogLogger struct {
	logger *slog.Logger
	level  LogLevel
}

func NewSlogLogger(logger *slog.Logger, level LogLevel) *slogLogger {
	return &slogLogger{
		logger: logger,
		level:  level,
	}
}

func (l *slogLogger) Log(ctx context.Context, level LogLevel, msg string, args ...any) {
	if l.level <= level {
		switch level {
		case LogLevelInfo:
			l.logger.InfoContext(ctx, msg, slog.Any("args", args))
			return
		case LogLevelWarn:
			l.logger.WarnContext(ctx, msg, slog.Any("args", args))
			return
		case LogLevelError:
			l.logger.ErrorContext(ctx, msg, slog.Any("args", args))
			return
		case LogLevelDebug:
			l.logger.DebugContext(ctx, msg, slog.Any("args", args))
			return
		}
	}
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelInfo {
		l.logger.InfoContext(ctx, msg, slog.Any("args", args))
	}
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelWarn {
		l.logger.WarnContext(ctx, msg, slog.Any("args", args))
	}
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelError {
		l.logger.ErrorContext(ctx, msg, slog.Any("args", args))
	}
}

func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelDebug {
		l.logger.DebugContext(ctx, msg, slog.Any("args", args))
	}
}
