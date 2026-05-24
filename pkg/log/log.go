package log

import (
	"context"
	"log/slog"
)

type logCtxKey struct{}

func New() *slog.Logger {
	logger := slog.New(slog.Default().Handler())
	return logger
}

func SetContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	l := ctx.Value(logCtxKey{})
	log, ok := l.(*slog.Logger)

	if !ok {
		panic("logger not in the context")
	}

	return log
}
