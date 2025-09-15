package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// NewTextMultiLogger creates a slog logger that writes identical text output
// to both stdout and a file. No external packages, no JSON.
func NewTextMultiLogger(path, level string, addSource bool) (*slog.Logger, *os.File, error) {
	// Ensure directory exists
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, nil, err
		}
	}

	// Open/append the log file
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}

	// One TextHandler that writes to both stdout and file -> identical bytes
	mw := io.MultiWriter(os.Stdout, f)

	h := slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: addSource,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// Match your Echo middleware time style: 2006/01/02 15:04:05
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006/01/02 15:04:05"))
			}
			return a
		},
	})

	return slog.New(h), f, nil
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
