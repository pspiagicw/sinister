package manage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
)

const (
	youtubeFeedURLTemplate = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
)

var (
	reChannelIDInPath = regexp.MustCompile(`youtube\.com/channel/(UC[\w-]+)`)
	reChannelIDJSON   = regexp.MustCompile(`"channelId":"(UC[\w-]+)"`)
	reExternalIDJSON  = regexp.MustCompile(`"externalId":"(UC[\w-]+)"`)
)

type AddOptions struct {
	ConfigPath string
	VideoURL   string
}

type oEmbedResponse struct {
	AuthorURL string `json:"author_url"`
}

func Add(opts AddOptions) {
	channelID, err := resolveChannelID(opts.VideoURL)
	if err != nil {
		goreland.LogFatal("Error resolving channel: %v", err)
	}

	feedURL := fmt.Sprintf(youtubeFeedURLTemplate, channelID)
	conf, confPath := config.ReadConfigForUpdate(opts.ConfigPath)

	for _, existing := range conf.URLS {
		if strings.TrimSpace(existing) == feedURL {
			goreland.LogInfo("Feed already exists in config: %s", feedURL)
			return
		}
	}

	conf.URLS = append(conf.URLS, feedURL)
	config.WriteConfig(confPath, conf)

	goreland.LogSuccess("Added feed URL to config: %s", feedURL)
}

func resolveChannelID(videoURL string) (string, error) {
	if id, err := resolveChannelIDFromOEmbed(videoURL); err == nil && id != "" {
		return id, nil
	}

	body, err := fetchURL(videoURL)
	if err != nil {
		return "", err
	}

	if id := extractChannelID(body); id != "" {
		return id, nil
	}

	return "", fmt.Errorf("unable to resolve channel ID from URL")
}

func resolveChannelIDFromOEmbed(videoURL string) (string, error) {
	oembedURL := "https://www.youtube.com/oembed?url=" + neturl.QueryEscape(videoURL) + "&format=json"
	body, err := fetchURL(oembedURL)
	if err != nil {
		return "", err
	}

	var payload oEmbedResponse
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return "", err
	}

	if strings.TrimSpace(payload.AuthorURL) == "" {
		return "", fmt.Errorf("missing author url in oEmbed response")
	}

	channelPage, err := fetchURL(payload.AuthorURL)
	if err != nil {
		return "", err
	}

	if id := extractChannelID(channelPage); id != "" {
		return id, nil
	}

	return "", fmt.Errorf("unable to resolve channel id from author page")
}

func fetchURL(rawURL string) (string, error) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("request failed with status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

func extractChannelID(body string) string {
	if matches := reExternalIDJSON.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	if matches := reChannelIDInPath.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	if matches := reChannelIDJSON.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	return ""
}
