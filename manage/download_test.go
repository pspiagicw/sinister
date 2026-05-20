package manage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

// --- pure function unit tests ---

func TestGetVideoID(t *testing.T) {
	cases := []struct {
		url     string
		want    string
		wantErr bool
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://youtube.com/watch?v=abc123", "abc123", false},
		{"https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/shorts/abc123", "abc123", false},
		{"https://www.youtube.com/watch?v=", "", true},
		{"https://vimeo.com/123", "", true},
		{"not-a-url", "", true},
	}

	for _, tc := range cases {
		id, err := getVideoID(tc.url)
		if tc.wantErr {
			if err == nil {
				t.Errorf("getVideoID(%q): expected error, got id=%q", tc.url, id)
			}
		} else {
			if err != nil {
				t.Errorf("getVideoID(%q): unexpected error: %v", tc.url, err)
			}
			if id != tc.want {
				t.Errorf("getVideoID(%q): got %q, want %q", tc.url, id, tc.want)
			}
		}
	}
}

func TestIsShortURL(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"https://www.youtube.com/shorts/abc", true},
		{"https://youtube.com/shorts/xyz", true},
		{"https://www.youtube.com/watch?v=abc", false},
		{"https://youtu.be/abc", false},
		{"https://vimeo.com/shorts/abc", false},
	}

	for _, tc := range cases {
		got := isShortURL(tc.url)
		if got != tc.want {
			t.Errorf("isShortURL(%q) = %v, want %v", tc.url, got, tc.want)
		}
	}
}

func TestParseTargetHeight(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"hd720", 720},
		{"1080p", 1080},
		{"HD720", 720},
		{"", 0},
		{"unknown", 0},
		{"4k", 4},
	}

	for _, tc := range cases {
		got := parseTargetHeight(tc.input)
		if got != tc.want {
			t.Errorf("parseTargetHeight(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestOutputBaseName(t *testing.T) {
	e := feed.Entry{
		Author: feed.Author{Name: "Rick Astley"},
		Title:  "Never Gonna Give You Up",
		Slug:   "never-gonna-give-you-up",
	}
	got := outputBaseName(e)
	if got == "" {
		t.Error("expected non-empty base name")
	}
}

func TestOutputBaseName_FallbackToSlug(t *testing.T) {
	e := feed.Entry{Slug: "fallback-slug"}
	got := outputBaseName(e)
	if got != "fallback-slug" {
		t.Errorf("expected fallback-slug, got %q", got)
	}
}

func TestBetterFormat(t *testing.T) {
	f := func(height, bitrate, audio int) youtube.Format {
		return youtube.Format{
			Height:        height,
			Bitrate:       bitrate,
			AudioChannels: audio,
		}
	}

	if !betterFormat(f(1080, 0, 0), f(720, 0, 0)) {
		t.Error("higher height should win")
	}
	if !betterFormat(f(720, 2000, 0), f(720, 1000, 0)) {
		t.Error("higher bitrate should win at equal height")
	}
	if !betterFormat(f(720, 1000, 2), f(720, 1000, 0)) {
		t.Error("audio channels should be tiebreak")
	}
	if betterFormat(f(720, 1000, 0), f(1080, 0, 0)) {
		t.Error("lower height should not win")
	}
}

func TestFilterDownloadEntries_Days(t *testing.T) {
	entries := []feed.Entry{
		{Title: "New", Published: time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339)},
		{Title: "Old", Published: time.Now().AddDate(0, 0, -10).UTC().Format(time.RFC3339)},
	}

	filtered := filterDownloadEntries(entries, DownloadOptions{Days: 3})
	if len(filtered) != 1 || filtered[0].Title != "New" {
		t.Errorf("expected only 'New', got %v", filtered)
	}
}

func TestFilterDownloadEntries_VideoLimit(t *testing.T) {
	entries := []feed.Entry{
		{Title: "A", Published: time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339)},
		{Title: "B", Published: time.Now().AddDate(0, 0, -2).UTC().Format(time.RFC3339)},
		{Title: "C", Published: time.Now().AddDate(0, 0, -3).UTC().Format(time.RFC3339)},
	}

	filtered := filterDownloadEntries(entries, DownloadOptions{Videos: 2})
	if len(filtered) != 2 {
		t.Errorf("expected 2 entries, got %d", len(filtered))
	}
}

func TestFilterDownloadEntries_SortedNewestFirst(t *testing.T) {
	entries := []feed.Entry{
		{Title: "Old", Published: time.Now().AddDate(0, 0, -5).UTC().Format(time.RFC3339)},
		{Title: "New", Published: time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339)},
	}

	filtered := filterDownloadEntries(entries, DownloadOptions{})
	if len(filtered) == 0 {
		t.Fatal("expected entries after filter")
	}
	if filtered[0].Title != "New" {
		t.Errorf("expected newest first, got %q", filtered[0].Title)
	}
}

// --- integration tests with mock VideoClient ---

func writeDownloadConfig(t *testing.T, videoFolder string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `videoFolder = "` + videoFolder + `"
quality = "hd720"
urls = ["https://example.com/feed"]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing config: %v", err)
	}
	return path
}

func seedUnwatchedEntry(t *testing.T, title, link, published string) feed.Entry {
	t.Helper()
	e := &feed.Entry{
		Author:    feed.Author{Name: "TestChannel"},
		Title:     title,
		Published: published,
		Link:      feed.Link{URL: link},
	}
	database.Add(e)
	for _, entry := range database.QueryUnwatched() {
		if entry.Title == title {
			return entry
		}
	}
	t.Fatalf("seeded entry %q not found", title)
	return feed.Entry{}
}

func TestDownload_Success(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeDownloadConfig(t, videoFolder)

	seedUnwatchedEntry(t, "Test Video", "https://www.youtube.com/watch?v=abc123",
		time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339))

	Download(DownloadOptions{
		ConfigPath: configPath,
		Client:     &mockVideoClient{video: makeMuxedVideo()},
	})

	entries := database.QueryDownloaded()
	if len(entries) != 1 {
		t.Fatalf("expected 1 downloaded entry, got %d", len(entries))
	}
	if entries[0].FilePath == "" {
		t.Error("expected non-empty filepath in DB")
	}
	if _, err := os.Stat(entries[0].FilePath); err != nil {
		t.Errorf("expected file on disk: %v", err)
	}
	if entries[0].Watched != 1 {
		t.Error("expected entry marked watched")
	}
}

func TestDownload_SkipShortURL(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeDownloadConfig(t, videoFolder)

	seedUnwatchedEntry(t, "Short Video", "https://www.youtube.com/shorts/abc123",
		time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339))

	Download(DownloadOptions{
		ConfigPath: configPath,
		Client:     &mockVideoClient{video: makeMuxedVideo()},
	})

	if database.CountUnwatched() != 1 {
		t.Error("expected short URL entry to remain unwatched")
	}
}

func TestDownload_SkipTooShortDuration(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeDownloadConfig(t, videoFolder)

	seedUnwatchedEntry(t, "Tiny Video", "https://www.youtube.com/watch?v=tiny1",
		time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339))

	shortVideo := makeMuxedVideo()
	shortVideo.Duration = 30 * time.Second

	Download(DownloadOptions{
		ConfigPath: configPath,
		Client:     &mockVideoClient{video: shortVideo},
	})

	if database.CountUnwatched() != 1 {
		t.Error("expected short-duration video to remain unwatched")
	}
}

func TestDownload_SkipTooLongDuration(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeDownloadConfig(t, videoFolder)

	seedUnwatchedEntry(t, "Long Video", "https://www.youtube.com/watch?v=long1",
		time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339))

	longVideo := makeMuxedVideo()
	longVideo.Duration = 2 * time.Hour

	Download(DownloadOptions{
		ConfigPath: configPath,
		Client:     &mockVideoClient{video: longVideo},
	})

	if database.CountUnwatched() != 1 {
		t.Error("expected long-duration video to remain unwatched")
	}
}

func TestDownload_DaysFilter(t *testing.T) {
	setupTestDB(t)
	videoFolder := t.TempDir()
	configPath := writeDownloadConfig(t, videoFolder)

	seedUnwatchedEntry(t, "Recent", "https://www.youtube.com/watch?v=recent",
		time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339))
	seedUnwatchedEntry(t, "Old", "https://www.youtube.com/watch?v=old123",
		time.Now().AddDate(0, 0, -10).UTC().Format(time.RFC3339))

	Download(DownloadOptions{
		ConfigPath: configPath,
		Days:       3,
		Client:     &mockVideoClient{video: makeMuxedVideo()},
	})

	downloaded := database.QueryDownloaded()
	if len(downloaded) != 1 || downloaded[0].Title != "Recent" {
		t.Errorf("expected only 'Recent' to be downloaded, got %v", downloaded)
	}
}
