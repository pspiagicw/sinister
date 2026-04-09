package manage

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

func Update(configPath string) {
	conf := config.ParseConfig(configPath)

	for _, url := range conf.URLS {
		goreland.LogInfo("Fetching %s", url)

		f := fetchFeed(url)
		for _, entry := range f.Entries {
			// Some feeds only define author at the feed level.
			if entry.Author.Name == "" {
				entry.Author = f.Author
			}
			database.Add(&entry)
		}
	}
}

func fetchFeed(url string) *feed.Feed {
	body := getContents(url)

	var f feed.Feed
	if err := xml.Unmarshal(body, &f); err != nil {
		goreland.LogFatal("Error while parsing feed: %v", err)
	}

	return &f
}

func getContents(url string) []byte {
	client := http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
	defer closeResponse(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		goreland.LogFatal("Error while fetching feed: %s", resp.Status)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		goreland.LogFatal("Error reading response body: %v", err)
	}

	return contents
}

func closeResponse(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		goreland.LogError("Error closing response body: %v", err)
	}
}
