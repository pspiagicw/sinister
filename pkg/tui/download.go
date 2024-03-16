package tui

import (
	"io"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/manifoldco/promptui"
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
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

	videoID := getVideoID(entry.Link.URL)

	downloadVideo(videoID, opts)

	database.UpdateWatched(entry)
}

func downloadVideo(videoID string, opts *argparse.Opts) {

	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		goreland.LogFatal("Error getting video: %s", err)
	}

	formats := video.Formats.Quality("720p").WithAudioChannels()

	formats.Sort()

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		goreland.LogFatal("Error getting stream: %s", err)
	}

	defer stream.Close()

	fileName := videoID + ".mp4"

	file, err := os.Create(fileName)
	if err != nil {
		goreland.LogFatal("Error creating file: %s", err)
	}

	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		goreland.LogFatal("Error copying stream to file: %s", err)
	}
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
