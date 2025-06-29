package run

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/subcommands"
	"github.com/open-cmi/gobase/initial"
)

type RunCommand struct {
	ConfigFile string
}

func Wait() {
	// 初始化后，等待信号
	sigs := make(chan os.Signal, 1)

	//注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigs
}

func Fini() {
}

func (c *RunCommand) Synopsis() string {
	return "run main program"
}

func (c *RunCommand) Run() error {

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.StringVar(&c.ConfigFile, "config", c.ConfigFile, "config file")

	err := runCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	if c.ConfigFile == "" {
		runCmd.Usage()
		return errors.New("config file must not be empty")
	}

	err = config.Init(c.ConfigFile)
	if err != nil {
		logger.Errorf("config init failed: %s\n", err.Error())
		return err
	}

	err = initial.Init()
	if err != nil {
		logger.Errorf("%s\n", err.Error())
		return err
	}

	defer Fini()

	logger.Infof("start successfully, wait for signal to terminate progress")
	Wait()
	return nil
}

func init() {
	subcommands.RegisterCommand("run", &RunCommand{})
	subcommands.DefaultCommand(&RunCommand{})
}
