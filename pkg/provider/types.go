package provider

import (
	"errors"
	"os"
)

// A FileInfo describes a file given by a provider
type FileInfo struct {
	Path string
	Mode os.FileMode
}

// WalkFunc is the type of the function called for each file or directory
// visited by Walk.
// path is relative
type WalkFunc func(info *FileInfo) error

// AccessProvider describes the access methods of a Provider
// This methods shouldn't change the state of the provider
type AccessProvider interface {
	GetLatestVersion() (string, error) // Get the latest version (Should be accessible even if Open() was not called)
	Walk(walkFn WalkFunc) error
	Retrieve(srcPath string, destPath string) error
}

// A Provider describes an interface for providing files
type Provider interface {
	AccessProvider
	Open() error
	Close() error
}

var (
	// ErrProviderUnavailable is a generic error when a provider is not available
	ErrProviderUnavailable = errors.New("provider not available")
	// ErrFileNotFound is a generic error when a file is not found
	ErrFileNotFound = errors.New("file not found")

	// ErrNotOpenned is a generic error when a provider have been called while it is not openned
	ErrNotOpenned = errors.New("provider have not been openned")
)
