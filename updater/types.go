package updater

import "github.com/mouuff/easy-update/provider"

type Updater interface {
	AddProvider(provider.Provider)
	Run() error
}
