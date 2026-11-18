package up

import (
	"flag"
	"fmt"
	"os"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/essential/migrate"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/essential/subcommands"
)

type UpCommand struct {
}

var configfile string

func (c *UpCommand) Synopsis() string {
	return "up migration to databse"
}

func (c *UpCommand) Run() error {
	var format string
	var input string
	var count int
	upCmd := flag.NewFlagSet("up", flag.ContinueOnError)
	upCmd.StringVar(&input, "input", input, "if use sql, should specify sql directory")
	upCmd.StringVar(&format, "format", format, "format, go or sql")
	upCmd.StringVar(&configfile, "config", configfile, "config file")
	upCmd.IntVar(&count, "count", count, "migrate up count")

	err := upCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	configfile, err = config.ResolveConfigFile(configfile)
	if err != nil {
		upCmd.Usage()
		return err
	}

	err = config.Init(configfile)
	if err != nil {
		fmt.Printf("init config failed: %s\n", err.Error())
		return err
	}

	opt := migrate.NewUpOpt(sqldb.GetDB(), format, input, count)
	err = opt.Run()
	return err
}

func init() {
	subcommands.RegisterCommand("up", &UpCommand{})
}
