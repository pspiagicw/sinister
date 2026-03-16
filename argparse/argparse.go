package argparse

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/pspiagicw/sinister/manage"
)

type Opts struct {
	Version string
	// Args    []string
	//
	// // Alternate config file
	Config string
	//
	// // Subcommand fields
	// NoSpinner bool
	//
	// // Auto fields
	// Days        int
	// NoSync      bool
	// MarkWatched bool
	// Format      bool
}

type VersionCMD struct {
}

func (v *VersionCMD) Run(o *Opts) {
	fmt.Printf("%s version: '%s'\n", "sinister", o.Version)
}

type StatusCMD struct {
}

func (s *StatusCMD) Run(o *Opts) {
	manage.Status(o.Config)

}

type UpdateCMD struct {
}

type MarkCMD struct {
}

var CLI struct {
	Config string `help:"Alternate config file."`

	Version VersionCMD `cmd:"" help:"Print version information."`
	Status  StatusCMD  `cmd:"" help:"Print status of downloaded videos."`
	Update  UpdateCMD  `cmd:"" help:"Update video database."`
	Mark    MarkCMD    `cmd:"" help:"Mark video status."`
}

func Run(version string) {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Opts{Version: version})

	ctx.FatalIfErrorf(err)
}
