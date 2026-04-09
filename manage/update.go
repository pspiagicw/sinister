package manage

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

type UpdateOptions struct {
	ConfigPath string
	URLs       []string
	Limit      int
	SinceDays  int
	Retries    int
	Timeout    int
	DryRun     bool
	JSON       bool
}

type UpdateSummary struct {
	FeedsProcessed int  `json:"feedsProcessed"`
	FeedsSkipped   int  `json:"feedsSkipped"`
	EntriesSeen    int  `json:"entriesSeen"`
	FilteredOut    int  `json:"filteredOut"`
	Inserted       int  `json:"inserted"`
	Skipped        int  `json:"skipped"`
	DryRun         bool `json:"dryRun"`
}

func Update(opts UpdateOptions) {
	urls := resolveUpdateURLs(opts)
	summary := UpdateSummary{DryRun: opts.DryRun}

	for _, url := range urls {
		goreland.LogInfo("Fetching %s", url)
		summary.FeedsProcessed++

		f, err := fetchFeedWithRetry(url, opts)
		if err != nil {
			summary.FeedsSkipped++
			goreland.LogError("Skipping feed %s: %v", url, err)
			continue
		}
		entries := applyLimit(f.Entries, opts.Limit)

		for _, entry := range entries {
			if entry.Author.Name == "" {
				entry.Author = f.Author
			}

			if !withinSinceDays(entry.Published, opts.SinceDays) {
				summary.FilteredOut++
				continue
			}

			summary.EntriesSeen++
			if opts.DryRun {
				if database.ExistsByTitle(entry.Title) {
					summary.Skipped++
				} else {
					summary.Inserted++
				}
				continue
			}

			if database.Add(&entry) {
				summary.Inserted++
			} else {
				summary.Skipped++
			}
		}
	}

	printUpdateSummary(summary, opts.JSON)
}

func resolveUpdateURLs(opts UpdateOptions) []string {
	if len(opts.URLs) > 0 {
		return opts.URLs
	}

	conf := config.ParseConfig(opts.ConfigPath)
	return conf.URLS
}

func applyLimit(entries []feed.Entry, limit int) []feed.Entry {
	if limit <= 0 || len(entries) <= limit {
		return entries
	}
	return entries[:limit]
}

func printUpdateSummary(summary UpdateSummary, asJSON bool) {
	if asJSON {
		payload, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			goreland.LogFatal("Error while serializing update JSON: %v", err)
		}
		if _, err := fmt.Fprintln(os.Stdout, string(payload)); err != nil {
			goreland.LogFatal("Error while writing update output: %v", err)
		}
		return
	}

	goreland.LogSuccess(
		"Update complete: feeds=%d skipped-feeds=%d entries=%d filtered=%d inserted=%d skipped=%d dry-run=%t",
		summary.FeedsProcessed,
		summary.FeedsSkipped,
		summary.EntriesSeen,
		summary.FilteredOut,
		summary.Inserted,
		summary.Skipped,
		summary.DryRun,
	)
}

func fetchFeedWithRetry(url string, opts UpdateOptions) (*feed.Feed, error) {
	attempts := opts.Retries + 1
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		f, err := fetchFeed(url, opts.Timeout)
		if err == nil {
			return f, nil
		}

		lastErr = err
		if attempt < attempts {
			goreland.LogError("Fetch failed for %s (attempt %d/%d): %v", url, attempt, attempts, err)
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	return nil, lastErr
}

func fetchFeed(url string, timeoutSec int) (*feed.Feed, error) {
	body, err := getContents(url, timeoutSec)
	if err != nil {
		return nil, err
	}

	var f feed.Feed
	if err := xml.Unmarshal(body, &f); err != nil {
		return nil, fmt.Errorf("error while parsing feed: %w", err)
	}

	return &f, nil
}

func getContents(url string, timeoutSec int) ([]byte, error) {
	timeout := 30 * time.Second
	if timeoutSec > 0 {
		timeout = time.Duration(timeoutSec) * time.Second
	}
	client := http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while connecting: %w", err)
	}
	defer closeResponse(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return contents, nil
}

func withinSinceDays(published string, sinceDays int) bool {
	if sinceDays <= 0 {
		return true
	}

	t, err := time.Parse(time.RFC3339, published)
	if err != nil {
		return false
	}

	cutoff := time.Now().AddDate(0, 0, -sinceDays)
	return t.After(cutoff)
}

func closeResponse(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		goreland.LogError("Error closing response body: %v", err)
	}
}
