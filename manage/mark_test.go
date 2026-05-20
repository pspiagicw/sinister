package manage

import (
	"testing"
	"time"

	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

// --- pure function unit tests ---

func TestMakeSet(t *testing.T) {
	s := makeSet([]string{"a", "b", "a", "  c  "})
	if !s["a"] || !s["b"] || !s["c"] {
		t.Errorf("makeSet missing entries: %v", s)
	}
	if len(s) != 3 {
		t.Errorf("makeSet expected 3 unique entries, got %d", len(s))
	}
}

func TestHasSelectionFlags(t *testing.T) {
	if hasSelectionFlags(MarkOptions{}) {
		t.Error("empty opts should return false")
	}
	if !hasSelectionFlags(MarkOptions{Creator: "Alice"}) {
		t.Error("Creator should be a selection flag")
	}
	if !hasSelectionFlags(MarkOptions{Slugs: []string{"s1"}}) {
		t.Error("Slugs should be a selection flag")
	}
	if !hasSelectionFlags(MarkOptions{AllUnwatched: true}) {
		t.Error("AllUnwatched should be a selection flag")
	}
}

func TestSelectMarkTargets_BySlug(t *testing.T) {
	entries := testEntries()
	targets := selectMarkTargets(entries, MarkOptions{Slugs: []string{"slug-a"}})
	if len(targets) != 1 || targets[0].Slug != "slug-a" {
		t.Errorf("expected slug-a, got %v", targets)
	}
}

func TestSelectMarkTargets_ByCreator(t *testing.T) {
	entries := testEntries()
	targets := selectMarkTargets(entries, MarkOptions{Creator: "Alice"})
	if len(targets) != 2 {
		t.Errorf("expected 2 Alice entries, got %d", len(targets))
	}
}

func TestSelectMarkTargets_AllUnwatched(t *testing.T) {
	entries := testEntries()
	targets := selectMarkTargets(entries, MarkOptions{AllUnwatched: true})
	if len(targets) != len(entries) {
		t.Errorf("expected all %d entries, got %d", len(entries), len(targets))
	}
}

func TestSelectMarkTargets_NoFlags_ReturnsAll(t *testing.T) {
	entries := testEntries()
	// When hasSelectionFlags is false, selectMarkTargets is not called from Mark().
	// But if called directly with no filters, it returns all entries.
	targets := selectMarkTargets(entries, MarkOptions{})
	if len(targets) != len(entries) {
		t.Errorf("expected all entries with no filter, got %d", len(targets))
	}
}

// --- integration tests ---

func TestMark_BySlug(t *testing.T) {
	setupTestDB(t)

	addTestEntries(t)

	Mark(MarkOptions{Slugs: []string{"video-one"}})

	entries := database.QueryUnwatched()
	for _, e := range entries {
		if e.Slug == "video-one" {
			t.Error("video-one should be marked watched")
		}
	}
}

func TestMark_AllUnwatched(t *testing.T) {
	setupTestDB(t)
	addTestEntries(t)

	Mark(MarkOptions{AllUnwatched: true})

	if database.CountUnwatched() != 0 {
		t.Errorf("expected 0 unwatched, got %d", database.CountUnwatched())
	}
}

func TestMark_MarkAllUnwatched_Reset(t *testing.T) {
	setupTestDB(t)
	addTestEntries(t)
	Mark(MarkOptions{AllUnwatched: true})

	Mark(MarkOptions{MarkAllUnwatched: true})

	if database.CountUnwatched() != database.TotalEntries() {
		t.Error("MarkAllUnwatched should reset everything to unwatched")
	}
}

func TestMark_DryRun(t *testing.T) {
	setupTestDB(t)
	addTestEntries(t)
	before := database.CountUnwatched()

	Mark(MarkOptions{AllUnwatched: true, DryRun: true})

	if database.CountUnwatched() != before {
		t.Error("dry run should not change watched status")
	}
}

// helpers

func testEntries() []feed.Entry {
	published := time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339)
	return []feed.Entry{
		{Author: feed.Author{Name: "Alice"}, Title: "Video One", Slug: "slug-a", Published: published, Link: feed.Link{URL: "https://youtube.com/watch?v=a1"}},
		{Author: feed.Author{Name: "Alice"}, Title: "Video Two", Slug: "slug-b", Published: published, Link: feed.Link{URL: "https://youtube.com/watch?v=a2"}},
		{Author: feed.Author{Name: "Bob"}, Title: "Video Three", Slug: "slug-c", Published: published, Link: feed.Link{URL: "https://youtube.com/watch?v=b1"}},
	}
}

func addTestEntries(t *testing.T) {
	t.Helper()
	for _, e := range testEntries() {
		entryCopy := e
		database.Add(&entryCopy)
	}
}
