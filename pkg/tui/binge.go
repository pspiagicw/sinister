package tui

import (
	"flag"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/help"
)

func parseBingeArgs(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister binge", flag.ExitOnError)
	flag.Usage = help.HelpBinge
	flag.Parse(opts.Args[1:])
}
func Binge(opts *argparse.Opts) {

	parseBingeArgs(opts)

	entries := database.QueryAll()

	names := []string{}

	for _, entry := range entries {
		names = append(names, entry.Title+" by "+entry.Author.Name)
	}

	selected := promptMultiple("Select video", names)

	for _, index := range selected {
		entry := entries[index]
		goreland.LogInfo("Downloading %s by %s", entry.Title, entry.Author.Name)
		performDownload(opts, &entry)
	}
}
