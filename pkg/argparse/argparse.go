package argparse

import (
	"flag"

	"github.com/pspiagicw/sinister/pkg/help"
)

func ParseArgs(VERSION string) []string {
	flag.Usage = func() {
		help.PrintHelp(VERSION)
	}

	flag.Parse()

	return flag.Args()
}
