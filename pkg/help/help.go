package help

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Usage(version string) {
	Version(version)
	fmt.Println("Watch youtube the Unix way")
	fmt.Println()
	fmt.Println("Usage: sinister [subcommand] [args]")
	fmt.Println()
	fmt.Println("COMMANDS")

	commands := `
update:
download:
status:
help:`

	messages := `
Update subscriptions
Download a video
Show status
Show this help message`

	commandCol := lipgloss.NewStyle().Align(lipgloss.Left).SetString(commands).MarginLeft(2).String()
	messageCol := lipgloss.NewStyle().Align(lipgloss.Left).SetString(messages).MarginLeft(5).String()

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Bottom, commandCol, messageCol))
	fmt.Println()
	fmt.Println("MORE HELP")
	fmt.Println("  Use 'sinister help [command]' for more info about a command.")
	fmt.Println()
	fmt.Println("EXAMPLES")
	fmt.Println("  $ sinister update")
	fmt.Println("  $ sinister download")
	fmt.Println()

}
func Version(version string) {
	fmt.Printf("sinister version: '%s'\n", version)
}
func HelpConfig() {
	fmt.Println("HELP CONFIG NOT IMPLEMENTED YET!")
}
func HelpUpdate() {
	fmt.Println("Update subscriptions")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  sinister update")
	fmt.Println()
	fmt.Println("DESCRIPTION")
	fmt.Println("  Update subscriptions according to the RSS feeds.")
	fmt.Println()
}
func HelpDownload() {
	fmt.Println("Download a video")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  sinister download")
	fmt.Println()
	fmt.Println("DESCRIPTION")
	fmt.Println("  Prompt for a creator and a video to download.")
	fmt.Println("  The video will be marked as watched, after successfull download.")
	fmt.Println()
}
func HelpStatus() {
	fmt.Println("Show the status of the subscriptions")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  sinister status")
	fmt.Println()
	fmt.Println("DESCRIPTION")
	fmt.Println("  Show statistics about the subscriptions.")
	fmt.Println()
}
