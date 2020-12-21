package provider

import (
	"errors"
)

// WalkFunc is the type of the function called for each file or directory
// visited by Walk.
// path is relative
type WalkFunc func(path string, isDir bool) error

// Provider describes an interface for providing files
type Provider interface {
	Open() error
	Close() error
	GetLatestVersion() (string, error)
	Walk(walkFn WalkFunc) error
	Retrieve(srcPath string, destPath string) error
}

var (
	// ErrProviderUnavaiable is a generic error when a provider is not avaiable
	ErrProviderUnavaiable = errors.New("provider not avaiable")
)
