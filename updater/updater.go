package updater

import (
	"github.com/mouuff/easy-update/provider"
)

type Updater struct {
	provider provider.Provider
	version  string
}

func NewUpdater(p provider.Provider, version string) *Updater {
	return &Updater{
		provider: p,
		version:  version,
	}
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

	return nil
}
