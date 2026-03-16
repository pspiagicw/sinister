package help

import (
	"github.com/pspiagicw/pelp"
)

const EXAMPLE_CONFIG = `
--------

videoFolder = "~/Videos"
urls = [
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCeeFfhMcJa1kjtfZAGskOCA",
    "https://www.youtube.com/feeds/videos.xml?channel_id=UCdBK94H6oZT2Q7l0-b0xmMg",
]

quality = "720"

----
`

func Usage(version string) {
	Version(version)
	pelp.Print("Watch youtube the Unix way!")
	pelp.HeaderWithDescription("Usage", []string{"sinister [subcommand] [args]"})
	pelp.Aligned("commands",
		[]string{"update:", "download:", "status:", "mark:", "auto:", "help:"},
		[]string{
			"Update subscriptions",
			"Download a video",
			"Show status",
			"Mark videos as watched (without downloading)",
			"Download unwatched videos automatically",
			"Show this help message",
		})
	pelp.Flags("flag", []string{"config"}, []string{"Provide alternate path to config file."})
	pelp.HeaderWithDescription("more help", []string{"Use 'sinister help [command]' for more info about a command."})
	pelp.Examples("examples", []string{"sinister update", "sinister download"})

}
func Version(version string) {
	pelp.Version("sinister", version)
}
func HelpConfig() {
	// fmt.Println("HELP CONFIG NOT IMPLEMENTED YET!")
	pelp.Print("Configure `sinister`")
	pelp.HeaderWithDescription("Path", []string{
		"The configuration will be searched at `~/.config/sinister/config.toml``",
		"Use the `--config` flag to provide an alternate path to the configuration file.",
	})
	pelp.HeaderWithDescription("Format", []string{
		"The configuration file is in TOML format.",
		"The configuration file should have the following fields",
	})

	pelp.Aligned(
		"Fields",
		[]string{
			"videoFolder:",
			"urls:",
			"quality:",
		},
		[]string{
			"The folder where the videos will be downloaded.",
			"A list of URLs to subscribe to.",
			"The quality of the video to download. (e.g. 720p, 1080p)",
		})

	pelp.Print("Example configuration file.")
	pelp.Print(EXAMPLE_CONFIG)
}
func HelpUpdate() {
	pelp.Print("Update subscriptions")
	pelp.HeaderWithDescription("Usage", []string{"sinister update"})
	pelp.HeaderWithDescription("Description", []string{"Update subscriptions according to the RSS feeds."})
}
func HelpDownload() {
	pelp.Print("Download a video")
	pelp.HeaderWithDescription("Usage", []string{"sinister download"})
	pelp.Flags("flag", []string{"no-spinner", "select-format"}, []string{"Disable spinner", "Select format manually"})
	pelp.HeaderWithDescription("Description",
		[]string{
			"Prompt for a creator and a video to download.",
			"The video will be marked as watched, after successfull download.",
		})
}
func HelpStatus() {
	pelp.Print("Show the status of the subscriptions")
	pelp.HeaderWithDescription("Usage", []string{"sinister status"})
	pelp.HeaderWithDescription("Description", []string{"Show statistics about the subscriptions."})
}
func HelpMark() {
	pelp.Print("Mark videos as watched")
	pelp.HeaderWithDescription("Usage", []string{"sinister mark"})
	pelp.HeaderWithDescription("Description",
		[]string{
			"Prompt for a creator and videos to mark as watched.",
			"The videos will be marked as watched.",
		})
}
func HelpAuto() {
	pelp.Print("Download unwatched videos automatically")
	pelp.HeaderWithDescription("Usage", []string{"sinister auto"})
	pelp.HeaderWithDescription("Description",
		[]string{
			"Download unwatched videos automatically.",
		})

	pelp.Flags(
		"flag",
		[]string{"days", "no-sync", "mark-watched"},
		[]string{"Maximum number of days to download", "Don't sync", "Mark videos as watched when rejected"})
}
func HelpBinge() {
	pelp.Print("Select a bunch of videos to download.")
	pelp.HeaderWithDescription("Usage", []string{"sinister binge"})
	pelp.HeaderWithDescription("Description",
		[]string{
			"Select a bunch of videos to download.",
		})
}
func HandleHelp(args []string, version string) {
	if len(args) == 0 {
		Usage(version)
	} else {
		switch args[0] {
		case "config":
			HelpConfig()
		case "update":
			HelpUpdate()
		case "download":
			HelpDownload()
		case "status":
			HelpStatus()
		case "mark":
			HelpMark()
		case "auto":
			HelpAuto()
		default:
			Usage(version)
		}
	}
}
