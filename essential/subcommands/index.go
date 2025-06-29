package subcommands

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Command interface {
	Run() error
	Synopsis() string
}

var cmdmapping map[string]Command = make(map[string]Command)
var RootCommand Command

func DefaultCommand(cmd Command) {
	RootCommand = cmd
}

func RegisterCommand(key string, cmd Command) error {
	_, found := cmdmapping[key]
	if found {
		return errors.New("command has been registered")
	}

	cmdmapping[key] = cmd

	return nil
}

func Run() error {
	prog := path.Base(os.Args[0])
	if len(os.Args) < 2 {
		return fmt.Errorf("see '%s help'", prog)
	}

	var err error
	key := os.Args[1]
	cmd, ok := cmdmapping[key]
	if ok {
		err = cmd.Run()
	} else {
		return fmt.Errorf("'%s' is not a sub command\nsee '%s help'", os.Args[1], prog)
	}

	return err
}
