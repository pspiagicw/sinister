package database

import (
	"path/filepath"
	"testing"

	"github.com/pspiagicw/sinister/feed"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	SetDBPath(filepath.Join(t.TempDir(), "test.db"))
	t.Cleanup(func() { SetDBPath("") })
}

func makeEntry(author, title, link, published string) *feed.Entry {
	return &feed.Entry{
		Author:    feed.Author{Name: author},
		Title:     title,
		Published: published,
		Link:      feed.Link{URL: link},
	}
}

func TestInsertAndQueryUnwatched(t *testing.T) {
	setupTestDB(t)

	e := makeEntry("Author", "Title One", "https://youtube.com/watch?v=aaa", "2024-01-01T00:00:00Z")
	inserted := Add(e)
	if !inserted {
		t.Fatal("expected entry to be inserted")
	}

	entries := QueryUnwatched()
	if len(entries) != 1 {
		t.Fatalf("expected 1 unwatched entry, got %d", len(entries))
	}
	if entries[0].Title != "Title One" {
		t.Errorf("unexpected title: %s", entries[0].Title)
	}
}

func TestInsertDedup(t *testing.T) {
	setupTestDB(t)

	e := makeEntry("Author", "Same Title", "https://youtube.com/watch?v=bbb", "2024-01-01T00:00:00Z")
	Add(e)
	second := Add(e)
	if second {
		t.Error("duplicate insert should return false")
	}

	if TotalEntries() != 1 {
		t.Errorf("expected 1 entry, got %d", TotalEntries())
	}
}

func TestUpdateWatched(t *testing.T) {
	setupTestDB(t)

	e := makeEntry("Author", "Watch Me", "https://youtube.com/watch?v=ccc", "2024-01-01T00:00:00Z")
	Add(e)

	entries := QueryUnwatched()
	UpdateWatched(&entries[0])

	if CountUnwatched() != 0 {
		t.Errorf("expected 0 unwatched, got %d", CountUnwatched())
	}
	if TotalEntries() != 1 {
		t.Errorf("expected 1 total, got %d", TotalEntries())
	}
}

func TestUpdateAndClearFilePath(t *testing.T) {
	setupTestDB(t)

	e := makeEntry("Author", "File Video", "https://youtube.com/watch?v=ddd", "2024-01-01T00:00:00Z")
	Add(e)

	entries := QueryUnwatched()
	UpdateWatched(&entries[0])
	UpdateFilePath(entries[0].Slug, "/tmp/video.mp4")

	downloaded := QueryDownloaded()
	if len(downloaded) != 1 {
		t.Fatalf("expected 1 downloaded, got %d", len(downloaded))
	}
	if downloaded[0].FilePath != "/tmp/video.mp4" {
		t.Errorf("unexpected filepath: %s", downloaded[0].FilePath)
	}

	ClearFilePath(entries[0].Slug)
	if len(QueryDownloaded()) != 0 {
		t.Error("expected 0 downloaded after clearing filepath")
	}
}

func TestQueryDownloadedByCreator(t *testing.T) {
	setupTestDB(t)

	e1 := makeEntry("Alice", "Video A", "https://youtube.com/watch?v=e1", "2024-01-01T00:00:00Z")
	e2 := makeEntry("Bob", "Video B", "https://youtube.com/watch?v=e2", "2024-01-01T00:00:00Z")
	Add(e1)
	Add(e2)

	all := QueryUnwatched()
	for _, e := range all {
		entryCopy := e
		UpdateWatched(&entryCopy)
		UpdateFilePath(e.Slug, "/tmp/"+e.Slug+".mp4")
	}

	alice := QueryDownloadedByCreator("Alice")
	if len(alice) != 1 || alice[0].Author.Name != "Alice" {
		t.Errorf("expected 1 Alice entry, got %v", alice)
	}

	bob := QueryDownloadedByCreator("Bob")
	if len(bob) != 1 || bob[0].Author.Name != "Bob" {
		t.Errorf("expected 1 Bob entry, got %v", bob)
	}
}

func TestMarkAllUnwatched(t *testing.T) {
	setupTestDB(t)

	for i, title := range []string{"V1", "V2", "V3"} {
		e := makeEntry("Author", title, "https://youtube.com/watch?v="+string(rune('a'+i)), "2024-01-01T00:00:00Z")
		Add(e)
	}

	for _, e := range QueryUnwatched() {
		entryCopy := e
		UpdateWatched(&entryCopy)
	}
	if CountUnwatched() != 0 {
		t.Fatal("expected all watched before reset")
	}

	n := MarkAllUnwatched()
	if n != 3 {
		t.Errorf("expected 3 reset, got %d", n)
	}
	if CountUnwatched() != 3 {
		t.Errorf("expected 3 unwatched after reset, got %d", CountUnwatched())
	}
}

func TestMigrateIdempotent(t *testing.T) {
	setupTestDB(t)
	// Opening the DB twice should not panic or fail on duplicate column.
	db := openDB()
	closeDB(db)
	db2 := openDB()
	closeDB(db2)
}

func TestQueryCreatorStats(t *testing.T) {
	setupTestDB(t)

	Add(makeEntry("Alice", "A1", "https://youtube.com/watch?v=f1", "2024-01-01T00:00:00Z"))
	Add(makeEntry("Alice", "A2", "https://youtube.com/watch?v=f2", "2024-01-01T00:00:00Z"))
	Add(makeEntry("Bob", "B1", "https://youtube.com/watch?v=f3", "2024-01-01T00:00:00Z"))

	// Mark one Alice entry watched
	entries := QueryUnwatched()
	for _, e := range entries {
		if e.Author.Name == "Alice" && e.Title == "A1" {
			entryCopy := e
			UpdateWatched(&entryCopy)
			break
		}
	}

	stats := QueryCreatorStats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 creator stats, got %d", len(stats))
	}

	byName := map[string]CreatorStat{}
	for _, s := range stats {
		byName[s.Name] = s
	}

	if byName["Alice"].Total != 2 || byName["Alice"].Unwatched != 1 {
		t.Errorf("unexpected Alice stats: %+v", byName["Alice"])
	}
	if byName["Bob"].Total != 1 || byName["Bob"].Unwatched != 1 {
		t.Errorf("unexpected Bob stats: %+v", byName["Bob"])
	}
}

func TestUpdateUnwatched(t *testing.T) {
	setupTestDB(t)

	e := makeEntry("Author", "Reset Me", "https://youtube.com/watch?v=g1", "2024-01-01T00:00:00Z")
	Add(e)

	entries := QueryUnwatched()
	UpdateWatched(&entries[0])
	if CountUnwatched() != 0 {
		t.Fatal("expected 0 unwatched")
	}

	UpdateUnwatched(&entries[0])
	if CountUnwatched() != 1 {
		t.Errorf("expected 1 unwatched after reset, got %d", CountUnwatched())
	}
}
