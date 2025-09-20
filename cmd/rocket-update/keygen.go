package main

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/mouuff/go-rocket-update/internal/crypto"
	"github.com/mouuff/go-rocket-update/internal/fileio"
)

// Keygen describes the keygen subcommand
// this command is used to generate a private and a public key
type Keygen struct {
	flagSet *flag.FlagSet

	keyName string
}

// Name gets the name of the command
func (cmd *Keygen) Name() string {
	return "keygen"
}

// Init initializes the command
func (cmd *Keygen) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.keyName, "name", "id_rsa", "name of the key to generate")

	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Keygen) Run() error {

	privateKeyPath := cmd.keyName
	publicKeyPath := cmd.keyName + ".pub"

	if fileio.FileExists(privateKeyPath) {
		return fmt.Errorf("key '%s' already exists", privateKeyPath)
	}

	if fileio.FileExists(publicKeyPath) {
		return fmt.Errorf("key '%s' already exists", publicKeyPath)
	}

	log.Println("Generating keys...")
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("could not generate private key: %w", err)
	}

	privPem := crypto.ExportPrivateKeyAsPem(priv)
	err = os.WriteFile(privateKeyPath, privPem, 0600)
	if err != nil {
		return fmt.Errorf("could not write private key: %w", err)
	}
	log.Println("Created private key: " + privateKeyPath)

	pubPem, err := crypto.ExportPublicKeyAsPem(&priv.PublicKey)
	if err != nil {
		return fmt.Errorf("could not export public key: %w", err)
	}
	err = os.WriteFile(publicKeyPath, pubPem, 0644)
	if err != nil {
		return fmt.Errorf("could not write public key: %w", err)
	}

	log.Println("Created public key: " + publicKeyPath)

	return nil
}
