package logger

import (
	"context"
	"log/slog"
)

// AddFields adds provided fields to context logger
func AddFields(ctx context.Context, fields ...any) context.Context {
	return WithContext(ctx, FromContext(ctx).With(fields...))
}

func Address(value string) slog.Attr {
	return slog.String("address", value)
}
