package help

import (
	"fmt"

	"github.com/pspiagicw/pelp"
)

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
	pelp.HeaderWithDescription("more help", []string{"Use 'sinister help [command]' for more info about a command."})
	pelp.Examples("examples", []string{"sinister update", "sinister download"})

}
func Version(version string) {
	pelp.Version("sinister", version)
}
func HelpConfig() {
	fmt.Println("HELP CONFIG NOT IMPLEMENTED YET!")
}
func HelpUpdate() {
	pelp.Print("Update subscriptions")
	pelp.HeaderWithDescription("Usage", []string{"sinister update"})
	pelp.HeaderWithDescription("Description", []string{"Update subscriptions according to the RSS feeds."})
}
func HelpDownload() {
	pelp.Print("Download a video")
	pelp.HeaderWithDescription("Usage", []string{"sinister download"})
	pelp.Flags("flag", []string{"no-spinner"}, []string{"Disable spinner"})
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
		case "binge":
			HelpBinge()
		default:
			Usage(version)
		}
	}
}
