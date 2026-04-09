package main

import (
	"github.com/pspiagicw/sinister/argparse"
)

var VERSION string = "unversioned"

func main() {
	argparse.Run(VERSION)
}
