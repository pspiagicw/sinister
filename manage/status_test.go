package manage

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

func writeStatusConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `videoFolder = "/tmp/videos"
quality = "hd720"
urls = ["https://example.com/feed"]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	return path
}

func TestStatus_Counts(t *testing.T) {
	setupTestDB(t)
	addTestEntries(t)

	// Mark one entry watched
	entries := database.QueryUnwatched()
	database.UpdateWatched(&entries[0])

	summary := buildStatusSummary(&config.Config{
		VideoFolder: "/tmp",
		URLS:        []string{"https://x.com"},
	}, "")

	if summary.TotalVideos != 3 {
		t.Errorf("expected 3 total, got %d", summary.TotalVideos)
	}
	if summary.Unwatched != 2 {
		t.Errorf("expected 2 unwatched, got %d", summary.Unwatched)
	}
	if summary.Watched != 1 {
		t.Errorf("expected 1 watched, got %d", summary.Watched)
	}
}

func TestStatus_ByCreator(t *testing.T) {
	setupTestDB(t)
	addTestEntries(t)

	summary := buildStatusSummary(&config.Config{
		VideoFolder: "/tmp",
		URLS:        []string{"https://x.com"},
	}, "Alice")

	if summary.TotalVideos != 2 {
		t.Errorf("expected 2 Alice videos, got %d", summary.TotalVideos)
	}
	if summary.TotalCreators != 0 {
		t.Error("per-creator summary should not include TotalCreators")
	}
}

func TestStatus_JSON(t *testing.T) {
	setupTestDB(t)

	e := &feed.Entry{
		Author:    feed.Author{Name: "JSONChan"},
		Title:     "JSON Video",
		Published: daysAgo(1),
		Link:      feed.Link{URL: "https://youtube.com/watch?v=json1"},
	}
	database.Add(e)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Status(StatusOptions{
		ConfigPath: writeStatusConfig(t),
		JSON:       true,
	})

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	var out StatusSummary
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}
	if out.TotalVideos != 1 {
		t.Errorf("expected 1 video in JSON, got %d", out.TotalVideos)
	}
}
