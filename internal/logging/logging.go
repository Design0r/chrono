package logging

import (
	"log/slog"
	"os"
	"strings"
)

func NewTextLogger(level string, addSource bool) *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: addSource,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006/01/02 15:04:05"))
			}
			return a
		},
	})

	return slog.New(h)
}

func parseLevel(s string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
