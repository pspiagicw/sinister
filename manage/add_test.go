package manage

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pspiagicw/sinister/config"
)

func TestExtractChannelID_ExternalID(t *testing.T) {
	body := `some html "externalId":"UCxxxxYYYYzzzz12" more html`
	got := extractChannelID(body)
	if got != "UCxxxxYYYYzzzz12" {
		t.Errorf("expected UCxxxxYYYYzzzz12, got %q", got)
	}
}

func TestExtractChannelID_ChannelInPath(t *testing.T) {
	body := `<link href="https://www.youtube.com/channel/UCabcDEF123456">`
	got := extractChannelID(body)
	if got != "UCabcDEF123456" {
		t.Errorf("expected UCabcDEF123456, got %q", got)
	}
}

func TestExtractChannelID_ChannelIDJSON(t *testing.T) {
	body := `{"channelId":"UCfooBarBaz999"}`
	got := extractChannelID(body)
	if got != "UCfooBarBaz999" {
		t.Errorf("expected UCfooBarBaz999, got %q", got)
	}
}

func TestExtractChannelID_NotFound(t *testing.T) {
	got := extractChannelID("no channel info here")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestAdd_AppendsToConfig(t *testing.T) {
	channelID := "UCtest1234567890"

	// Serve a fake channel page that contains the channel ID
	channelSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html>"externalId":"%s"</html>`, channelID)
	}))
	defer channelSrv.Close()

	// Serve a fake oEmbed response pointing to the channel page
	oembedSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"author_url":"%s"}`, channelSrv.URL)
	}))
	defer oembedSrv.Close()

	// Write a minimal config
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	content := `videoFolder = "/tmp/videos"
quality = "hd720"
urls = []
`
	os.WriteFile(configPath, []byte(content), 0644)

	// Patch the oEmbed URL by pointing the videoURL to the oEmbed server.
	// resolveChannelIDFromOEmbed constructs the oEmbed URL internally, so we
	// test through the fallback path (fetchURL on videoURL directly) instead.
	videoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html>"externalId":"%s"</html>`, channelID)
	}))
	defer videoSrv.Close()

	Add(AddOptions{
		ConfigPath: configPath,
		VideoURL:   videoSrv.URL,
	})

	conf, _ := config.ReadConfigForUpdate(configPath)
	found := false
	for _, u := range conf.URLS {
		if strings.Contains(u, channelID) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected feed URL with channelID %q in config, got %v", channelID, conf.URLS)
	}
}

func TestAdd_SkipsDuplicate(t *testing.T) {
	channelID := "UCduplicateTest123"

	videoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"externalId":"%s"`, channelID)
	}))
	defer videoSrv.Close()

	feedURL := fmt.Sprintf(youtubeFeedURLTemplate, channelID)
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	content := fmt.Sprintf(`videoFolder = "/tmp/videos"
quality = "hd720"
urls = ["%s"]
`, feedURL)
	os.WriteFile(configPath, []byte(content), 0644)

	Add(AddOptions{ConfigPath: configPath, VideoURL: videoSrv.URL})

	conf, _ := config.ReadConfigForUpdate(configPath)
	count := 0
	for _, u := range conf.URLS {
		if u == feedURL {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 occurrence of feed URL, got %d", count)
	}
}
