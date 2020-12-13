package updater

import (
	"github.com/mouuff/easy-update/provider"
)

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
