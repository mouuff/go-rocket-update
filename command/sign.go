package command

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/mouuff/go-rocket-update/crypto"
)

type Sign struct {
	flagSet *flag.FlagSet

	path string
	key  string
}

func (cmd *Sign) Name() string {
	return "sign"
}

func (cmd *Sign) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.path, "path", "", "path to the package (required)")
	cmd.flagSet.StringVar(&cmd.key, "key", "", "path to the private key (required)")

	return cmd.flagSet.Parse(args)
}

func (cmd *Sign) Run() error {

	privkeyBytes, err := ioutil.ReadFile(cmd.key)
	if err != nil {
		return err
	}
	privkey, err := crypto.ParsePemPrivateKey(privkeyBytes)
	if err != nil {
		return err
	}
	signatures, err := crypto.GetFolderSignature(privkey, cmd.path)
	if err != nil {
		return err
	}
	fmt.Println("path: ", cmd.path, "!")
	fmt.Println("key: ", cmd.key, "!")
	fmt.Println(signatures)

	return nil
}
