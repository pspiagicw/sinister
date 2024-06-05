package tui

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kkdai/youtube/v2"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/config"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/feed"
	"github.com/pspiagicw/sinister/pkg/help"
)

func parseDownloadArgs(opts *argparse.Opts) {
	flag := flag.NewFlagSet("sinister download", flag.ExitOnError)
	flag.BoolVar(&opts.NoSpinner, "no-spinner", false, "Disable spinner")
	flag.BoolVar(&opts.Format, "select-format", false, "Select format manually")
	flag.Usage = help.HelpDownload
	flag.Parse(opts.Args[1:])
}

func Download(opts *argparse.Opts) {

	parseDownloadArgs(opts)

	entry := selectEntry()

	confirmDownload()

	performDownload(opts, entry)

}

func selectCreator() string {

	creators := database.QueryCreators()

	if len(creators) == 0 {
		goreland.LogFatal("No creators with unwatched videos")
	}

	selected := promptSelection("Select creator", creators)

	return creators[selected]

}
func selectVideo(selected string) string {
	videos := database.QueryVideos(selected)

	if len(videos) == 0 {
		goreland.LogFatal("No unwatched videos for creator")
	}

	selectedVideo := promptSelection("Select video", videos)

	return videos[selectedVideo]
}

func selectEntry() *feed.Entry {

	creator := selectCreator()

	video := selectVideo(creator)

	entry := database.QueryEntry(creator, video)

	return entry
}

func startSpinner(opts *argparse.Opts) func() {

	if opts.NoSpinner || opts.Format {
		return func() {}
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	s.Start()

	return func() {
		s.Stop()
	}
}

func performDownload(opts *argparse.Opts, entry *feed.Entry) {

	stopFunc := startSpinner(opts)

	downloadVideo(entry, opts)

	stopFunc()

	goreland.LogSuccess("Download complete")

	database.UpdateWatched(entry)
}
func getVideo(entry *feed.Entry) (*youtube.Client, *youtube.Video) {

	videoID := getVideoID(entry.Link.URL)

	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		goreland.LogFatal("Error getting video: %s", err)
	}

	return &client, video
}
func formatListString(formats youtube.FormatList) []string {
	results := []string{}
	for _, format := range formats {
		line := fmt.Sprintf("Quality: %s, Bitrate: %d", format.Quality, format.Bitrate)
		results = append(results, line)
	}

	return results
}
func selectFormat(video *youtube.Video, opts *argparse.Opts) *youtube.Format {
	formats := video.Formats

	formats.Sort()

	if opts.Format {
		index := promptSelection("Select format", formatListString(formats))
		return &formats[index]
	}

	conf := config.ParseConfig(opts)

	formatQuery := conf.Quality

	formats = formats.Quality(formatQuery).WithAudioChannels()

	if len(formats) == 0 {
		goreland.LogFatal("No suitable formats found")
	}

	return &formats[0]
}
func getStream(client *youtube.Client, video *youtube.Video, format *youtube.Format) io.ReadCloser {
	stream, _, err := client.GetStream(video, format)

	if err != nil {
		goreland.LogFatal("Error getting stream: %s", err)
	}

	return stream
}

func downloadVideo(entry *feed.Entry, opts *argparse.Opts) {

	client, video := getVideo(entry)

	format := selectFormat(video, opts)

	stream := getStream(client, video, format)

	defer stream.Close()

	fileName := getDownloadPath(opts, entry)

	file := openFile(fileName)

	defer file.Close()

	copyStreamToFile(stream, file)

}
func copyStreamToFile(stream io.ReadCloser, file *os.File) {
	_, err := io.Copy(file, stream)
	if err != nil {
		goreland.LogFatal("Error copying stream to file: %s", err)
	}
}
func openFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		goreland.LogFatal("Error creating file: %s", err)
	}
	return file
}

func getVideoID(url string) string {
	id, found := strings.CutPrefix(url, "https://www.youtube.com/watch?v=")
	if !found {
		goreland.LogFatal("Invalid URL")
	}
	return id
}
func getDownloadPath(opts *argparse.Opts, entry *feed.Entry) string {

	conf := config.ParseConfig(opts)

	return filepath.Join(conf.VideoFolder, entry.Slug+".mp4")
}
