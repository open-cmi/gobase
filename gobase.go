package gobase

import (
	"github.com/open-cmi/gobase/essential/subcommands"

	_ "github.com/open-cmi/gobase/commands"
)

func Run() error {
	return subcommands.Run()
}
