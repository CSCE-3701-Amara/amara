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
	logger := New(&output, slog.LevelInfo, "text").With("service", "amara-api", "environment", "development")
	logger.Info("HTTP request", "request_id", "ABC", "method", "POST", "route", "/v1/auth/sessions", "status", 429, "duration_ms", 0.066)
	text := strings.TrimSpace(output.String())
	if strings.HasPrefix(text, "{") {
		t.Fatalf("text format must not emit JSON, got: %s", text)
	}
	if !strings.Contains(text, "INFO  api    POST /v1/auth/sessions → 429  0.066ms  id=ABC") {
		t.Fatalf("request summary is unreadable: %s", text)
	}
	if strings.Contains(text, "environment=") || strings.Contains(text, "method=") {
		t.Fatalf("summary fields were duplicated: %s", text)
	}
}

func TestTextFormatRetainsErrorAndGroups(t *testing.T) {
	var output bytes.Buffer
	logger := New(&output, slog.LevelInfo, "text").With("service", "amara-worker")
	logger.WithGroup("grading").With("attempt", "42").Error("grading failed", "error", "provider unavailable")
	line := output.String()
	for _, expected := range []string{"ERROR worker grading failed", "grading.attempt=42", `grading.error="provider unavailable"`} {
		if !strings.Contains(line, expected) {
			t.Fatalf("missing %q in %s", expected, line)
		}
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
