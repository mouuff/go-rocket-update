package fileio_test

import (
	"github.com/mouuff/easy-update/provider"
)

type UpdaterFolder struct {
	provider provider.Provider
}

func (u *UpdaterFolder) SetProvider(p provider.Provider) {
	u.provider = p
}

func (u *UpdaterFolder) Run() error {
	return nil
}
