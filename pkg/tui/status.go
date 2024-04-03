package tui

import (
	"flag"
	"fmt"

	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/config"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/help"
)

func parseStatusArgs(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister status", flag.ExitOnError)
	flag.Usage = help.HelpStatus
	flag.Parse(opts.Args[1:])
}

func Status(opts *argparse.Opts) {

	parseStatusArgs(opts)
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
func printFeed(feed string, value config.Feed) {
	fmt.Println("Feed: ", feed)
	fmt.Println("URL: ", value.URL)
	fmt.Println("Tags: ")
	for _, tag := range value.Tags {
		fmt.Printf("- %s\n", tag)
	}
	fmt.Println()
}
