package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestNewJSONFormatWritesJSONLines(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, slog.LevelInfo, "json")
	logger.Info("HTTP request", "route", "/v1/jobs", "status", 200)
	for line := range strings.SplitSeq(strings.TrimSpace(output.String()), "\n") {
		if !json.Valid([]byte(line)) {
			t.Fatalf("JSON format must emit JSON lines, got: %s", line)
		}
	}
	if !strings.Contains(output.String(), `"route":"/v1/jobs"`) {
		t.Fatalf("attributes were not emitted: %s", output.String())
	}
}

func TestNewTextFormatWritesPlainLines(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, slog.LevelInfo, "text")
	logger.Info("HTTP request", "route", "/v1/jobs")
	text := strings.TrimSpace(output.String())
	if strings.HasPrefix(text, "{") {
		t.Fatalf("text format must not emit JSON, got: %s", text)
	}
	if !strings.Contains(text, "HTTP request") || !strings.Contains(text, "/v1/jobs") {
		t.Fatalf("message and attributes were not emitted: %s", text)
	}
}

func TestNewUnknownFormatFallsBackToText(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, slog.LevelInfo, "csv")
	logger.Info("HTTP request")
	text := strings.TrimSpace(output.String())
	if strings.HasPrefix(text, "{") {
		t.Fatalf("unknown formats must fall back to the text handler, got: %s", text)
	}
}

func TestNewRespectsLevel(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, slog.LevelWarn, "json")
	logger.Info("filtered")
	logger.Warn("emitted")
	if strings.Contains(output.String(), "filtered") {
		t.Fatal("records below the minimum level must be dropped")
	}
	if !strings.Contains(output.String(), "emitted") {
		t.Fatal("records at or above the minimum level must be emitted")
	}
}
