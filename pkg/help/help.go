package help

import (
	"fmt"

	"github.com/pspiagicw/goreland"
)

func Usage(version string) {
	goreland.LogFatal("Help printing not implemented yet!")
}
func Version(version string) {
	fmt.Println("sinister version: '%s'", version)
}
func HelpConfig() {
	fmt.Println("HELP CONFIG NOT IMPLEMENTED YET!")
}
