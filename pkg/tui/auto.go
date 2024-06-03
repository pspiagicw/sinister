package tui

import (
	"flag"
	"fmt"
	"time"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/feed"
	"github.com/pspiagicw/sinister/pkg/help"
)

func parseAutoOpts(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister auto", flag.ExitOnError)
	flag.IntVar(&opts.Days, "days", 0, "Maximum number of days to download")
	flag.BoolVar(&opts.NoSync, "no-sync", false, "Disable spinner")
	flag.BoolVar(&opts.MarkWatched, "mark-watched", false, "Mark videos as watched when rejected")
	flag.Usage = help.HelpAuto
	flag.Parse(opts.Args[1:])
}

func Auto(opts *argparse.Opts) {

	parseAutoOpts(opts)

	goreland.LogInfo("Starting auto mode!!")
	goreland.LogInfo("Synching with feeds...")

	if !opts.NoSync {
		performUpdate(opts)
	}

	goreland.LogSuccess("Synching with feeds completed!")
	goreland.LogInfo("Downloading videos...")

	entries := database.QueryAll()

	entries = filterEntries(entries, opts)

	goreland.LogInfo("%d videos matched your filter", len(entries))

	for _, entry := range entries {

		if softConfirm(fmt.Sprintf("Download %s by %s ?", entry.Title, entry.Author.Name)) {
			performDownload(opts, &entry)
		} else {
			if opts.MarkWatched {
				database.UpdateWatched(&entry)
			}
		}
	}

	goreland.LogSuccess("All downloads completed!")
}
func filterEntries(entries []feed.Entry, opts *argparse.Opts) []feed.Entry {
	filtered := make([]feed.Entry, 0)

	for _, entry := range entries {

		t := entry.Date()

		if validDate(t, opts) && entry.Watched == 0 {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func validDate(t time.Time, opts *argparse.Opts) bool {
	now := time.Now()

	nDaysAgo := now.AddDate(0, 0, -opts.Days)

	if t.After(nDaysAgo) {
		return true
	}

	return false
}
