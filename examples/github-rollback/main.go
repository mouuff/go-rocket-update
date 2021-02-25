package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

// verifyInstallation verifies if the executable is installed correctly
// we are going to run the newly installed program by running it with -version
// if it outputs the good version then we assume the installation is good
func verifyInstallation(u *updater.Updater) error {
	latestVersion, err := u.GetLatestVersion()
	if err != nil {
		return err
	}
	executable, err := u.GetExecutable()
	if err != nil {
		return err
	}
	cmd := exec.Cmd{
		Path: executable,
		Args: []string{executable, "-version"},
	}
	// Should be replaced with Output() as soon as test project is updated
	output, err := cmd.Output()
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
		log.Println("Updated to latest version!")
	}
	return nil
}

func main() {
	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
		},
		ExecutableName: fmt.Sprintf("go-rocket-update-example_%s_%s", runtime.GOOS, runtime.GOARCH),
		Version:        "v0.0.0",
	}

	versionFlag := false
	flag.BoolVar(&versionFlag, "version", false, "prints the version and exit")
	flag.Parse()

	if versionFlag {
		// we use this flag to verify the installation for this example:
		// https://github.com/mouuff/go-rocket-update/blob/master/examples/github-rollback/main.go
		fmt.Println(u.Version)
		return
	}

	log.Println("Current version: " + u.Version)
	log.Println("Looking for updates...")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := selfUpdate(u); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	// you can add code here, it should run during the update process without conflicts
	wg.Wait()
}
