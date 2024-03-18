package tui

import (
	"fmt"

	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/config"
	"github.com/pspiagicw/sinister/pkg/database"
)

func Status(opts *argparse.Opts) {
	conf := config.ParseConfig(opts)

	fmt.Println()
	fmt.Println("SINISTER")
	fmt.Println("Video: ", conf.VideoFolder)
	fmt.Println()
	fmt.Println("URLs: ")

	for _, url := range conf.URLS {
		fmt.Printf("- %s\n", url)
	}

	fmt.Println()
	totalEntries := database.TotalEntries()
	totalCreators := database.TotalCreators()
	fmt.Println("Total Videos: ", totalEntries)
	fmt.Println("Total Creators: ", totalCreators)

	unwatchedEntries := database.UnwatchedEntries()

	fmt.Println("Unwatched Videos: ", unwatchedEntries)
}
