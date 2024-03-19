package tui

import (
	"flag"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/database"
)

func parseAutoOpts(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister auto", flag.ExitOnError)
	flag.Parse(opts.Args[1:])
}

func Auto(opts *argparse.Opts) {

	parseAutoOpts(opts)

	goreland.LogInfo("Starting auto mode!!")
	goreland.LogInfo("Synching with feeds...")
	performUpdate(opts)
	goreland.LogSuccess("Synching with feeds completed!")
	goreland.LogInfo("Downloading videos...")

	entries := database.QueryAll()

	for _, entry := range entries {
		if entry.Watched == 0 {
			goreland.LogInfo("Downloading %s", entry.Title)
			performDownload(opts, &entry)
		}
	}

	goreland.LogSuccess("All downloads completed!")
}
