package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mouuff/go-rocket-update/command"
)

func runCommands(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []command.Command{
		command.NewSignPackage(),
	}

	subcommand := args[0]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(args[1:])
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
