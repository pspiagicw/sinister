package argparse

import (
	"flag"

	"github.com/pspiagicw/sinister/pkg/help"
)

type Opts struct {
	Version string
	Args    []string

	// Alternate config file
	Config string

	// Subcommand fields
	NoSpinner bool

	// Auto fields
	Days        int
	NoSync      bool
	MarkWatched bool
	Format      bool
}

func Parse(version string) *Opts {

	opts := &Opts{}

	flag.Usage = func() {
		help.Usage(version)
	}

	flag.StringVar(&opts.Config, "config", "", "Path to the configuration file")

	flag.Parse()
	opts.Args = flag.Args()
	opts.Version = version

	return opts
}
