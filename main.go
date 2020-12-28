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
		&command.SignPackage{},
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
