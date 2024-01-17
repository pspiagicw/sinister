package main

import (
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/handle"
)

var VERSION string

func main() {
	args := argparse.ParseArgs(VERSION)
	handle.HandleArgs(args, VERSION)
}
