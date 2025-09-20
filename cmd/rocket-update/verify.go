package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"os"

	"github.com/mouuff/go-rocket-update/internal/constant"
	"github.com/mouuff/go-rocket-update/internal/crypto"
)

// Verify describes the verify subcommand
// this command is used to verify if all files are signed within a folder
type Verify struct {
	flagSet *flag.FlagSet

	path      string
	publicKey string
}

// Name gets the name of the command
func (cmd *Verify) Name() string {
	return "verify"
}

// Init initializes the command
func (cmd *Verify) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.path, "path", "", "path to the package directory to verify (required)")
	cmd.flagSet.StringVar(&cmd.publicKey, "publicKey", "", "path to the public key (required)")

	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Verify) Run() error {
	log.Println("Reading public key...")
	pubkeyBytes, err := os.ReadFile(cmd.publicKey)
	if err != nil {
		return fmt.Errorf("could not read public key: %w", err)
	}
	publicKey, err := crypto.ParsePemPublicKey(pubkeyBytes)
	if err != nil {
		return fmt.Errorf("could not parse public key: %w", err)
	}
	signaturesPath := filepath.Join(cmd.path, constant.SignatureRelPath)
	log.Println("Reading " + signaturesPath + " ...")

	signatures, err := crypto.LoadSignaturesFromJSON(signaturesPath)
	if err != nil {
		return fmt.Errorf("could not load signatures: %w", err)
	}
	unverifiedFiles, err := signatures.VerifyFolder(publicKey, cmd.path)
	if err != nil {
		return fmt.Errorf("could not verify folder: %w", err)
	}
	if len(unverifiedFiles) <= 1 {
		// <= 1 because it is normal to have one unverified file because signatures file isnt verified
		fmt.Println("All files verified!")
		return nil
	}
	return errors.New("Some files could not be verified:\n" + strings.Join(unverifiedFiles, "\n"))
}
