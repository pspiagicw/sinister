package tui

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kkdai/youtube/v2"
	"github.com/manifoldco/promptui"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/config"
	"github.com/pspiagicw/sinister/pkg/database"
	"github.com/pspiagicw/sinister/pkg/feed"
)

func Download(opts *argparse.Opts) {

	entry := selectEntry()

	performDownload(opts, entry)

}

func selectCreator() string {
	creators := database.QueryCreators()

	if len(creators) == 0 {
		goreland.LogFatal("No creators with unwatched videos")
	}

	selected := promptSelection("Select creator", creators)

	return selected

}
func selectVideo(selected string) string {
	videos := database.QueryVideos(selected)

	if len(videos) == 0 {
		goreland.LogFatal("No unwatched videos for creator")
	}

	selectedVideo := promptSelection("Select video", videos)

	return selectedVideo
}

func selectEntry() *feed.Entry {

	creator := selectCreator()

	video := selectVideo(creator)

	entry := database.QueryEntry(creator, video)

	return entry
}

func performDownload(opts *argparse.Opts, entry *feed.Entry) {

	confirmDownload()
	downloadVideo(entry, opts)

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
func selectFormat(video *youtube.Video) *youtube.Format {
	formats := video.Formats.Quality("720p").WithAudioChannels()

	formats.Sort()

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

	format := selectFormat(video)

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
func promptSelection(label string, creators []string) string {
	prompt := promptui.Select{Label: label, Items: creators}

	_, value, err := prompt.Run()

	if err != nil {
		goreland.LogFatal("Something went wrong: %q", err)
	}

	return value
}
func getDownloadPath(opts *argparse.Opts, entry *feed.Entry) string {

	conf := config.ParseConfig(opts)

	return filepath.Join(conf.VideoFolder, entry.Slug+".mp4")
}

func confirmDownload() {
	confirm := false
	prompt := survey.Confirm{
		Message: "Download the video?",
	}
	survey.AskOne(&prompt, &confirm)
	if !confirm {
		goreland.LogFatal("User cancelled the download.")
	}
}
