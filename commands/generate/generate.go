package generate

import (
	"errors"
	"flag"
	"os"

	"github.com/open-cmi/gobase/essential/migrate"
	"github.com/open-cmi/gobase/essential/subcommands"
)

type GenerateCommand struct {
}

func (c *GenerateCommand) Synopsis() string {
	return "generate migration file"
}

func (c *GenerateCommand) Run() error {
	var name string
	var format string
	var output string

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	generateCmd.StringVar(&output, "output", output, "output directory")
	generateCmd.StringVar(&format, "format", format, "format, go or sql")
	generateCmd.StringVar(&name, "name", name, "script name")

	generateCmd.Parse(os.Args[2:])

	if name == "" {
		generateCmd.Usage()
		return errors.New("name cant't be empty")
	}
	gen := migrate.NewGenerateOpt(name, format, output)
	err := gen.Run()
	return err
}

func init() {
	subcommands.RegisterCommand("generate", &GenerateCommand{})
}
