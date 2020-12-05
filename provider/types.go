package provider

import "path/filepath"

// Provider describes an interface for providing files
type Provider interface {
	Open() error
	Close() error
	GetVersion() (string, error)

	Walk(walkFn filepath.WalkFunc) error
	Retrieve(srcPath string, destPath string) error
}
