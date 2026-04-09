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

func (u *UpdateCMD) Run(o *Opts) error {
	manage.Update(o.Config)
	return nil
}

type MarkCMD struct {
}

type ExportCMD struct {
}

func (e *ExportCMD) Run(o *Opts) error {
	manage.Export()
	return nil
}

var CLI struct {
	Config string `help:"Alternate config file."`

	Version VersionCMD `cmd:"" help:"Print version information."`
	Status  StatusCMD  `cmd:"" help:"Print status of downloaded videos."`
	Update  UpdateCMD  `cmd:"" help:"Update video database."`
	Mark    MarkCMD    `cmd:"" help:"Mark video status."`
	Export  ExportCMD  `cmd:"" help:"Export unwatched video URLs to urls.txt and mark them watched."`
}

func Run(version string) {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Opts{Version: version})

	ctx.FatalIfErrorf(err)
}
