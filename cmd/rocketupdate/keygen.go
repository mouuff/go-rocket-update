package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"

	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
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

	privKeyPath := cmd.keyName
	pubKeyPath := cmd.keyName + ".pub"

	if fileio.FileExists(privKeyPath) {
		return errors.New("Key '" + privKeyPath + "' already exists")
	}

	if fileio.FileExists(pubKeyPath) {
		return errors.New("Key '" + pubKeyPath + "' already exists")
	}

	log.Println("Generating keys...")
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return err
	}

	privPem := crypto.ExportPrivateKeyAsPem(priv)
	err = ioutil.WriteFile(privKeyPath, privPem, 0600)
	if err != nil {
		return err
	}
	log.Println("Created private key: " + privKeyPath)

	pubPem, err := crypto.ExportPublicKeyAsPem(&priv.PublicKey)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(pubKeyPath, pubPem, 0644)
	if err != nil {
		return err
	}

	log.Println("Created public key: " + pubKeyPath)

	return nil
}
