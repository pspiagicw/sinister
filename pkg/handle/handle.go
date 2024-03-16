package handle

import (
	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/help"
	"github.com/pspiagicw/sinister/pkg/tui"
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
		},
		"help":     notImplemented,
		"status":   notImplemented,
		"update":   tui.Update,
		"download": tui.Download,
	}

	cmd := opts.Args[0]
	handleFunc, ok := handler[cmd]
	if !ok {
		help.Usage(opts.Version)
		goreland.LogError("subcommand %s not found", cmd)
	}
	handleFunc(opts)
}

func notImplemented(opts *argparse.Opts) {
	goreland.LogFatal("Not implemented yet!")
}
