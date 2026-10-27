package telemetry

import (
	"context"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"
	sentryslog "github.com/getsentry/sentry-go/slog"
)

var meter = sentry.NewMeter(context.Background())

func Init() {
	meter = sentry.NewMeter(context.Background())
}

func WithSentryLogs(log *slog.Logger) *slog.Logger {
	sentryHandler := sentryslog.Option{
		LogLevel:  []slog.Level{slog.LevelInfo, slog.LevelWarn, slog.LevelError, sentryslog.LevelFatal},
		AddSource: true,
	}.NewSentryHandler(context.Background())

	return slog.New(slog.NewMultiHandler(log.Handler(), sentryHandler))
}

func Count(ctx context.Context, name string, attrs ...attribute.Builder) {
	meter.WithCtx(ctx).Count(name, 1, sentry.WithAttributes(attrs...))
}

func Duration(ctx context.Context, name string, d time.Duration, attrs ...attribute.Builder) {
	meter.WithCtx(ctx).Distribution(
		name,
		float64(d.Microseconds())/1000,
		sentry.WithUnit(sentry.UnitMillisecond),
		sentry.WithAttributes(attrs...),
	)
}
