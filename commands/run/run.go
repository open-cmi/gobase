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
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/initial"
)

type RunCommand struct {
	ConfigFile string
}

func Wait() error {
	// 初始化后，等待信号或异常错误
	sigs := make(chan os.Signal, 1)

	// 注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	errChan := webserver.WaitError()

	select {
	case sig := <-sigs:
		logger.Infof("received signal %s, exiting...\n", sig.String())
		return nil
	case err := <-errChan:
		if err != nil {
			logger.Errorf("service encountered runtime error: %s\n", err.Error())
			return err
		}
		return nil
	}
}

func Fini() {
	webserver.Close()
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

	logger.Infof("server started successfully, press Ctrl+C to stop\n")
	return Wait()
}

func init() {
	subcommands.RegisterCommand("run", &RunCommand{})
	subcommands.DefaultCommand(&RunCommand{})
}
