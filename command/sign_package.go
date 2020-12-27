package command

import (
	"flag"
	"fmt"
)

func NewSignPackage() *SignPackage {
	cmd := &SignPackage{
		name:    "sign",
		flagSet: flag.NewFlagSet("sign", flag.ContinueOnError),
	}

	cmd.flagSet.StringVar(&cmd.name, "path", "", "path to the package")

	return cmd
}

type SignPackage struct {
	flagSet *flag.FlagSet

	name string
}

func (cmd *SignPackage) Name() string {
	return cmd.flagSet.Name()
}

func (cmd *SignPackage) Init(args []string) error {
	return cmd.flagSet.Parse(args)
}

func (cmd *SignPackage) Run() error {
	fmt.Println("Hello", cmd.name, "!")
	return nil
}
