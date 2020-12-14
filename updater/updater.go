package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mouuff/easy-update/provider"
)

type Updater struct {
	Provider   provider.Provider
	BinaryName string
	Version    string
}

func (u *Updater) getBinaryName() string {
	return u.BinaryName + "-" + GetPlatformName()
}

func (u *Updater) findBinaryProviderPath() (string, error) {
	binaryPath := ""
	err := u.Provider.Walk(func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", filePath, err)
			return err
		}
		if strings.Contains(filePath, u.getBinaryName()) {
			binaryPath = filePath
		}
		return nil
	})
	if err != nil {
		return binaryPath, err
	}
	return binaryPath, nil
}

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
	tmpDir, err := ioutil.TempDir("", "updater") //TODO replace appname
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
