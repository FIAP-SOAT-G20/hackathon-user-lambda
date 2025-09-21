package logger

import (
	"context"
	"flag"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(env string) *Logger {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
	} else {
		handler = NewPrettyHandler(os.Stdout, PrettyHandlerOptions{SlogOpts: slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}})
	}
	if flag.Lookup("test.v") != nil {
		handler = slog.DiscardHandler
	}
	return &Logger{logger: slog.New(handler)}
}

func (l *Logger) Debug(msg string, args ...any) { l.logger.Debug(msg, args...) }
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}
func (l *Logger) Info(msg string, args ...any) { l.logger.Info(msg, args...) }
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}
func (l *Logger) Warn(msg string, args ...any) { l.logger.Warn(msg, args...) }
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}
func (l *Logger) Error(msg string, args ...any) { l.logger.Error(msg, args...) }
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
