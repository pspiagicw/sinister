package manage

import (
	"fmt"
	"io"
	"mime"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

const maxDownloadDuration = time.Hour

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
		if err := downloadEntry(&client, conf.VideoFolder, conf.Quality, entry); err != nil {
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

func downloadEntry(client *youtube.Client, videoFolder, quality string, entry feed.Entry) error {
	if isShortURL(entry.Link.URL) {
		return fmt.Errorf("shorts are skipped")
	}

	videoID, err := getVideoID(entry.Link.URL)
	if err != nil {
		return err
	}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return fmt.Errorf("error getting video metadata: %w", err)
	}
	if video.Duration > maxDownloadDuration {
		return fmt.Errorf("video duration %s exceeds 1h limit", video.Duration.Round(time.Second))
	}

	videoFormat, err := getBestVideoFormat(video, quality)
	if err != nil {
		return err
	}

	if videoFormat.AudioChannels > 0 {
		outputPath := getOutputPath(videoFolder, entry, videoFormat)
		return downloadStreamToFile(client, video, videoFormat, outputPath)
	}

	audioFormat, err := getBestAudioFormat(video)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "sinister-download-*")
	if err != nil {
		return fmt.Errorf("error creating temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	videoTempPath := filepath.Join(tempDir, "video."+extensionFromMimeType(videoFormat.MimeType))
	audioTempPath := filepath.Join(tempDir, "audio."+extensionFromMimeType(audioFormat.MimeType))
	outputPath := filepath.Join(videoFolder, entry.Slug+".mkv")

	if err := downloadStreamToFile(client, video, videoFormat, videoTempPath); err != nil {
		return err
	}
	if err := downloadStreamToFile(client, video, audioFormat, audioTempPath); err != nil {
		return err
	}
	if err := mergeVideoAndAudio(videoTempPath, audioTempPath, outputPath); err != nil {
		return err
	}

	return nil
}

func getBestVideoFormat(video *youtube.Video, quality string) (*youtube.Format, error) {
	videoFormats := youtube.FormatList{}
	for _, format := range video.Formats {
		if isVideoFormat(format) {
			videoFormats = append(videoFormats, format)
		}
	}

	if len(videoFormats) == 0 {
		return nil, fmt.Errorf("no downloadable audio+video format found")
	}

	targetHeight := parseTargetHeight(quality)
	if targetHeight < 1080 {
		targetHeight = 1080
	}

	candidates := make([]youtube.Format, 0, len(videoFormats))
	for _, format := range videoFormats {
		if format.Height >= targetHeight {
			candidates = append(candidates, format)
		}
	}

	if len(candidates) == 0 {
		candidates = append(candidates, videoFormats...)
	}

	best := candidates[0]
	for i := 1; i < len(candidates); i++ {
		if betterFormat(candidates[i], best) {
			best = candidates[i]
		}
	}

	return &best, nil
}

func getBestAudioFormat(video *youtube.Video) (*youtube.Format, error) {
	audioFormats := youtube.FormatList{}
	for _, format := range video.Formats {
		if isAudioFormat(format) {
			audioFormats = append(audioFormats, format)
		}
	}

	if len(audioFormats) == 0 {
		return nil, fmt.Errorf("no downloadable audio format found")
	}

	best := audioFormats[0]
	for i := 1; i < len(audioFormats); i++ {
		if audioFormats[i].Bitrate > best.Bitrate {
			best = audioFormats[i]
		}
	}

	return &best, nil
}

func isVideoFormat(format youtube.Format) bool {
	parsed, _, err := mime.ParseMediaType(format.MimeType)
	if err != nil {
		return false
	}

	return strings.HasPrefix(parsed, "video/")
}

func isAudioFormat(format youtube.Format) bool {
	parsed, _, err := mime.ParseMediaType(format.MimeType)
	if err != nil {
		return false
	}

	return strings.HasPrefix(parsed, "audio/")
}

func downloadStreamToFile(client *youtube.Client, video *youtube.Video, format *youtube.Format, outputPath string) error {
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf("error getting stream: %w", err)
	}
	defer stream.Close()

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

func mergeVideoAndAudio(videoPath, audioPath, outputPath string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg is required for high-quality downloads with audio")
	}

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", videoPath,
		"-i", audioPath,
		"-c", "copy",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg merge failed: %v (%s)", err, strings.TrimSpace(string(output)))
	}

	return nil
}

func isShortURL(rawURL string) bool {
	parsed, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}

	host := strings.ToLower(parsed.Host)
	if host != "www.youtube.com" && host != "youtube.com" && host != "m.youtube.com" {
		return false
	}

	return strings.HasPrefix(strings.ToLower(parsed.Path), "/shorts/")
}

func betterFormat(a, b youtube.Format) bool {
	if a.Height != b.Height {
		return a.Height > b.Height
	}
	if a.Bitrate != b.Bitrate {
		return a.Bitrate > b.Bitrate
	}
	return a.AudioChannels > b.AudioChannels
}

func parseTargetHeight(quality string) int {
	quality = strings.TrimSpace(strings.ToLower(quality))
	if quality == "" {
		return 0
	}

	re := regexp.MustCompile(`\d+`)
	matches := re.FindString(quality)
	if matches == "" {
		return 0
	}

	value, err := strconv.Atoi(matches)
	if err != nil {
		return 0
	}
	return value
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
