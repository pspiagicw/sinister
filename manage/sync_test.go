package manage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pspiagicw/sinister/database"
)

// TestSync_EndToEnd wires all mocks together and verifies the full
// update → download → delete pipeline in a single run.
//
// Setup:
//   - RSS feed with 3 entries: 2 recent (1 day old), 1 old (20 days old)
//   - sync --days 7: downloads videos ≤7 days old, deletes downloaded videos >7 days old
//
// Expected outcome after sync:
//   - DB has 3 total entries
//   - 2 recent entries: watched=1, filepath set, files on disk
//   - 1 old entry: watched=0 (too old to download), no file
func TestSync_EndToEnd(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeSyncConfig(t, videoFolder)

	rssFeed := sampleRSSFeed("SyncChannel",
		rssEntry("Recent One", "https://www.youtube.com/watch?v=sync1", daysAgo(1))+
			rssEntry("Recent Two", "https://www.youtube.com/watch?v=sync2", daysAgo(2))+
			rssEntry("Old Video", "https://www.youtube.com/watch?v=sync3", daysAgo(20)),
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		fmt.Fprint(w, rssFeed)
	}))
	defer srv.Close()

	fetcher := func(_ string, _ int) ([]byte, error) {
		resp, err := http.Get(srv.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		buf := make([]byte, 0, 4096)
		tmp := make([]byte, 512)
		for {
			n, readErr := resp.Body.Read(tmp)
			buf = append(buf, tmp[:n]...)
			if readErr != nil {
				break
			}
		}
		return buf, nil
	}

	Sync(SyncOptions{
		ConfigPath: configPath,
		Update: UpdateOptions{
			URLs:    []string{srv.URL},
			Fetcher: fetcher,
		},
		Download: DownloadOptions{
			ConfigPath: configPath,
			Days:       7,
			Client:     &mockVideoClient{video: makeMuxedVideo()},
		},
		Delete: DeleteOptions{
			Days: 7,
		},
	})

	// All 3 entries should be in the DB
	total := database.TotalEntries()
	if total != 3 {
		t.Errorf("expected 3 total entries, got %d", total)
	}

	// 2 recent entries downloaded and marked watched
	downloaded := database.QueryDownloaded()
	if len(downloaded) != 2 {
		t.Errorf("expected 2 downloaded entries, got %d", len(downloaded))
	}
	for _, e := range downloaded {
		if _, err := os.Stat(e.FilePath); err != nil {
			t.Errorf("expected file on disk for %q: %v", e.Title, err)
		}
	}

	// 1 old entry still unwatched and not downloaded
	unwatched := database.QueryUnwatched()
	if len(unwatched) != 1 || unwatched[0].Title != "Old Video" {
		t.Errorf("expected 1 unwatched 'Old Video', got %v", unwatched)
	}
}

// TestSync_AutoDeleteAfterDownload verifies that when --days is set, videos
// downloaded on a previous sync run are removed from disk on the next run.
func TestSync_AutoDeleteAfterDownload(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeSyncConfig(t, videoFolder)

	// First sync: one video, 10 days old. days=30 so it downloads.
	oldFeed := sampleRSSFeed("Chan",
		rssEntry("Will Be Deleted", "https://www.youtube.com/watch?v=old99", daysAgo(10)),
	)
	srv1 := serveStaticRSS(t, oldFeed)
	defer srv1.Close()

	fetcher1 := staticFetcher(srv1)
	Sync(SyncOptions{
		ConfigPath: configPath,
		Update: UpdateOptions{URLs: []string{srv1.URL}, Fetcher: fetcher1},
		Download: DownloadOptions{
			ConfigPath: configPath,
			Days:       30,
			Client:     &mockVideoClient{video: makeMuxedVideo()},
		},
		Delete: DeleteOptions{Days: 30},
	})

	downloaded := database.QueryDownloaded()
	if len(downloaded) != 1 {
		t.Fatalf("expected 1 downloaded entry after first sync, got %d", len(downloaded))
	}
	oldFilePath := downloaded[0].FilePath

	// Second sync: days=5 — the 10-day-old video should now be deleted.
	srv2 := serveStaticRSS(t, oldFeed)
	defer srv2.Close()

	Sync(SyncOptions{
		ConfigPath: configPath,
		Update: UpdateOptions{URLs: []string{srv2.URL}, Fetcher: staticFetcher(srv2)},
		Download: DownloadOptions{
			ConfigPath: configPath,
			Days:       5,
			Client:     &mockVideoClient{video: makeMuxedVideo()},
		},
		Delete: DeleteOptions{Days: 5},
	})

	if _, err := os.Stat(oldFilePath); !os.IsNotExist(err) {
		t.Error("expected old video file to be deleted on second sync")
	}
	if len(database.QueryDownloaded()) != 0 {
		t.Error("expected no downloaded entries after second sync cleanup")
	}
}

// helpers

func writeSyncConfig(t *testing.T, videoFolder string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `videoFolder = "` + videoFolder + `"
quality = "hd720"
urls = ["https://example.com/feed"]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing sync config: %v", err)
	}
	return path
}

func serveStaticRSS(t *testing.T, xml string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/atom+xml")
		fmt.Fprint(w, xml)
	}))
}

func staticFetcher(srv *httptest.Server) func(string, int) ([]byte, error) {
	return func(_ string, _ int) ([]byte, error) {
		resp, err := http.Get(srv.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		buf := make([]byte, 0, 4096)
		tmp := make([]byte, 512)
		for {
			n, readErr := resp.Body.Read(tmp)
			buf = append(buf, tmp[:n]...)
			if readErr != nil {
				break
			}
		}
		return buf, nil
	}
}
