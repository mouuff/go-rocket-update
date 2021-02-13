package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

// verifyInstallation verifies if the executable is installed correctly
func verifyInstallation(u *updater.Updater) error {
	latestVersion, err := u.Provider.GetLatestVersion()
	if err != nil {
		return err
	}
	executable, err := u.GetExecutable()
	if err != nil {
		return err
	}
	cmd := exec.Cmd{
		Path: executable,
		Args: []string{executable, "--verify"},
	}
	// Should be replaced with Output() as soon as test project is updated
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	strOutput := string(output)

	if !strings.Contains(strOutput, latestVersion) {
		return errors.New("Version not found in program output")
	}
	return nil
}

func selfUpdate(u *updater.Updater) error {
	err := verifyInstallation(u)
	updateStatus, err := u.Update()
	if err != nil {
		return err
	}
	if updateStatus == updater.Updated {
		if err := verifyInstallation(u); err != nil {
			log.Println(err)
			log.Println("Rolling back...")
			return u.Rollback()
		}
		log.Println("Installation OK")
	}
	return nil
}

func main() {
	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries_" + runtime.GOOS + ".zip",
		},
		ExecutableName: "go-rocket-update-example",
		Version:        "v0.3.0",
	}

	fmt.Println(u.Version)
	if len(os.Args) > 1 && os.Args[1] == "--verify" {
		os.Exit(0)
	}
	err := selfUpdate(u)
	if err != nil {
		log.Println(err)
	}
}
