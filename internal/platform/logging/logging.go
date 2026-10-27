package logging

import (
	"io"
	"log/slog"
)

// New creates a structured JSON logger with the configured minimum level.
func New(output io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level}))
}
