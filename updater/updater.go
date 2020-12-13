package updater

import (
	"runtime"

	"github.com/kardianos/osext"
	"github.com/mouuff/easy-update/provider"
)

func GetPlatformName() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

func GetExecutable() (string, error) {
	execPath, err := osext.Executable()
	if err != nil {
		return "", err
	}
	return execPath, nil
}

type Updater struct {
	provider provider.Provider
}

func NewUpdater(p provider.Provider) *Updater {
	return &Updater{
		provider: p,
	}
}

func (u *Updater) Run() error {

	return nil
}
