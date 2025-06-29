package commands

import (
	"fmt"

	"github.com/open-cmi/gobase/essential/subcommands"
)

var Version string

type VersionCommand struct {
}

func (c *VersionCommand) Synopsis() string {
	return "print version"
}

func (c *VersionCommand) Run() error {
	fmt.Printf("version: %s\n", Version)
	return nil
}

func init() {
	subcommands.RegisterCommand("version", &VersionCommand{})
}
