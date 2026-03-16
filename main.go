package main

import (
	"github.com/pspiagicw/sinister/argparse"
	"github.com/pspiagicw/sinister/handle"
)

var VERSION string = "unversioned"

func main() {
	args := argparse.Run(VERSION)
}
