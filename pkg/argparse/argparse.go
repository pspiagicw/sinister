package argparse

import (
	"flag"

	"github.com/pspiagicw/sinister/pkg/help"
)

type Opts struct {
	Version string
	Args    []string
}

func Parse(version string) *Opts {

	opts := &Opts{}

	flag.Usage = func() {
		help.Usage(version)
	}

	flag.Parse()
	opts.Args = flag.Args()
	opts.Version = version

	return opts
}
