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
		[]string{"update:", "download:", "status:", "help:"},
		[]string{"Update subscriptions", "Download a video", "Show status", "Show this help message"})
	pelp.HeaderWithDescription("more help", []string{"Use 'sinister help [command]' for more info about a command."})
	pelp.Examples("examples", []string{"sinister update", "sinister download"})

}
func Version(version string) {
	fmt.Printf("sinister version: '%s'\n", version)
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
