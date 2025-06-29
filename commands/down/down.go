package down

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/essential/migrate"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/essential/subcommands"
)

type DownCommand struct {
}

var configfile string

func (c *DownCommand) Synopsis() string {
	return "down database migration"
}

func (c *DownCommand) Run() error {

	var format string
	var input string
	var count int
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	downCmd.StringVar(&configfile, "config", configfile, "config file")
	downCmd.StringVar(&format, "format", format, "format, go or sql")
	downCmd.StringVar(&input, "input", input, "if use sql, should specify sql directory")
	downCmd.IntVar(&count, "count", count, "migrate count")

	err := downCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	if configfile == "" {
		return errors.New("config file must not be empty")
	}

	err = config.Init(configfile)
	if err != nil {
		fmt.Printf("init config failed: %s\n", err.Error())
		return err
	}
	downopt := migrate.NewDownOpt(sqldb.GetDB(), format, input, count)
	err = downopt.Run()
	return err
}

func init() {
	subcommands.RegisterCommand("down", &DownCommand{})
}
