package tui

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pspiagicw/goreland"
)

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
func promptSelection(label string, options []string) string {
	prompt := &survey.Select{
		Message: label,
		Options: options,
	}

	var selected string
	survey.AskOne(prompt, &selected)

	return selected
}
func promptMultiple(label string, options []string) []string {
	choices := make([]string, 0)

	prompt := &survey.MultiSelect{
		Message: label,
		Options: options,
	}
	survey.AskOne(prompt, &choices)

	if len(choices) == 0 {
		goreland.LogFatal("No option selected")
	}

	return choices
}
