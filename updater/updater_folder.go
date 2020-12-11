package updater

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/mouuff/easy-update/provider"
)

type UpdaterFolder struct {
	provider      provider.Provider
	projectFolder string
}

func NewUpdaterFolder() *UpdaterFolder {
	return &UpdaterFolder{}
}

func (u *UpdaterFolder) SetProjectFolder(projectFolder string) {
	u.projectFolder = projectFolder
}

func (u *UpdaterFolder) SetProvider(p provider.Provider) {
	u.provider = p
}

func (u *UpdaterFolder) Run() error {

	err := u.provider.Walk(func(relPath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", relPath, err)
			return err
		}
		destPath := path.Join(u.projectFolder, relPath)
		if info.IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			err = u.provider.Retrieve(relPath, destPath)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
