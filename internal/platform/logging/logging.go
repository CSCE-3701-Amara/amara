package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"
)

func New(output io.Writer, level slog.Level, format string) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	switch format {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = &consoleHandler{output: output, level: level, mu: &sync.Mutex{}}
	}

	return slog.New(handler)
}

// consoleHandler keeps structured JSON available for log collectors while
// presenting local container logs in the order people read them.
type consoleHandler struct {
	output io.Writer
	level  slog.Level
	mu     *sync.Mutex
	attrs  []slog.Attr
	groups []string
}

func (h *consoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *consoleHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, len(h.attrs)+record.NumAttrs())
	attrs = append(attrs, h.attrs...)
	record.Attrs(func(attr slog.Attr) bool {
		attrs = appendFlattened(attrs, strings.Join(h.groups, "."), attr)
		return true
	})
	values := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		values[attr.Key] = attr.Value.Resolve().String()
	}
	stamp := record.Time.UTC().Format("2006-01-02 15:04:05.000")
	if record.Time.IsZero() {
		stamp = time.Now().UTC().Format("2006-01-02 15:04:05.000")
	}
	service := strings.TrimPrefix(values["service"], "amara-")
	if service == "" {
		service = "app"
	}
	var line strings.Builder
	fmt.Fprintf(&line, "%s %-5s %-6s ", stamp, record.Level.String(), service)
	if record.Message == "HTTP request" && values["method"] != "" && values["route"] != "" {
		fmt.Fprintf(&line, "%s %s → %s  %sms", values["method"], values["route"], values["status"], values["duration_ms"])
		if values["request_id"] != "" {
			fmt.Fprintf(&line, "  id=%s", values["request_id"])
		}
	} else {
		line.WriteString(record.Message)
	}
	for _, attr := range attrs {
		if attr.Key == "service" || attr.Key == "environment" ||
			(record.Message == "HTTP request" && (attr.Key == "method" || attr.Key == "route" || attr.Key == "status" || attr.Key == "duration_ms" || attr.Key == "request_id")) {
			continue
		}
		value := attr.Value.Resolve().String()
		if strings.ContainsAny(value, " \n\r\t\"=") {
			value = strconv.Quote(value)
		}
		fmt.Fprintf(&line, "  %s=%s", attr.Key, value)
	}
	line.WriteByte('\n')
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.output, line.String())
	return err
}

func (h *consoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	copy := *h
	copy.attrs = append([]slog.Attr(nil), h.attrs...)
	for _, attr := range attrs {
		copy.attrs = appendFlattened(copy.attrs, strings.Join(h.groups, "."), attr)
	}
	return &copy
}

func (h *consoleHandler) WithGroup(name string) slog.Handler {
	copy := *h
	copy.groups = append(append([]string(nil), h.groups...), name)
	return &copy
}

func appendFlattened(attrs []slog.Attr, prefix string, attr slog.Attr) []slog.Attr {
	if attr.Equal(slog.Attr{}) {
		return attrs
	}
	key := attr.Key
	if prefix != "" {
		key = prefix + "." + key
	}
	value := attr.Value.Resolve()
	if value.Kind() == slog.KindGroup {
		for _, child := range value.Group() {
			attrs = appendFlattened(attrs, key, child)
		}
		return attrs
	}
	return append(attrs, slog.Attr{Key: key, Value: value})
}
