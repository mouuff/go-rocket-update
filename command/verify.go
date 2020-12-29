package command

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/mouuff/go-rocket-update/crypto"
)

type Verify struct {
	flagSet *flag.FlagSet

	path   string
	pubkey string
}

func (cmd *Verify) Name() string {
	return "verify"
}

func (cmd *Verify) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(&cmd.path, "path", "", "path to the package to verify (required)")
	cmd.flagSet.StringVar(&cmd.pubkey, "pubkey", "", "path to the public key (required)")

	return cmd.flagSet.Parse(args)
}

func (cmd *Verify) Run() error {

	log.Println("Reading public key...")
	pubkeyBytes, err := ioutil.ReadFile(cmd.pubkey)
	if err != nil {
		return err
	}
	pubkey, err := crypto.ParsePemPublicKey(pubkeyBytes)
	if err != nil {
		return err
	}
	signaturesPath := filepath.Join(cmd.path, "signatures.json")
	log.Println("Reading " + signaturesPath + " ...")
	signaturesJSON, err := ioutil.ReadFile(signaturesPath)
	if err != nil {
		return err
	}
	signatures := &crypto.Signatures{}
	err = json.Unmarshal(signaturesJSON, signatures)
	if err != nil {
		return err
	}

	unverifiedFiles, err := signatures.VerifyFolder(pubkey, cmd.path)
	if err != nil {
		return err
	}
	if len(unverifiedFiles) <= 1 {
		// <= 1 because it is normal to have one unverified file because signatures.json isnt verified
		fmt.Println("All files verified!")
		return nil
	}
	return errors.New("Some files could not be verified:\n" + strings.Join(unverifiedFiles, "\n"))
}
