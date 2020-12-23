package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mouuff/easy-update/fileio"
	"github.com/mouuff/easy-update/provider"
)

// Updater struct
type Updater struct {
	Provider   provider.Provider
	BinaryName string
	Version    string
}

// getBinaryName gets the name used to find the right binary
func (u *Updater) getBinaryName() string {
	return u.BinaryName + "_" + GetPlatformName()
}

// findBinaryProviderPath finds the right binary using the provider
func (u *Updater) findBinaryProviderPath() (string, error) {
	binaryPath := ""
	fmt.Println(u.getBinaryName())
	err := u.Provider.Walk(func(info *provider.FileInfo) error {
		if !info.Mode.IsDir() && strings.Contains(info.Path, u.getBinaryName()) {
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
	lastestVersion, err := u.Provider.GetLatestVersion()
	if err != nil {
		return false, err
	}
	if u.Version != lastestVersion {
		return true, nil
	}
	return false, nil
}

// Run runs the updater
// It will update the current application if an update is found
func (u *Updater) Run() error {
	if err := u.Provider.Open(); err != nil {
		return err
	}
	defer u.Provider.Close()
	canUpdate, err := u.CanUpdate()
	if err != nil {
		return err
	}
	if !canUpdate {
		return nil
	}
	binaryProviderPath, err := u.findBinaryProviderPath()
	if err != nil {
		return err
	}
	tmpDir, err := fileio.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	binaryPath := filepath.Join(tmpDir, filepath.Base(binaryProviderPath))
	err = u.Provider.Retrieve(binaryProviderPath, binaryPath)
	if err != nil {
		return err
	}
	err = ReplaceExecutableWith(binaryPath)
	if err != nil {
		return err
	}
	return nil
}
