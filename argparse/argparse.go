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
	JSON    bool   `name:"json" help:"Print status in JSON format."`
	Creator string `name:"creator" help:"Show status for a single creator."`
}

func (s *StatusCMD) Run(o *Opts) error {
	manage.Status(manage.StatusOptions{
		ConfigPath: o.Config,
		JSON:       s.JSON,
		Creator:    s.Creator,
	})
	return nil
}

type UpdateCMD struct {
	URL    []string `name:"url" help:"Fetch these RSS feeds instead of config URLs."`
	Limit  int      `name:"limit" default:"0" help:"Process at most N feed entries per URL (0 = all)."`
	DryRun bool     `name:"dry-run" help:"Show what would be inserted without writing to the database."`
	JSON   bool     `name:"json" help:"Print update summary in JSON format."`
}

func (u *UpdateCMD) Run(o *Opts) error {
	manage.Update(manage.UpdateOptions{
		ConfigPath: o.Config,
		URLs:       u.URL,
		Limit:      u.Limit,
		DryRun:     u.DryRun,
		JSON:       u.JSON,
	})
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
