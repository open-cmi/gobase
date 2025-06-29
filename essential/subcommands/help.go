package subcommands

import (
	"fmt"
	"os"
	"path"
)

type HelpCommand struct {
}

func (c *HelpCommand) Synopsis() string {
	return "print usage"
}

func (c *HelpCommand) Run() error {
	prog := path.Base(os.Args[0])
	fmt.Printf("Usage: %s <subcommand> <subcommand args>\n\n", prog)

	fmt.Printf("Subcommands:\n")
	for name, command := range cmdmapping {
		fmt.Printf("\t%-15s\t%s\n", name, command.Synopsis())
	}
	return nil
}

func init() {
	RegisterCommand("help", &HelpCommand{})
}
