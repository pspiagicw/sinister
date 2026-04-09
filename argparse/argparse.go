package argparse

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/pspiagicw/sinister/manage"
)

type Opts struct {
	Version string
	Config  string
}

type VersionCMD struct {
}

func (v *VersionCMD) Run(o *Opts) error {
	fmt.Printf("sinister version: '%s'\n", o.Version)
	return nil
}

type StatusCMD struct {
}

func (s *StatusCMD) Run(o *Opts) error {
	manage.Status(o.Config)
	return nil
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
