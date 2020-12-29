package rocketupdate

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Command interface {
	Init([]string) error
	Run() error
	Name() string
}

func runCommands(args []string) error {

	cmds := []Command{
		&Sign{},
		&Keygen{},
		&Verify{},
	}

	if len(args) < 1 {
		cmdNames := ""
		for i, cmd := range cmds {
			if i > 0 {
				cmdNames += ", "
			}
			cmdNames += cmd.Name()
		}
		return errors.New("You must pass a sub-command (" + cmdNames + ")")

	}

	subcommand := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			err := cmd.Init(args[1:])
			if err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
	if err := runCommands(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
