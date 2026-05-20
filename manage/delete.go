package manage

import (
	"os"
	"strings"
	"time"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

type DeleteOptions struct {
	Creator       string
	Slugs         []string
	Days          int
	DryRun        bool
	MarkUnwatched bool
}

func Delete(opts DeleteOptions) {
	if !hasDeleteSelectionFlags(opts) {
		runInteractiveDelete(opts.DryRun, opts.MarkUnwatched)
		return
	}

	entries := database.QueryDownloaded()
	targets := filterDeleteTargets(entries, opts)

	if len(targets) == 0 {
		goreland.LogInfo("No downloaded videos matched the provided filters")
		return
	}

	applyDelete(targets, opts.DryRun, opts.MarkUnwatched)
}

func hasDeleteSelectionFlags(opts DeleteOptions) bool {
	return opts.Creator != "" || len(opts.Slugs) > 0 || opts.Days > 0
}

func runInteractiveDelete(dryRun, markUnwatched bool) {
	creators := database.QueryDownloadedCreators()
	if len(creators) == 0 {
		goreland.LogFatal("No downloaded videos to delete")
	}

	selected := promptSelection("Select creator", creators)
	creator := creators[selected]

	entries := database.QueryDownloadedByCreator(creator)
	titles := make([]string, len(entries))
	for i, e := range entries {
		titles[i] = e.Title
	}

	selectedIndexes := promptMultiple("Select videos to delete", titles)
	targets := make([]feed.Entry, 0, len(selectedIndexes))
	for _, idx := range selectedIndexes {
		targets = append(targets, entries[idx])
	}

	applyDelete(targets, dryRun, markUnwatched)
}

func filterDeleteTargets(entries []feed.Entry, opts DeleteOptions) []feed.Entry {
	slugSet := makeSet(opts.Slugs)
	creator := strings.TrimSpace(opts.Creator)

	var cutoff time.Time
	if opts.Days > 0 {
		cutoff = time.Now().AddDate(0, 0, -opts.Days)
	}

	var results []feed.Entry
	for _, entry := range entries {
		if creator != "" && entry.Author.Name != creator {
			continue
		}
		if len(slugSet) > 0 && !slugSet[entry.Slug] {
			continue
		}
		if opts.Days > 0 {
			t, ok := parsePublished(entry.Published)
			if !ok || !t.Before(cutoff) {
				continue
			}
		}
		results = append(results, entry)
	}
	return results
}

func applyDelete(entries []feed.Entry, dryRun, markUnwatched bool) {
	if dryRun {
		for _, entry := range entries {
			goreland.LogInfo("[dry-run] Delete: %s by %s (%s)", entry.Title, entry.Author.Name, entry.FilePath)
		}
		goreland.LogSuccess("[dry-run] %d videos would be deleted", len(entries))
		return
	}

	successCount := 0
	failedCount := 0

	for _, entry := range entries {
		if err := os.Remove(entry.FilePath); err != nil && !os.IsNotExist(err) {
			goreland.LogError("Failed to delete %s: %v", entry.FilePath, err)
			failedCount++
			continue
		}

		database.ClearFilePath(entry.Slug)
		if markUnwatched {
			entryCopy := entry
			database.UpdateUnwatched(&entryCopy)
		}

		goreland.LogInfo("Deleted: %s", entry.Title)
		successCount++
	}

	goreland.LogSuccess("Deletion complete. success=%d failed=%d", successCount, failedCount)
}
