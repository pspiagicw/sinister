package main

import (
	"github.com/pspiagicw/sinister/pkg/argparse"
	"github.com/pspiagicw/sinister/pkg/handle"
)

var VERSION string = "unversioned"

func main() {
	args := argparse.Parse(VERSION)
	handle.Handle(args)
}
