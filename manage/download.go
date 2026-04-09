package manage

import (
	"fmt"
	"io"
	"mime"
	neturl "net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

type DownloadOptions struct {
	ConfigPath string
	Days       int
	Videos     int
}

func Download(opts DownloadOptions) {
	conf := config.ParseConfig(opts.ConfigPath)
	entries := database.QueryUnwatched()
	entries = filterDownloadEntries(entries, opts)

	if len(entries) == 0 {
		goreland.LogInfo("No unwatched videos matched the requested filters")
		return
	}

	if err := os.MkdirAll(conf.VideoFolder, 0755); err != nil {
		goreland.LogFatal("Error while creating video folder: %v", err)
	}

	client := youtube.Client{}
	successCount := 0
	failedCount := 0

	for _, entry := range entries {
		if err := downloadEntry(&client, conf.VideoFolder, entry); err != nil {
			failedCount++
			goreland.LogError("Skipping %s: %v", entry.Title, err)
			continue
		}

		entryCopy := entry
		database.UpdateWatched(&entryCopy)
		successCount++
		goreland.LogSuccess("Downloaded: %s", entry.Title)
	}

	goreland.LogSuccess("Download complete. success=%d failed=%d", successCount, failedCount)
}

func filterDownloadEntries(entries []feed.Entry, opts DownloadOptions) []feed.Entry {
	sort.SliceStable(entries, func(i, j int) bool {
		ti, okI := parsePublished(entries[i].Published)
		tj, okJ := parsePublished(entries[j].Published)

		if okI && okJ {
			return ti.After(tj)
		}
		if okI {
			return true
		}
		if okJ {
			return false
		}
		return false
	})

	filtered := entries
	if opts.Days > 0 {
		cutoff := time.Now().AddDate(0, 0, -opts.Days)
		tmp := make([]feed.Entry, 0, len(filtered))
		for _, entry := range filtered {
			t, ok := parsePublished(entry.Published)
			if ok && t.After(cutoff) {
				tmp = append(tmp, entry)
			}
		}
		filtered = tmp
	}

	if opts.Videos > 0 && len(filtered) > opts.Videos {
		filtered = filtered[:opts.Videos]
	}

	return filtered
}

func parsePublished(value string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func downloadEntry(client *youtube.Client, videoFolder string, entry feed.Entry) error {
	videoID, err := getVideoID(entry.Link.URL)
	if err != nil {
		return err
	}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return fmt.Errorf("error getting video metadata: %w", err)
	}

	format, err := getBestFormat(video)
	if err != nil {
		return err
	}

	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf("error getting stream: %w", err)
	}
	defer stream.Close()

	outputPath := getOutputPath(videoFolder, entry, format)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", outputPath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, stream); err != nil {
		return fmt.Errorf("error writing file %s: %w", outputPath, err)
	}

	return nil
}

func getBestFormat(video *youtube.Video) (*youtube.Format, error) {
	muxedVideoFormats := youtube.FormatList{}
	for _, format := range video.Formats {
		if isMuxedVideoFormat(format) {
			muxedVideoFormats = append(muxedVideoFormats, format)
		}
	}

	if len(muxedVideoFormats) == 0 {
		return nil, fmt.Errorf("no downloadable video+audio format found")
	}

	muxedVideoFormats.Sort()
	best := muxedVideoFormats[len(muxedVideoFormats)-1]
	return &best, nil
}

func isMuxedVideoFormat(format youtube.Format) bool {
	if format.AudioChannels == 0 {
		return false
	}

	parsed, _, err := mime.ParseMediaType(format.MimeType)
	if err != nil {
		return false
	}

	return strings.HasPrefix(parsed, "video/")
}

func getOutputPath(videoFolder string, entry feed.Entry, format *youtube.Format) string {
	ext := extensionFromMimeType(format.MimeType)
	if ext == "" {
		ext = "mp4"
	}
	return filepath.Join(videoFolder, entry.Slug+"."+ext)
}

func extensionFromMimeType(mimeType string) string {
	parsed, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return ""
	}

	if exts, err := mime.ExtensionsByType(parsed); err == nil && len(exts) > 0 {
		return strings.TrimPrefix(exts[0], ".")
	}

	parts := strings.Split(parsed, "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func getVideoID(url string) (string, error) {
	parsed, err := neturl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := strings.ToLower(parsed.Host)
	switch host {
	case "www.youtube.com", "youtube.com", "m.youtube.com":
		if id := strings.TrimSpace(parsed.Query().Get("v")); id != "" {
			return id, nil
		}

		pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(pathParts) == 2 && pathParts[0] == "shorts" {
			if id := strings.TrimSpace(pathParts[1]); id != "" {
				return id, nil
			}
		}
	case "youtu.be", "www.youtu.be":
		if id := strings.TrimSpace(strings.Trim(parsed.Path, "/")); id != "" {
			return id, nil
		}
	}

	return "", fmt.Errorf("unsupported YouTube URL: %s", url)
}
