package tui

import (
	"encoding/xml"
	"flag"
	"io"
	"net/http"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/config"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/feed"
	"github.com/pspiagicw/sinister/pkg/help"
)

func parseUpdateArgs(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister update", flag.ExitOnError)
	flag.Usage = help.HelpUpdate
	flag.Parse(opts.Args[1:])
}
func Update(opts *argparse.Opts) {
	parseUpdateArgs(opts)
	performUpdate(opts)
}
func performUpdate(opts *argparse.Opts) {
	conf := config.ParseConfig(opts)

	for _, url := range conf.GetURLs() {
		goreland.LogInfo("Fetching %s", url)
		fetch(url)
	}

}

func fetch(url string) {

	contents := getContents(url)

	feed := decodeFeed(contents)

	for _, entry := range feed.Entries {
		database.Add(&entry)
	}

}
func decodeFeed(contents []byte) *feed.Feed {
	var feed feed.Feed

	err := xml.Unmarshal(contents, &feed)
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
	return &feed
}

func getContents(url string) []byte {

	resp := makeRequest(url)
	contents := readResponse(resp)

	return contents
}
func makeRequest(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		goreland.LogFatal("Error while connecting: %v", err)
	}
	return resp
}
func readResponse(resp *http.Response) []byte {

	contents := readBody(resp)

	closeResponse(resp)

	return contents
}
func readBody(resp *http.Response) []byte {
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		goreland.LogFatal("Error reading response body: %v", err)
	}

	return contents
}
func closeResponse(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		goreland.LogError("Error closing response body: %v", err)
	}
}
