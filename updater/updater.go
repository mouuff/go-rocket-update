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
	provider   provider.Provider
	binaryName string
	version    string
}

func NewUpdater(p provider.Provider, binaryName, version string) *Updater {
	return &Updater{
		provider:   p,
		binaryName: binaryName,
		version:    version,
	}
}

func (u *Updater) getBinaryName() string {
	return u.binaryName + "-" + GetPlatformName()
}

func (u *Updater) findBinaryProviderPath() (string, error) {
	binaryPath := ""
	err := u.provider.Walk(func(filePath string, info os.FileInfo, err error) error {
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
	lastestVersion, err := u.provider.GetLatestVersion()
	if err != nil {
		return false, err
	}
	if u.version != lastestVersion {
		return true, nil
	}
	return false, nil
}

func (u *Updater) Run() error {
	if err := u.provider.Open(); err != nil {
		return err
	}
	defer u.provider.Close()
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

	err = u.provider.Retrieve(binaryProviderPath, binaryPath)
	if err != nil {
		return err
	}

	err = ReplaceExecutableWith(binaryPath)
	if err != nil {
		return err
	}

	return nil
}
