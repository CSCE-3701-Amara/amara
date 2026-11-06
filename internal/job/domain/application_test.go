package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestMissingAmaratRequiresEveryTitleAndMatchesCaseInsensitively(t *testing.T) {
	missing := MissingAmarat([]string{" Go ", "PostgreSQL", "Rust"}, []string{"go", " postgresql "})
	if want := []string{"Rust"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("MissingAmarat()=%v, want %v", missing, want)
	}
}
func TestMissingAmaratReportsAllWhenCandidateHasNone(t *testing.T) {
	missing := MissingAmarat([]string{"Go", "SQL"}, nil)
	if want := []string{"Go", "SQL"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("MissingAmarat()=%v, want %v", missing, want)
	}
}
func TestApplicationCapReached(t *testing.T) {
	limit := int32(3)
	for _, tc := range []struct {
		name     string
		max      *int32
		accepted int64
		want     bool
	}{
		{"unlimited", nil, 900, false}, {"below", &limit, 2, false}, {"at cap", &limit, 3, true}, {"above cap", &limit, 4, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ApplicationCapReached(tc.max, tc.accepted); got != tc.want {
				t.Fatalf("ApplicationCapReached()=%v, want %v", got, tc.want)
			}
		})
	}
}
func TestApplicationIsOpen(t *testing.T) {
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	before, after := now.Add(-time.Minute), now.Add(time.Minute)
	for _, tc := range []struct {
		name, status      string
		opening, deadline *time.Time
		want              bool
	}{
		{"no bounds", PostingPublished, nil, nil, true}, {"draft", PostingDraft, nil, nil, false},
		{"not opened", PostingPublished, &after, nil, false}, {"opens now", PostingPublished, &now, nil, true},
		{"deadline passed", PostingPublished, nil, &before, false}, {"deadline now", PostingPublished, nil, &now, false},
		{"deadline future", PostingPublished, nil, &after, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ApplicationIsOpen(tc.status, tc.opening, tc.deadline, now); got != tc.want {
				t.Fatalf("ApplicationIsOpen()=%v, want %v", got, tc.want)
			}
		})
	}
}
