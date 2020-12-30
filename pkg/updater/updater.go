package updater

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mouuff/go-rocket-update/internal/fileio"
	"github.com/mouuff/go-rocket-update/pkg/provider"
)

// Updater struct
type Updater struct {
	Provider   provider.Provider
	BinaryName string
	Version    string
}

// getBinaryName gets the name used to find the right binary
func (u *Updater) getBinaryName() string {
	return u.BinaryName + "_" + runtime.GOOS + "_" + runtime.GOARCH
}

// findBinaryProviderPath finds the right binary using the provider
func (u *Updater) findBinaryProviderPath() (string, error) {
	binaryPath := ""
	err := u.Provider.Walk(func(info *provider.FileInfo) error {
		if info.Mode.IsRegular() && strings.Contains(filepath.Base(info.Path), u.getBinaryName()) {
			binaryPath = info.Path
		}
		return nil
	})
	if err != nil {
		return binaryPath, err
	}
	return binaryPath, nil
}

// CanUpdate checks if the updater found a new version
func (u *Updater) CanUpdate() (bool, error) {
	latestVersion, err := u.Provider.GetLatestVersion()
	if err != nil {
		return false, err
	}
	if u.Version != latestVersion {
		log.Printf("Found update: %s", latestVersion)
		return true, nil
	}
	return false, nil
}

// updateExecutable updates the current executable with the new one
func (u *Updater) updateExecutable() (err error) {
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return
	}
	defer os.RemoveAll(tmpDir)
	binaryProviderPath, err := u.findBinaryProviderPath()
	if err != nil {
		return
	}
	binaryTmpPath := filepath.Join(tmpDir, filepath.Base(binaryProviderPath))
	err = u.Provider.Retrieve(binaryProviderPath, binaryTmpPath)
	if err != nil {
		return
	}
	err = fileio.ReplaceExecutableWith(binaryTmpPath)
	return
}

// Run runs the updater
// It will update the current application if an update is found
func (u *Updater) Run() (err error) {
	canUpdate, err := u.CanUpdate()
	if err != nil || !canUpdate {
		return
	}
	log.Printf("Updating...")
	if err = u.Provider.Open(); err != nil {
		return
	}
	defer u.Provider.Close()

	err = u.updateExecutable()
	log.Printf("Updated!")
	return
}
