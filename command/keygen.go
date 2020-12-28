package command

import (
	"flag"
	"io/ioutil"

	"github.com/mouuff/go-rocket-update/crypto"
)

type Keygen struct {
	flagSet *flag.FlagSet

	keyName string
}

func (cmd *Keygen) Name() string {
	return "keygen"
}

func (cmd *Keygen) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.keyName, "name", "id_rsa", "name of the key to generate")

	return cmd.flagSet.Parse(args)
}

func (cmd *Keygen) Run() error {
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return err
	}

	privPem := crypto.ExportPrivateKeyAsPem(priv)
	err = ioutil.WriteFile(cmd.keyName, privPem, 0600)
	if err != nil {
		return err
	}

	pubPem, err := crypto.ExportPublicKeyAsPem(&priv.PublicKey)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cmd.keyName+".pub", pubPem, 0644)
	if err != nil {
		return err
	}

	return nil
}
