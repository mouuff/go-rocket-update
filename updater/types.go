package updater

import "github.com/mouuff/easy-update/provider"

type Runnable interface {
	Run() error
}

type Updater interface {
	Runnable
	SetProvider(provider.Provider)
}
