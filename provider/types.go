package provider

import (
	"errors"
	"path/filepath"
)

// Provider describes an interface for providing files
type Provider interface {
	Open() error
	Close() error
	GetLatestVersion() (string, error)
	Walk(walkFn filepath.WalkFunc) error
	Retrieve(srcPath string, destPath string) error
}

var (
	// ErrProviderUnavaiable is a generic error when a provider is not avaiable
	ErrProviderUnavaiable = errors.New("provider not avaiable")
)
