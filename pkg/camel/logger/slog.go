package logger

import (
	"context"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"log/slog"
)

// slogLogger represents a wrapper for std slog.
type slogLogger struct {
	logger *slog.Logger
	level  api.LogLevel
}

func NewSlog(logger *slog.Logger, level api.LogLevel) *slogLogger {
	return &slogLogger{
		logger: logger,
		level:  level,
	}
}

func (l *slogLogger) Log(ctx context.Context, level api.LogLevel, msg string, args ...any) {
	if l.level <= level {
		switch level {
		case api.LogLevelInfo:
			if len(args) == 0 {
				l.logger.InfoContext(ctx, msg)
			} else {
				l.logger.InfoContext(ctx, msg, args...)
			}
			return
		case api.LogLevelWarn:
			if len(args) == 0 {
				l.logger.WarnContext(ctx, msg)
			} else {
				l.logger.WarnContext(ctx, msg, args...)
			}
			return
		case api.LogLevelError:
			if len(args) == 0 {
				l.logger.ErrorContext(ctx, msg)
			} else {
				l.logger.ErrorContext(ctx, msg, args...)
			}
			return
		case api.LogLevelDebug:
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
	if l.level >= api.LogLevelInfo {
		if len(args) == 0 {
			l.logger.InfoContext(ctx, msg)
		} else {
			l.logger.InfoContext(ctx, msg, slog.Any("args", args))
		}
	}
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= api.LogLevelWarn {
		if len(args) == 0 {
			l.logger.WarnContext(ctx, msg)
		} else {
			l.logger.WarnContext(ctx, msg, slog.Any("args", args))
		}
	}
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= api.LogLevelError {
		if len(args) == 0 {
			l.logger.ErrorContext(ctx, msg)
		} else {
			l.logger.ErrorContext(ctx, msg, slog.Any("args", args))
		}
	}
}

func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	if l.level >= api.LogLevelDebug {
		if len(args) == 0 {
			l.logger.DebugContext(ctx, msg)
		} else {
			l.logger.DebugContext(ctx, msg, slog.Any("args", args))
		}
	}
}
