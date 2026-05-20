package manage

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

// setupTestDB redirects all DB operations to a fresh isolated SQLite file
// for this test, restoring the default path when the test ends.
func setupTestDB(t *testing.T) {
	t.Helper()
	database.SetDBPath(filepath.Join(t.TempDir(), "test.db"))
	t.Cleanup(func() { database.SetDBPath("") })
}

// mockVideoClient satisfies VideoClient without touching the network.
type mockVideoClient struct {
	video    *youtube.Video
	videoErr error
	stream   []byte
}

func (m *mockVideoClient) GetVideo(_ string) (*youtube.Video, error) {
	return m.video, m.videoErr
}

func (m *mockVideoClient) GetStream(_ *youtube.Video, _ *youtube.Format) (io.ReadCloser, int64, error) {
	data := m.stream
	if data == nil {
		data = []byte("fake-video-bytes")
	}
	return io.NopCloser(bytes.NewReader(data)), int64(len(data)), nil
}

// makeMuxedVideo returns a Video with a 5-minute duration and a single
// muxed format (AudioChannels > 0) so the ffmpeg path is never taken.
func makeMuxedVideo() *youtube.Video {
	return &youtube.Video{
		Duration: 5 * time.Minute,
		Formats: youtube.FormatList{
			{
				AudioChannels: 2,
				Quality:       "medium",
				QualityLabel:  "360p",
				MimeType:      "video/mp4; codecs=\"avc1.42001E, mp4a.40.2\"",
				Bitrate:       500000,
				Height:        360,
			},
		},
	}
}

// sampleRSSFeed returns a minimal Atom feed XML string.
func sampleRSSFeed(author, entries string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <author><name>` + author + `</name></author>
  ` + entries + `
</feed>`
}

// rssEntry returns a single <entry> element for use in sampleRSSFeed.
func rssEntry(title, link, published string) string {
	return `<entry>
    <title>` + title + `</title>
    <link href="` + link + `"/>
    <published>` + published + `</published>
    <author><name></name></author>
  </entry>`
}

// daysAgo returns an RFC3339 timestamp N days in the past.
func daysAgo(n int) string {
	return time.Now().AddDate(0, 0, -n).UTC().Format(time.RFC3339)
}

// makeStringEntries creates n dummy feed.Entry values for slice tests.
func makeStringEntries(n int) []feed.Entry {
	entries := make([]feed.Entry, n)
	for i := range entries {
		entries[i] = feed.Entry{Title: fmt.Sprintf("Video %d", i)}
	}
	return entries
}
