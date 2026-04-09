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
	Slug         []string `name:"slug" help:"Mark entries by slug."`
	URL          []string `name:"url" help:"Mark entries by video URL."`
	Creator      string   `name:"creator" help:"Mark entries for this creator."`
	AllUnwatched bool     `name:"all-unwatched" help:"Mark all unwatched entries."`
	DryRun       bool     `name:"dry-run" help:"Show what would be marked without updating the database."`
}

func (m *MarkCMD) Run(o *Opts) error {
	manage.Mark(manage.MarkOptions{
		Slugs:        m.Slug,
		URLs:         m.URL,
		Creator:      m.Creator,
		AllUnwatched: m.AllUnwatched,
		DryRun:       m.DryRun,
	})
	return nil
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
