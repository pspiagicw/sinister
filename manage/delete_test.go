package manage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

// --- pure function unit tests ---

func TestHasDeleteSelectionFlags(t *testing.T) {
	if hasDeleteSelectionFlags(DeleteOptions{}) {
		t.Error("empty opts should return false")
	}
	if !hasDeleteSelectionFlags(DeleteOptions{Creator: "Alice"}) {
		t.Error("Creator should be a selection flag")
	}
	if !hasDeleteSelectionFlags(DeleteOptions{Slugs: []string{"s1"}}) {
		t.Error("Slugs should be a selection flag")
	}
	if !hasDeleteSelectionFlags(DeleteOptions{Days: 7}) {
		t.Error("Days should be a selection flag")
	}
}

func TestFilterDeleteTargets_ByCreator(t *testing.T) {
	entries := deleteTestEntries()
	targets := filterDeleteTargets(entries, DeleteOptions{Creator: "Alice"})
	if len(targets) != 1 || targets[0].Author.Name != "Alice" {
		t.Errorf("expected 1 Alice entry, got %v", targets)
	}
}

func TestFilterDeleteTargets_BySlug(t *testing.T) {
	entries := deleteTestEntries()
	targets := filterDeleteTargets(entries, DeleteOptions{Slugs: []string{"slug-del-a"}})
	if len(targets) != 1 || targets[0].Slug != "slug-del-a" {
		t.Errorf("expected slug-del-a only, got %v", targets)
	}
}

func TestFilterDeleteTargets_ByDays(t *testing.T) {
	entries := deleteTestEntries()
	// entries[0] is 20 days old, entries[1] is 1 day old
	targets := filterDeleteTargets(entries, DeleteOptions{Days: 7})
	if len(targets) != 1 || targets[0].Title != "Old Video" {
		t.Errorf("expected only old video (>7 days), got %v", targets)
	}
}

// --- integration tests ---

func seedDownloadedEntry(t *testing.T, title, link, published string) feed.Entry {
	t.Helper()
	e := &feed.Entry{
		Author:    feed.Author{Name: "TestChan"},
		Title:     title,
		Published: published,
		Link:      feed.Link{URL: link},
	}
	database.Add(e)
	for _, entry := range database.QueryUnwatched() {
		if entry.Title == title {
			database.UpdateWatched(&entry)
			tmpFile := filepath.Join(t.TempDir(), entry.Slug+".mp4")
			if err := os.WriteFile(tmpFile, []byte("video"), 0644); err != nil {
				t.Fatalf("creating fake video file: %v", err)
			}
			database.UpdateFilePath(entry.Slug, tmpFile)
			// Return the updated entry
			for _, dl := range database.QueryDownloaded() {
				if dl.Title == title {
					return dl
				}
			}
		}
	}
	t.Fatalf("could not seed downloaded entry %q", title)
	return feed.Entry{}
}

func TestDelete_RemovesFile(t *testing.T) {
	setupTestDB(t)

	entry := seedDownloadedEntry(t, "Delete Me", "https://youtube.com/watch?v=del1",
		daysAgo(1))

	Delete(DeleteOptions{Slugs: []string{entry.Slug}})

	if _, err := os.Stat(entry.FilePath); !os.IsNotExist(err) {
		t.Error("expected file to be deleted from disk")
	}
	if len(database.QueryDownloaded()) != 0 {
		t.Error("expected filepath cleared in DB")
	}
}

func TestDelete_FileAlreadyGone(t *testing.T) {
	setupTestDB(t)

	entry := seedDownloadedEntry(t, "Already Gone", "https://youtube.com/watch?v=gone1",
		daysAgo(1))

	// Remove the file before calling Delete
	os.Remove(entry.FilePath)

	// Should not panic or error — tolerate missing file
	Delete(DeleteOptions{Slugs: []string{entry.Slug}})

	if len(database.QueryDownloaded()) != 0 {
		t.Error("expected filepath cleared even when file was already missing")
	}
}

func TestDelete_MarkUnwatched(t *testing.T) {
	setupTestDB(t)

	entry := seedDownloadedEntry(t, "Rewatch Me", "https://youtube.com/watch?v=rw1",
		daysAgo(1))

	Delete(DeleteOptions{
		Slugs:         []string{entry.Slug},
		MarkUnwatched: true,
	})

	if database.CountUnwatched() != 1 {
		t.Errorf("expected entry marked unwatched after delete, got %d unwatched", database.CountUnwatched())
	}
}

func TestDelete_DryRun(t *testing.T) {
	setupTestDB(t)

	entry := seedDownloadedEntry(t, "Keep Me", "https://youtube.com/watch?v=keep1",
		daysAgo(1))

	Delete(DeleteOptions{
		Slugs:  []string{entry.Slug},
		DryRun: true,
	})

	if _, err := os.Stat(entry.FilePath); err != nil {
		t.Error("dry run should not delete the file")
	}
	if len(database.QueryDownloaded()) != 1 {
		t.Error("dry run should not clear filepath in DB")
	}
}

func TestDelete_DaysFilter(t *testing.T) {
	setupTestDB(t)

	old := seedDownloadedEntry(t, "Old Download", "https://youtube.com/watch?v=old99",
		daysAgo(20))
	recent := seedDownloadedEntry(t, "Recent Download", "https://youtube.com/watch?v=rec99",
		daysAgo(1))

	Delete(DeleteOptions{Days: 7})

	// Old should be gone, recent should remain
	if _, err := os.Stat(old.FilePath); !os.IsNotExist(err) {
		t.Error("old download should have been deleted")
	}
	remaining := database.QueryDownloaded()
	if len(remaining) != 1 || remaining[0].Title != "Recent Download" {
		t.Errorf("expected recent download to remain, got %v", remaining)
	}
	_ = recent
}

// helpers

func deleteTestEntries() []feed.Entry {
	return []feed.Entry{
		{
			Author:    feed.Author{Name: "Alice"},
			Title:     "Old Video",
			Slug:      "slug-del-a",
			Published: time.Now().AddDate(0, 0, -20).UTC().Format(time.RFC3339),
			FilePath:  "/tmp/old.mp4",
		},
		{
			Author:    feed.Author{Name: "Bob"},
			Title:     "New Video",
			Slug:      "slug-del-b",
			Published: time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339),
			FilePath:  "/tmp/new.mp4",
		},
	}
}
