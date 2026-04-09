package manage

import (
	"net/http"
	"time"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
)

type CleanOptions struct {
	ConfigPath string
	DryRun     bool
}

func Clean(opts CleanOptions) {
	conf, configPath := config.ReadConfigForUpdate(opts.ConfigPath)

	client := http.Client{Timeout: 20 * time.Second}
	kept := make([]string, 0, len(conf.URLS))
	removed := make([]string, 0)

	for _, feedURL := range conf.URLS {
		statusCode, err := getStatusCode(&client, feedURL)
		if err != nil {
			goreland.LogError("Failed to check %s: %v (keeping)", feedURL, err)
			kept = append(kept, feedURL)
			continue
		}

		if statusCode == http.StatusNotFound {
			removed = append(removed, feedURL)
			goreland.LogInfo("Removing 404 feed: %s", feedURL)
			continue
		}

		kept = append(kept, feedURL)
	}

	if opts.DryRun {
		goreland.LogSuccess("[dry-run] Would remove %d feed(s)", len(removed))
		return
	}

	conf.URLS = kept
	config.WriteConfig(configPath, conf)
	goreland.LogSuccess("Removed %d feed(s), kept %d", len(removed), len(kept))
}

func getStatusCode(client *http.Client, url string) (int, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
