package handle

import (
	"fmt"

	"github.com/pspiagicw/goreland"
	"github.com/pspiagicw/sinister/pkg/help"
	"github.com/pspiagicw/sinister/pkg/update"
)

func HandleArgs(args []string, VERSION string) {
	handler := map[string]func(args []string){
		"version": func(args []string) {
			fmt.Println(args)
		},
		"help": func(args []string) {
			goreland.LogFatal("Not implemented yet!")
		},
		"status": func(args []string) {
			goreland.LogFatal("Not implemented yet!")
		},
		"update": func(args []string) {
			update.Update(args)
		},
	}

	if len(args) == 0 {
		help.PrintHelp(VERSION)
	} else {
		cmd := args[0]
		handleCmd, ok := handler[cmd]
		if !ok {
			goreland.LogError("subcommand %s not found", cmd)
			help.PrintHelp(VERSION)
		}
		handleCmd(args)
	}
}
