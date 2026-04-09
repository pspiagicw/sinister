package manage

import (
	"strings"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/database"
	"github.com/pspiagicw/sinister/feed"
)

type MarkOptions struct {
	Slugs        []string
	URLs         []string
	Creator      string
	AllUnwatched bool
	DryRun       bool
}

func Mark(opts MarkOptions) {
	if !hasSelectionFlags(opts) {
		runInteractiveMark(opts.DryRun)
		return
	}

	entries := database.QueryUnwatched()
	targets := selectMarkTargets(entries, opts)

	if len(targets) == 0 {
		goreland.LogInfo("No unwatched videos matched the provided filters")
		return
	}

	applyMark(targets, opts.DryRun)
}

func hasSelectionFlags(opts MarkOptions) bool {
	return len(opts.Slugs) > 0 || len(opts.URLs) > 0 || opts.Creator != "" || opts.AllUnwatched
}

func runInteractiveMark(dryRun bool) {
	creator := selectCreatorForMark()
	videos := database.QueryVideos(creator)

	selectedVideos := promptMultiple("Select videos to mark watched", videos)
	targets := make([]feed.Entry, 0, len(selectedVideos))

	for _, index := range selectedVideos {
		entry := database.QueryEntry(creator, videos[index])
		targets = append(targets, *entry)
	}

	applyMark(targets, dryRun)
}

func selectCreatorForMark() string {
	creators := database.QueryCreators()
	if len(creators) == 0 {
		goreland.LogFatal("No creators with unwatched videos")
	}

	selected := promptSelection("Select creator", creators)
	return creators[selected]
}

func selectMarkTargets(entries []feed.Entry, opts MarkOptions) []feed.Entry {
	resultBySlug := map[string]feed.Entry{}
	slugSet := makeSet(opts.Slugs)
	urlSet := makeSet(opts.URLs)
	creator := strings.TrimSpace(opts.Creator)

	for _, entry := range entries {
		if creator != "" && entry.Author.Name != creator {
			continue
		}

		shouldInclude := false
		if opts.AllUnwatched {
			shouldInclude = true
		}
		if len(slugSet) > 0 && slugSet[entry.Slug] {
			shouldInclude = true
		}
		if len(urlSet) > 0 && urlSet[entry.Link.URL] {
			shouldInclude = true
		}
		if !opts.AllUnwatched && len(slugSet) == 0 && len(urlSet) == 0 {
			shouldInclude = true
		}

		if shouldInclude {
			resultBySlug[entry.Slug] = entry
		}
	}

	results := make([]feed.Entry, 0, len(resultBySlug))
	for _, entry := range entries {
		if selected, ok := resultBySlug[entry.Slug]; ok {
			results = append(results, selected)
		}
	}

	return results
}

func makeSet(values []string) map[string]bool {
	result := map[string]bool{}
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result[trimmed] = true
		}
	}
	return result
}

func applyMark(entries []feed.Entry, dryRun bool) {
	if dryRun {
		for _, entry := range entries {
			goreland.LogInfo("[dry-run] Mark watched: %s by %s", entry.Title, entry.Author.Name)
		}
		goreland.LogSuccess("[dry-run] %d videos would be marked watched", len(entries))
		return
	}

	for _, entry := range entries {
		entryCopy := entry
		database.UpdateWatched(&entryCopy)
		goreland.LogInfo("Marked watched: %s by %s", entry.Title, entry.Author.Name)
	}
	goreland.LogSuccess("Marked %d videos as watched", len(entries))
}
