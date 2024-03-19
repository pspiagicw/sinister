package tui

import (
	"flag"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/help"
)

func parseMarkArgs(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister mark", flag.ExitOnError)
	flag.Usage = help.HelpMark
	flag.Parse(opts.Args[1:])
}

func Mark(opts *argparse.Opts) {
	parseMarkArgs(opts)

	creator := selectCreator()

	videos := database.QueryVideos(creator)

	selectedVideos := promptMultiple("Select videos to mark watched", videos)

	for _, video := range selectedVideos {
		entry := database.QueryEntry(creator, video)
		goreland.LogInfo("Marking %s by %s as watched", entry.Title, entry.Author.Name)
		database.UpdateWatched(entry)
	}
}
