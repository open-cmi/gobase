package list

import (
	"flag"
	"os"

	"github.com/open-cmi/gobase/essential/migrate"
	"github.com/open-cmi/gobase/essential/subcommands"
)

type ListCommand struct {
}

func (c *ListCommand) Synopsis() string {
	return "list migrations in program"
}

func (c *ListCommand) Run() error {
	var format string
	var input string
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCmd.StringVar(&input, "input", input, "if use sql, should specify sql directory")
	listCmd.StringVar(&format, "format", format, "format, go or sql")

	err := listCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	opt := migrate.NewListOpt(format, input)
	err = opt.Run()
	return err
}

func init() {
	subcommands.RegisterCommand("list", &ListCommand{})
}
