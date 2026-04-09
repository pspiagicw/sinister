package manage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/config"
	"github.com/pspiagicw/sinister/database"
)

type StatusOptions struct {
	ConfigPath string
	JSON       bool
	Creator    string
}

type StatusSummary struct {
	VideoFolder   string   `json:"videoFolder"`
	URLs          []string `json:"urls"`
	Creator       string   `json:"creator,omitempty"`
	TotalVideos   int      `json:"totalVideos"`
	Unwatched     int      `json:"unwatchedVideos"`
	Watched       int      `json:"watchedVideos"`
	TotalCreators int      `json:"totalCreators,omitempty"`
}

func Status(opts StatusOptions) {
	conf := config.ParseConfig(opts.ConfigPath)
	summary := buildStatusSummary(conf, opts.Creator)

	if opts.JSON {
		printStatusJSON(summary)
		return
	}

	printStatusText(summary)
}

func buildStatusSummary(conf *config.Config, creator string) StatusSummary {
	summary := StatusSummary{
		VideoFolder: conf.VideoFolder,
		URLs:        conf.URLS,
		Creator:     creator,
	}

	if creator == "" {
		summary.TotalVideos = database.TotalEntries()
		summary.Unwatched = database.CountUnwatched()
		summary.TotalCreators = database.TotalCreators()
	} else {
		summary.TotalVideos = database.CountEntriesByCreator(creator)
		summary.Unwatched = database.CountUnwatchedByCreator(creator)
	}
	summary.Watched = summary.TotalVideos - summary.Unwatched

	return summary
}

func printStatusJSON(summary StatusSummary) {
	payload, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		goreland.LogFatal("Error while serializing status JSON: %v", err)
	}

	if _, err := fmt.Fprintln(os.Stdout, string(payload)); err != nil {
		goreland.LogFatal("Error while writing status output: %v", err)
	}
}

func printStatusText(summary StatusSummary) {
	fmt.Println()
	fmt.Println("SINISTER")
	fmt.Println("Video: ", summary.VideoFolder)
	fmt.Println()
	fmt.Println("URLs: ")

	for _, url := range summary.URLs {
		fmt.Printf("- %s\n", url)
	}

	fmt.Println()
	if summary.Creator != "" {
		fmt.Println("Creator: ", summary.Creator)
	}
	fmt.Println("Total Videos: ", summary.TotalVideos)
	if summary.TotalCreators > 0 {
		fmt.Println("Total Creators: ", summary.TotalCreators)
	}
	fmt.Println("Unwatched Videos: ", summary.Unwatched)
	fmt.Println("Watched Videos: ", summary.Watched)
}
