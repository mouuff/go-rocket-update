package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// SubCommand defines the interface to implement new subcommands
type SubCommand interface {
	Init([]string) error
	Run() error
	Name() string
}

// RunSubCommand runs the right subcommand from the args
// Example: args {"sign", "-name", "test"} will run the sign command
// with {"-name", "test"} as parameters
func RunSubCommand(args []string) error {

	cmds := []SubCommand{
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

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {
	if err := RunSubCommand(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
