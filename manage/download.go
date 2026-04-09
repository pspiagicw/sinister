package manage

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

func Download(configPath string) {
	conf := config.ParseConfig(configPath)
	entries := database.QueryUnwatched()

	if len(entries) == 0 {
		goreland.LogInfo("No unwatched videos to download")
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
	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return nil, fmt.Errorf("no downloadable format with audio found")
	}

	formats.Sort()
	best := formats[len(formats)-1]
	return &best, nil
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
	id, found := strings.CutPrefix(url, "https://www.youtube.com/watch?v=")
	if !found || strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("invalid YouTube URL: %s", url)
	}
	return id, nil
}
