package http

import (
	"encoding/json"
	"testing"
	"time"
)

// Posting windows are stored in TIMESTAMP columns, so an offset must be
// converted rather than dropped: 17:00+02:00 is 15:00 UTC, not 17:00 UTC.
func TestPostingTimeNormalizesOffsetToUTC(t *testing.T) {
	var value timeJSON
	if err := json.Unmarshal([]byte(`"2026-12-01T17:00:00+02:00"`), &value); err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 12, 1, 15, 0, 0, 0, time.UTC)
	if !value.Time.Equal(want) || value.Location() != time.UTC {
		t.Fatalf("parsed %v, want %v in UTC", value.Time, want)
	}
}
