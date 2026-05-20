package feed

import (
	"testing"
	"time"
)

func TestEntryDate_Valid(t *testing.T) {
	e := Entry{Published: "2024-03-15T10:30:00Z"}
	got := e.Date()

	want := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestEntryDate_WithOffset(t *testing.T) {
	e := Entry{Published: "2024-06-01T12:00:00+05:30"}
	got := e.Date()

	// Equivalent UTC time
	want := time.Date(2024, 6, 1, 6, 30, 0, 0, time.UTC)
	if !got.UTC().Equal(want) {
		t.Errorf("expected %v, got %v", want, got.UTC())
	}
}
