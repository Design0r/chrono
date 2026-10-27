package telemetry_test

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"

	"chrono/internal/telemetry"
)

type captured struct {
	mu      sync.Mutex
	logs    []sentry.Log
	metrics []sentry.Metric
}

func initSentry(t *testing.T) *captured {
	t.Helper()
	c := &captured{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: &sentry.MockTransport{},
		BeforeSendLog: func(l *sentry.Log) *sentry.Log {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.logs = append(c.logs, *l)
			return nil
		},
		BeforeSendMetric: func(m *sentry.Metric) *sentry.Metric {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.metrics = append(c.metrics, *m)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("failed to init sentry: %v", err)
	}
	telemetry.Init()
	t.Cleanup(func() { sentry.CurrentHub().BindClient(nil) })
	return c
}

func TestLogsAreSentToSentryFromInfoUpwards(t *testing.T) {
	c := initSentry(t)
	log := telemetry.WithSentryLogs(slog.New(slog.DiscardHandler))

	log.Debug("debug message")
	log.Info("request handled", "route", "/api/v1/users")
	log.Error("request failed")
	sentry.Flush(2 * time.Second)

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.logs) != 2 {
		t.Fatalf("expected 2 logs, got %d: %+v", len(c.logs), c.logs)
	}
	if c.logs[0].Body != "request handled" || c.logs[1].Body != "request failed" {
		t.Fatalf("unexpected log bodies: %q, %q", c.logs[0].Body, c.logs[1].Body)
	}
	if got := c.logs[0].Attributes["route"].AsString(); got != "/api/v1/users" {
		t.Fatalf("expected route attribute, got %q", got)
	}
}

func TestMetricsAreSentToSentry(t *testing.T) {
	c := initSentry(t)
	ctx := context.Background()

	telemetry.Count(ctx, "events.created", attribute.String("event.name", "urlaub"))
	telemetry.Duration(ctx, "http.server.request.duration", 1500*time.Microsecond)
	sentry.Flush(2 * time.Second)

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.metrics) != 2 {
		t.Fatalf("expected 2 metrics, got %d: %+v", len(c.metrics), c.metrics)
	}

	count := c.metrics[0]
	if v, _ := count.Value.Int64(); count.Name != "events.created" || v != 1 {
		t.Fatalf("unexpected counter %s=%v", count.Name, count.Value)
	}
	if got := count.Attributes["event.name"].AsString(); got != "urlaub" {
		t.Fatalf("expected event.name attribute, got %q", got)
	}

	duration := c.metrics[1]
	if v, _ := duration.Value.Float64(); v != 1.5 || duration.Unit != sentry.UnitMillisecond {
		t.Fatalf("unexpected distribution %v %s", duration.Value, duration.Unit)
	}
}
