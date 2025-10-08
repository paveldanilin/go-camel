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
			if len(args) == 0 {
				l.logger.InfoContext(ctx, msg)
			} else {
				l.logger.InfoContext(ctx, msg, args...)
			}
			return
		case LogLevelWarn:
			if len(args) == 0 {
				l.logger.WarnContext(ctx, msg)
			} else {
				l.logger.WarnContext(ctx, msg, args...)
			}
			return
		case LogLevelError:
			if len(args) == 0 {
				l.logger.ErrorContext(ctx, msg)
			} else {
				l.logger.ErrorContext(ctx, msg, args...)
			}
			return
		case LogLevelDebug:
			if len(args) == 0 {
				l.logger.DebugContext(ctx, msg)
			} else {
				l.logger.DebugContext(ctx, msg, args...)
			}
			return
		}
	}
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelInfo {
		if len(args) == 0 {
			l.logger.InfoContext(ctx, msg)
		} else {
			l.logger.InfoContext(ctx, msg, slog.Any("args", args))
		}
	}
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelWarn {
		if len(args) == 0 {
			l.logger.WarnContext(ctx, msg)
		} else {
			l.logger.WarnContext(ctx, msg, slog.Any("args", args))
		}
	}
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelError {
		if len(args) == 0 {
			l.logger.ErrorContext(ctx, msg)
		} else {
			l.logger.ErrorContext(ctx, msg, slog.Any("args", args))
		}
	}
}

func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	if l.level >= LogLevelDebug {
		if len(args) == 0 {
			l.logger.DebugContext(ctx, msg)
		} else {
			l.logger.DebugContext(ctx, msg, slog.Any("args", args))
		}
	}
}
