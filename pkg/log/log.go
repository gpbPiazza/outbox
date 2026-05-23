package log

import (
	"context"
	"log/slog"
)

type logCtxKey struct{}

func New(ctx context.Context) (*slog.Logger, context.Context) {
	logger := slog.New(slog.Default().Handler())
	ctx = context.WithValue(ctx, logCtxKey{}, logger)
	return logger, ctx
}

func FromContext(ctx context.Context) *slog.Logger {
	l := ctx.Value(logCtxKey{})
	log, ok := l.(*slog.Logger)
	if !ok {
		log, _ = New(ctx)
	}
	return log
}
