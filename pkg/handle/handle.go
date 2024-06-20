package handle

import (
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/help"
	"github.com/pspiagicw/sinister/pkg/manage"
)

func Handle(opts *argparse.Opts) {

	if len(opts.Args) == 0 {
		help.Usage(opts.Version)
		goreland.LogFatal("No subcommand given")
	} else {
		handleCmd(opts)
	}
}

func handleCmd(opts *argparse.Opts) {

	handler := map[string]func(opts *argparse.Opts){
		"version": func(opts *argparse.Opts) {
			help.Version(opts.Version)
		},
		"help": func(opts *argparse.Opts) {
			help.HandleHelp(opts.Args[1:], opts.Version)
		},
		"status":   manage.Status,
		"update":   manage.Update,
		"download": manage.Download,
		"mark":     manage.Mark,
		"auto":     manage.Auto,
	}

	cmd := opts.Args[0]
	handleFunc, ok := handler[cmd]

	if !ok {
		help.Usage(opts.Version)
		goreland.LogFatal("subcommand %s not found", cmd)
	}

	handleFunc(opts)
}

func notImplemented(opts *argparse.Opts) {
	goreland.LogFatal("Not implemented yet!")
}
