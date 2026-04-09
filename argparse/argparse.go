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
	URL       []string `name:"url" help:"Fetch these RSS feeds instead of config URLs."`
	Limit     int      `name:"limit" default:"0" help:"Process at most N feed entries per URL (0 = all)."`
	SinceDays int      `name:"since-days" default:"0" help:"Only process videos published in the last N days (0 = no filter)."`
	Retries   int      `name:"retries" default:"2" help:"Retry failed feed fetches this many times."`
	Timeout   int      `name:"timeout" default:"30" help:"HTTP timeout in seconds for each feed request."`
	DryRun    bool     `name:"dry-run" help:"Show what would be inserted without writing to the database."`
	JSON      bool     `name:"json" help:"Print update summary in JSON format."`
}

func (u *UpdateCMD) Run(o *Opts) error {
	manage.Update(manage.UpdateOptions{
		ConfigPath: o.Config,
		URLs:       u.URL,
		Limit:      u.Limit,
		SinceDays:  u.SinceDays,
		Retries:    u.Retries,
		Timeout:    u.Timeout,
		DryRun:     u.DryRun,
		JSON:       u.JSON,
	})
	return nil
}

type MarkCMD struct {
	Slug             []string `name:"slug" help:"Mark entries by slug."`
	URL              []string `name:"url" help:"Mark entries by video URL."`
	Creator          string   `name:"creator" help:"Mark entries for this creator."`
	AllUnwatched     bool     `name:"all-unwatched" help:"Mark all unwatched entries."`
	MarkAllUnwatched bool     `name:"mark-all-unwatched" help:"Mark every video in the database as unwatched."`
	DryRun           bool     `name:"dry-run" help:"Show what would be marked without updating the database."`
}

func (m *MarkCMD) Run(o *Opts) error {
	manage.Mark(manage.MarkOptions{
		Slugs:            m.Slug,
		URLs:             m.URL,
		Creator:          m.Creator,
		AllUnwatched:     m.AllUnwatched,
		MarkAllUnwatched: m.MarkAllUnwatched,
		DryRun:           m.DryRun,
	})
	return nil
}

type ExportCMD struct {
}

func (e *ExportCMD) Run(o *Opts) error {
	manage.Export()
	return nil
}

type DownloadCMD struct {
	Days   int `name:"days" default:"0" help:"Download only videos from the last N days (0 = no filter)."`
	Videos int `name:"videos" default:"0" help:"Download only the latest N unwatched videos (0 = no limit)."`
}

func (d *DownloadCMD) Run(o *Opts) error {
	manage.Download(manage.DownloadOptions{
		ConfigPath: o.Config,
		Days:       d.Days,
		Videos:     d.Videos,
	})
	return nil
}

type ListCMD struct {
}

func (l *ListCMD) Run(o *Opts) error {
	manage.List()
	return nil
}

type AddCMD struct {
	VideoURL string `arg:"" name:"video-url" help:"Any YouTube video URL."`
}

func (a *AddCMD) Run(o *Opts) error {
	manage.Add(manage.AddOptions{
		ConfigPath: o.Config,
		VideoURL:   a.VideoURL,
	})
	return nil
}

type CleanCMD struct {
	DryRun bool `name:"dry-run" help:"Preview which 404 feeds would be removed."`
}

func (c *CleanCMD) Run(o *Opts) error {
	manage.Clean(manage.CleanOptions{
		ConfigPath: o.Config,
		DryRun:     c.DryRun,
	})
	return nil
}

type SyncCMD struct {
	URL       []string `name:"url" help:"Fetch these RSS feeds instead of config URLs."`
	Limit     int      `name:"limit" default:"0" help:"Process at most N feed entries per URL (0 = all)."`
	SinceDays int      `name:"since-days" default:"0" help:"Only process videos published in the last N days (0 = no filter)."`
	Retries   int      `name:"retries" default:"2" help:"Retry failed feed fetches this many times."`
	Timeout   int      `name:"timeout" default:"30" help:"HTTP timeout in seconds for each feed request."`
	DryRun    bool     `name:"dry-run" help:"Run update in dry-run mode (download will use existing unwatched data)."`
	JSON      bool     `name:"json" help:"Print update summary in JSON format."`
	Days      int      `name:"days" default:"0" help:"Download only videos from the last N days (0 = no filter)."`
	Videos    int      `name:"videos" default:"0" help:"Download only the latest N unwatched videos (0 = no limit)."`
}

func (s *SyncCMD) Run(o *Opts) error {
	manage.Sync(manage.SyncOptions{
		ConfigPath: o.Config,
		Update: manage.UpdateOptions{
			URLs:      s.URL,
			Limit:     s.Limit,
			SinceDays: s.SinceDays,
			Retries:   s.Retries,
			Timeout:   s.Timeout,
			DryRun:    s.DryRun,
			JSON:      s.JSON,
		},
		Download: manage.DownloadOptions{
			Days:   s.Days,
			Videos: s.Videos,
		},
	})
	return nil
}

var CLI struct {
	Config string `help:"Alternate config file."`

	Version  VersionCMD  `cmd:"" help:"Print version information."`
	Status   StatusCMD   `cmd:"" help:"Print status of downloaded videos."`
	Update   UpdateCMD   `cmd:"" help:"Update video database."`
	Download DownloadCMD `cmd:"" help:"Download unwatched videos in highest available quality."`
	Mark     MarkCMD     `cmd:"" help:"Mark video status."`
	Export   ExportCMD   `cmd:"" help:"Export unwatched video URLs to urls.txt and mark them watched."`
	List     ListCMD     `cmd:"" help:"List channels and video counts."`
	Add      AddCMD      `cmd:"" help:"Add a channel feed URL to config from a YouTube video URL."`
	Clean    CleanCMD    `cmd:"" help:"Remove RSS feed URLs from config that return 404."`
	Sync     SyncCMD     `cmd:"" help:"Run update and then download in one command."`
}

func Run(version string) {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Opts{
		Version: version,
		Config:  CLI.Config,
	})

	ctx.FatalIfErrorf(err)
}
