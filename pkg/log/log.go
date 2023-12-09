package log

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

var _ Logger = (*JSONLogger)(nil)

type JSONLogger struct {
	logger *slog.Logger
}

func NewLogger() *JSONLogger {
	return &JSONLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
	}
}

func (l JSONLogger) With(args ...any) Logger {
	return &JSONLogger{
		logger: l.logger.With(args...),
	}
}

func (l JSONLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, l.withTrace(ctx, args)...)
}

func (l JSONLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, l.withTrace(ctx, args)...)
}

func (l JSONLogger) withTrace(ctx context.Context, args []any) []any {
	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.TraceID().IsValid() {
		args = append(args, slog.String("traceID", spanCtx.TraceID().String()))
	}

	if spanCtx.SpanID().IsValid() {
		args = append(args, slog.String("spanID", spanCtx.SpanID().String()))
	}

	return args
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

var _ Logger = (*NopLogger)(nil)

type NopLogger struct{}

func (l NopLogger) With(args ...any) Logger {
	return l
}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (l NopLogger) Info(ctx context.Context, msg string, args ...any) {}

func (l NopLogger) Error(ctx context.Context, msg string, args ...any) {}
